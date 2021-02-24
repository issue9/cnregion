// SPDX-License-Identifier: MIT

package db

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/id"
)

// Region 表示单个区域
type Region struct {
	ID        string
	Name      string
	supported int // 支持的版本号
	Items     []*Region

	// 以下数据不会写入数据文件中

	FullName string // 全名
	FullID   string
	db       *DB
	level    id.Level
}

// IsSupported 当前数据是否支持该年份
func (reg *Region) IsSupported(year int) bool {
	index := reg.db.VersionIndex(year)
	if index == -1 {
		return false
	}

	flag := 1 << index
	return reg.supported&flag == flag
}

func (reg *Region) addItem(id, name string, level id.Level, year int) error {
	index := reg.db.VersionIndex(year)
	if index == -1 {
		return fmt.Errorf("不支持该年份 %d 的数据", year)
	}

	for _, item := range reg.Items {
		if item.ID == id {
			return fmt.Errorf("已经存在相同 ID 的数据项：%s", id)
		}
	}

	reg.Items = append(reg.Items, &Region{
		ID:        id,
		Name:      name,
		supported: 1 << index,
		db:        reg.db,
		level:     level,
	})
	return nil
}

func (reg *Region) setSupported(year int) error {
	index := reg.db.VersionIndex(year)
	if index == -1 {
		return fmt.Errorf("不存在该年份 %d 的数据", year)
	}

	flag := 1 << index
	if reg.supported&flag == 0 {
		reg.supported += flag
	}
	return nil
}

func (reg *Region) findItem(regionID ...string) *Region {
	if len(regionID) == 0 {
		return reg
	}

	for _, item := range reg.Items {
		if item.ID == regionID[0] {
			return item.findItem(regionID[1:]...)
		}
	}

	return nil
}

func (reg *Region) marshal(buf *errwrap.Buffer) error {
	buf.Printf("%s:%s:%d:%d{", reg.ID, reg.Name, reg.supported, len(reg.Items))
	for _, item := range reg.Items {
		err := item.marshal(buf)
		if err != nil {
			return err
		}
	}
	buf.WByte('}')

	return nil
}

func (reg *Region) unmarshal(data []byte, parentName, parentID string) error {
	data, reg.ID = indexBytes(data, ':')

	data, reg.Name = indexBytes(data, ':')
	reg.FullName = reg.Name
	if parentName != "" {
		reg.FullName = parentName + reg.db.fullNameSeparator + reg.Name
	}
	parentID += reg.ID
	reg.FullID = id.Fill(parentID, id.Village)

	data, val := indexBytes(data, ':')
	supperted, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	reg.supported = supperted

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

			item := &Region{db: reg.db}
			if err := item.unmarshal(data[:index], reg.FullName, parentID); err != nil {
				return err
			}
			reg.Items = append(reg.Items, item)
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
