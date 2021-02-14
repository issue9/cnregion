// SPDX-License-Identifier: MIT

// Package id 针对 ID 的一些操作函数
package id

import "fmt"

// Level 表示一个 ID 表示的区域级别
type Level int8

// 预定义的区域级别
const (
	Village Level = iota
	Town
	County
	City
	Province
)

// ID 用于描述区域 ID
type ID struct {
	Level                                 Level
	Province, City, County, Town, Village string
}

// ParseID 用于将区域 ID 解析到 ID 对象
func ParseID(id string) *ID {
	if len(id) != 12 {
		panic(fmt.Sprintf("id 的长度只能为 12，当前为 %s", id))
	}

	ret := &ID{
		Province: id[:2],
		City:     id[2:4],
		County:   id[4:6],
		Town:     id[6:9],
		Village:  id[9:12],
	}

	switch {
	case ret.Province == "00":
		panic(fmt.Sprintf("省份不能为空：%s", id))
	case ret.City == "00":
		ret.Level = Province
	case ret.County == "00":
		ret.Level = City
	case ret.Town == "000":
		ret.Level = County
	case ret.Village == "000":
		ret.Level = Town
	}

	return ret
}

// Split 将一个区域 ID 按区域进行划分
func Split(id string) (province, city, county, town, village string) {
	if len(id) != 12 {
		panic(fmt.Sprintf("id 的长度只能为 12，当前为 %s", id))
	}

	return id[:2], id[2:4], id[4:6], id[6:9], id[9:12]
}
