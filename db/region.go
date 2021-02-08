// SPDX-License-Identifier: MIT

package db

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/issue9/errwrap"
)

// Region 表示单个区域
type Region struct {
	ID        string
	Name      string
	Supported int // 支持的版本号
	Items     []*Region
}

func (reg *Region) marshal(buf *errwrap.Buffer) error {
	buf.Printf("%s:%s:%d:%d{", reg.ID, reg.Name, reg.Supported, len(reg.Items))
	for _, item := range reg.Items {
		err := item.marshal(buf)
		if err != nil {
			return err
		}
	}
	buf.WByte('}')

	return nil
}

func (reg *Region) unmarshal(data []byte) error {
	index := indexByte(data, ':')
	reg.ID = string(data[:index])

	data = data[index+1:]
	index = indexByte(data, ':')
	reg.Name = string(data[:index])
	data = data[index+1:]

	index = indexByte(data, ':')
	supperted, err := strconv.Atoi(string(data[:index]))
	if err != nil {
		return err
	}
	reg.Supported = supperted
	data = data[index+1:]

	index = indexByte(data, '{')
	size, err := strconv.Atoi(string(data[:index]))
	if err != nil {
		return err
	}
	data = data[index+1:]

	if size > 0 {
		for i := 0; i < size; i++ {
			index := findEnd(data)
			if index < 0 {
				return errors.New("未找到结束符号 }")
			}

			item := &Region{}
			if err := item.unmarshal(data[:index]); err != nil {
				return err
			}
			reg.Items = append(reg.Items, item)
			data = data[index+1:]
		}
	}

	return nil
}

func indexByte(data []byte, b byte) int {
	index := bytes.IndexByte(data, b)
	if index == -1 {
		panic(fmt.Sprintf("在%s未找到：%s", string(data), string(b)))
	}
	return index
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
