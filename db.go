// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/v2/id"
)

// Version 数据文件的版本号
const Version = 1

// ErrIncompatible 数据文件版本不兼容
//
// 当数据文件中指定的版本号与当前的 Version 不相等时，返回此错误。
var ErrIncompatible = errors.New("数据文件版本不兼容")

// DB 区域数据库信息
//
// 数据格式：
//
//	1:[versions]:{id:name:yearIndex:size{}}
//
//	- 1 表示数据格式的版本，采用当前包的 Version 常量；
//	- versions 表示当前数据文件中的数据支持的年份列表，以逗号分隔；
//	- id 当前区域的 ID；
//	- name 当前区域的名称；
//	- yearIndex 此条数据支持的年份列表，每一个位表示一个年份在 versions 中的索引值；
//	- size 表示子元素的数量；
type DB struct {
	root     *Region
	versions []int // 支持的版本

	// 以下数据不会写入数据文件中

	fullNameSeparator string
	districts         []*Region

	// Load 指定的过滤版本，仅在 unmarshal 过程中使用，
	// 在完成 unmarshal 之的清空。
	filters []int
}

// NewDB 返回空的 [DB] 对象
func NewDB() *DB {
	db := &DB{versions: []int{}}
	db.root = &Region{db: db}
	return db
}

// Version 当前这份数据支持的年份列表
func (db *DB) Versions() []int { return db.versions }

// 指定年份在 Versions 中的下标
//
// 如果不存在，返回 -1
func (db *DB) versionIndex(ver int) int {
	for i, v := range db.versions { // TODO(go1.21) slices.IndexFunc
		if v == ver {
			return i
		}
	}
	return -1
}

// AddVersion 添加新的版本号
func (db *DB) AddVersion(ver int) (ok bool) {
	if db.versionIndex(ver) > -1 { // 检测 ver 是否已经存在
		return false
	}

	db.versions = append(db.versions, ver)
	return true
}

// Find 查找指定 ID 对应的信息
func (db *DB) Find(regionID string) *Region { return db.root.findItem(id.SplitFilter(regionID)...) }

var levelIndex = []id.Level{id.Province, id.City, id.County, id.Town, id.Village}

// AddItem 添加一条子项
func (db *DB) AddItem(regionID, name string, ver int) error {
	list := id.SplitFilter(regionID)
	item := db.root.findItem(list...)

	if item == nil {
		items := list[:len(list)-1] // 上一级
		item = db.root.findItem(items...)
		level := levelIndex[len(items)]
		return item.addItem(list[len(list)-1], name, level, ver)
	}

	return item.setSupported(ver)
}

func (db *DB) marshal() ([]byte, error) {
	versions := make([]string, 0, len(db.versions))
	for _, v := range db.versions {
		versions = append(versions, strconv.Itoa(v))
	}

	buf := errwrap.Buffer{Buffer: bytes.Buffer{}}
	buf.WString(strconv.Itoa(Version)).WByte(':')

	buf.WByte('[')
	buf.WString(strings.Join(versions, ","))
	buf.WByte(']').WByte(':')

	err := db.root.marshal(&buf)
	if err != nil {
		return nil, err
	}

	if buf.Err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (db *DB) unmarshal(data []byte) error {
	data, val := indexBytes(data, ':')
	ver, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if ver != Version {
		return ErrIncompatible
	}

	data, val = indexBytes(data, ':')
	versions := strings.Split(strings.Trim(val, "[]"), ",")
	db.versions = make([]int, 0, len(versions))
	for _, version := range versions {
		v, err := strconv.Atoi(version)
		if err != nil {
			return err
		}
		db.versions = append(db.versions, v)
	}

	if len(db.filters) == 0 {
		db.filters = db.versions
	} else {
	LOOP:
		for _, v := range db.filters {
			for _, v2 := range db.versions {
				if v2 == v {
					continue LOOP
				}
			}
			return fmt.Errorf("当前数据文件没有 %d 年份的数据", v)
		}
	}

	defer func() {
		db.versions = db.filters
		db.filters = db.filters[:0]
	}()

	db.root = &Region{db: db}
	return db.root.unmarshal(data, "", "", 0)
}

func (db *DB) filterVersions(versions []int) []int {
	vers := make([]int, 0, len(versions))
LOOP:
	for _, v := range versions {
		for _, v2 := range db.filters {
			if v2 == v {
				vers = append(vers, v)
				continue LOOP
			}
		}
	}

	return vers
}
