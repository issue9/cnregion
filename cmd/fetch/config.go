// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"time"
)

// 年份只能介于 [2009, 当前) 的区间之间。
var errInvalidYear = errors.New("无效的年份")

const (
	baseURL   = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"
	startYear = 2009
)

var lastYear = time.Now().Year() - 1

func allYears() []int {
	years := make([]int, 0, lastYear-startYear)
	for year := lastYear; year >= startYear; year-- {
		years = append(years, year)
	}

	return years
}
