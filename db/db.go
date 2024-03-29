// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package db 提供区域数据文件的相关操作
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
package db

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/id"
)

// Version 数据文件的版本号
const Version = 1

// ErrIncompatible 数据文件版本不兼容
//
// 当数据文件中指定的版本号与当前的 Version 不相等时，返回此错误。
var ErrIncompatible = errors.New("数据文件版本不兼容")

// DB 区域数据库信息
type DB struct {
	root     *Region
	versions []int // 支持的版本

	// 以下数据不会写入数据文件中

	fullNameSeparator string

	// Load 指定的过滤版本，仅在 unmarshal 过程中使用，
	// 在完成 unmarshal 之的清空。
	filters []int
}

// New 返回 DB 的空对象
func New() *DB {
	db := &DB{
		versions: []int{},
	}
	db.root = &Region{db: db}

	return db
}

// LoadFS 从数据文件加载数据
func LoadFS(f fs.FS, file, separator string, compress bool, version ...int) (*DB, error) {
	data, err := fs.ReadFile(f, file)
	if err != nil {
		return nil, err
	}
	return Load(data, separator, compress, version...)
}

// Load 将数据内容加载至 DB 对象
//
// version 仅加载指定年份的数据，如果为空，则加载所有数据；
func Load(data []byte, separator string, compress bool, version ...int) (*DB, error) {
	if compress {
		rd, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		data, err = io.ReadAll(rd)
		if err != nil {
			return nil, err
		}
	}

	db := &DB{
		fullNameSeparator: separator,
		filters:           version,
	}
	if err := db.unmarshal(data); err != nil {
		return nil, err
	}
	return db, nil
}

// LoadFile 从数据文件加载数据
func LoadFile(file, separator string, compress bool, version ...int) (*DB, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Load(data, separator, compress, version...)
}

// Dump 输出到文件
func (db *DB) Dump(file string, compress bool) error {
	data, err := db.marshal()
	if err != nil {
		return err
	}

	if compress {
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err = w.Write(data); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}

		data = buf.Bytes()
	}

	return os.WriteFile(file, data, os.ModePerm)
}

func (db *DB) Versions() []int { return db.versions }

// Unmarshal 解码 data 至 DB
//
// Deprecated: 请使用 Load 代替
func Unmarshal(data []byte, separator string, version ...int) (*DB, error) {
	return Load(data, separator, false, version...)
}

// VersionIndex 指定年份在 Versions 中的下标
//
// 如果不存在，返回 -1
func (db *DB) VersionIndex(ver int) int {
	for i, v := range db.versions {
		if v == ver {
			return i
		}
	}
	return -1
}

// AddVersion 添加新的版本号
func (db *DB) AddVersion(ver int) (ok bool) {
	if db.VersionIndex(ver) > -1 { // 检测 ver 是否已经存在
		return false
	}

	db.versions = append(db.versions, ver)
	return true
}

// Find 查找指定 ID 对应的信息
func (db *DB) Find(id ...string) *Region { return db.root.findItem(id...) }

var levelIndex = []id.Level{id.Province, id.City, id.County, id.Town, id.Village}

// AddItem 添加一条子项
func (db *DB) AddItem(regionID, name string, ver int) error {
	list := id.SplitFilter(regionID)
	item := db.Find(list...)

	if item == nil {
		items := list[:len(list)-1] // 上一级
		item = db.Find(items...)
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
