// SPDX-License-Identifier: MIT

// Package version 提供版本的相关信息
//
// 依据 http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/ 提供的数据，
// 以年作为单位进行更新，同时也以四位的年份作为版本号。
package version

import (
	"errors"
	"time"
)

// ErrInvalidYear 无效的年份版本
//
// 年份只能介于 [2009, 当前) 的区间之间。
var ErrInvalidYear = errors.New("无效的年份")

// Start 起始版本号
//
// 即提供的数据的起始年份。
const Start = 2009

// Last 最新的版本号
var Last = time.Now().Year() - 1

// All 返回支持的版本号列表
func All() []int {
	years := make([]int, 0, Last-Start)
	for year := Last; year >= Start; year-- {
		years = append(years, year)
	}

	return years
}
