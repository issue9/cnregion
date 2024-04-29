// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/v2/id"
)

// Region 表示单个区域
type Region struct {
	id       string
	name     string
	items    []*Region
	versions []int // 支持的版本号列表

	// 以下数据不会写入数据文件中

	fullName string // 全名
	fullID   string
	db       *DB
	level    id.Level
}

// Provinces 省份列表
func (db *DB) Provinces() []*Region { return db.root.items }

func (r *Region) ID() string       { return r.id }       // 区域的 ID，不包括后缀 0 和上一级的 ID
func (r *Region) Name() string     { return r.name }     // 区域的名称
func (r *Region) FullName() string { return r.fullName } // 区域的全称，包括上一级的名称
func (r *Region) FullID() string   { return r.fullID }   // 区域的 ID，包括后缀的 0 以及上一级的 ID，长度为 12
func (r *Region) Versions() []int  { return r.versions } // 支持的年份版本
func (r *Region) Items() []*Region { return r.items }    // 子项

// IsSupported 当前数据是否支持该年份
func (reg *Region) IsSupported(ver int) bool {
	for _, y := range reg.versions {
		if y == ver {
			return true
		}
	}
	return false
}

func (reg *Region) addItem(id, name string, level id.Level, ver int) error {
	if index := reg.db.versionIndex(ver); index == -1 {
		return fmt.Errorf("不支持该年份 %d 的数据", ver)
	}

	for _, item := range reg.items {
		if item.id == id {
			return fmt.Errorf("已经存在相同 ID 的数据项：%s", id)
		}
	}

	reg.items = append(reg.items, &Region{
		id:       id,
		name:     name,
		db:       reg.db,
		level:    level,
		versions: []int{ver},
	})
	return nil
}

func (reg *Region) setSupported(ver int) error {
	index := reg.db.versionIndex(ver)
	if index == -1 {
		return fmt.Errorf("不存在该年份 %d 的数据", ver)
	}

	if !reg.IsSupported(ver) {
		reg.versions = append(reg.versions, ver)
	}
	return nil
}

func (reg *Region) findItem(regionID ...string) *Region {
	if len(regionID) == 0 {
		return reg
	}

	for _, item := range reg.items {
		if item.id == regionID[0] {
			return item.findItem(regionID[1:]...)
		}
	}

	return nil
}

func (reg *Region) marshal(buf *errwrap.Buffer) error {
	supported := 0
	for _, ver := range reg.versions {
		index := reg.db.versionIndex(ver)
		if index == -1 {
			return fmt.Errorf("无效的年份 %d 位于 %s", ver, reg.fullName)
		}
		supported += 1 << index
	}
	buf.Printf("%s:%s:%d:%d{", reg.id, reg.name, supported, len(reg.items))
	for _, item := range reg.items {
		err := item.marshal(buf)
		if err != nil {
			return err
		}
	}
	buf.WByte('}')

	return nil
}

func (reg *Region) unmarshal(data []byte, parentName, parentID string, level id.Level) error {
	reg.level = level

	data, reg.id = indexBytes(data, ':')

	data, reg.name = indexBytes(data, ':')
	reg.fullName = reg.name
	if parentName != "" {
		reg.fullName = parentName + reg.db.fullNameSeparator + reg.name
	}
	parentID += reg.id
	reg.fullID = id.Fill(parentID, id.Village)

	// Versions
	data, val := indexBytes(data, ':')
	supported, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	versions := make([]int, 0, len(reg.db.versions))
	for i, v := range reg.db.versions {
		if flag := 1 << i; flag&supported == flag {
			versions = append(versions, v)
		}
	}
	reg.versions = reg.db.filterVersions(versions)

	data, val = indexBytes(data, '{')
	size, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if size > 0 {
		for i := 0; i < size; i++ {
			index := findEnd(data)
			if index < 0 {
				return errors.New("未找到结束符号 }")
			}

			// 下一级的 Level
			var next id.Level
			if level == 0 {
				next = id.Province
			} else {
				next = level >> 1
			}

			item := &Region{db: reg.db}
			if err := item.unmarshal(data[:index], reg.fullName, parentID, next); err != nil {
				return err
			}
			if len(item.versions) > 0 { // 表示该条数据不支持所有的年份
				reg.items = append(reg.items, item)
			}
			data = data[index+1:]
		}
	}

	return nil
}

func indexBytes(data []byte, b byte) ([]byte, string) {
	index := bytes.IndexByte(data, b)
	if index == -1 {
		panic(fmt.Sprintf("在%s未找到：%s", string(data), string(b)))
	}

	return data[index+1:], string(data[:index])
}

func findEnd(data []byte) int {
	deep := 0
	for i, b := range data {
		switch b {
		case '{':
			deep++
		case '}':
			deep--
			if deep == 0 {
				return i
			}
		}
	}

	return 0
}
