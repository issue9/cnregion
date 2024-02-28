// SPDX-FileCopyrightText: 2021-2024 caixw
//
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

	AllLevel = Village + Town + County + City + Province
)

var lengths = map[Level]int{
	Village:  12,
	Town:     9,
	County:   6,
	City:     4,
	Province: 2,
}

// Length 获取各个类型 ID 的有效果长度
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

// SplitFilter 将 id 按区域进行划分且过滤掉零值的区域
//
//	330312123000 => 33 03 12 123
//
// 如果传递的是零值，则返回空数组。
func SplitFilter(id string) []string {
	province, city, county, town, village := Split(id)
	return filterZero(province, city, county, town, village)
}

func filterZero(id ...string) []string {
	for index, i := range id { // 过滤掉数组中的零值
		if isZero(i) {
			id = id[:index]
			break
		}
	}
	return id
}

// Parent 获取 id 的上一级行政区域的 ID
//
//	330312123456 => 330312123
func Parent(id string) string {
	list := SplitFilter(id)
	return strings.Join(list[:len(list)-1], "")
}

// Prefix 获取 ID 的非零前缀
//
//	330312123456 => 330312123456
//	330312123000 => 330312123
func Prefix(id string) string {
	return strings.Join(SplitFilter(id), "")
}

// Fill 为 id 填充后缀的 0
//
// id 为原始值；
// level 为需要达到的行政级别，最终的长度为 Length(level)。
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

// isZero 判断一组字符串是否都由 0 组成
func isZero(id string) bool {
	for _, r := range id {
		if r != '0' {
			return false
		}
	}
	return true
}
