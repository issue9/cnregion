// SPDX-License-Identifier: MIT

// Package id 针对 ID 的一些操作函数
package id

import (
	"fmt"
	"strings"
)

// Level 表示区域的级别
type Level uint8

// 对区域级别的定义
const (
	Village Level = 1 << iota
	Town
	County
	City
	Province
)

var lengths = map[Level]int{
	Village:  12,
	Town:     9,
	County:   6,
	City:     4,
	Province: 2,
}

// Length 获取各个类型 ID 的实际有效果长度
func Length(level Level) int {
	if _, found := lengths[level]; !found {
		panic("无效的 level 参数")
	}

	return lengths[level]
}

// Split 将一个区域 ID 按区域进行划分
func Split(id string) (province, city, county, town, village string) {
	if len(id) != Length(Village) {
		panic(fmt.Sprintf("id 的长度只能为 %d，当前为 %s", Length(Village), id))
	}

	return id[:Length(Province)],
		id[Length(Province):Length(City)],
		id[Length(City):Length(County)],
		id[Length(County):Length(Town)],
		id[Length(Town):Length(Village)]
}

// Fill 为 id 填充后缀的 0
func Fill(id string, level Level) string {
	rem := Length(level) - len(id)
	switch {
	case rem == 0:
		return id
	case rem > Length(level) || rem < 2:
		panic(fmt.Sprintf("无效的 id %s，无法为其填充 0", id))
	default:
		return id + strings.Repeat("0", rem)
	}
}

// IsZero 判断一组字符串是否都由 0 组成
func IsZero(id string) bool {
	for _, r := range id {
		if r != '0' {
			return false
		}
	}
	return true
}
