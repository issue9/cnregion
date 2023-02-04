// SPDX-License-Identifier: MIT

// Package version 提供版本的相关信息
//
// 依据 http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/ 提供的数据，
// 以年作为单位进行更新，同时也以四位的年份作为版本号。
package version

import "fmt"

// ErrInvalidYear 无效的年份版本
//
// 年份只能介于 [2009, 当前年份-1) 的区间之间。
var ErrInvalidYear = fmt.Errorf("无效的版本号，必须是介于 [%d,%d] 之间的整数", start, latest)

// 起始版本号，即提供的数据的起始年份。
const start = 2009

// 最新的有效年份，每次更新数据之后，需要手动更新此值。
var latest = 2022

// All 返回支持的版本号列表
func All() []int { return Range(start, latest) }

// IsValid 验证年份是否为一个有效的版本号
func IsValid(year int) bool { return year >= start && year <= latest }

// BeginWith 从 begin 开始直到最新年份
func BeginWith(begin int) []int {
	return Range(begin, latest)
}

// Range 获取指定范围内的版本号
func Range(begin, end int) []int {
	if !IsValid(begin) {
		panic(ErrInvalidYear)
	}

	if !IsValid(end) {
		panic(ErrInvalidYear)
	}

	if begin > end {
		panic(ErrInvalidYear)
	}

	years := make([]int, 0, end-begin+1)
	for year := end; year >= begin; year-- {
		years = append(years, year)
	}
	return years
}
