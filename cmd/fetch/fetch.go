// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/issue9/cnregion/version"
	"github.com/issue9/term/v3/colors"
)

const baseURL = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"

// 拉取指定年份的数据
//
// years 为指定的一个或多个年份，如果为空，则表示所有的年份。
// 年份时间为 http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/
// 上存在的时间，从 2009 开始，到当前年份的上一年。
func fetch(dir string, interval time.Duration, years ...int) error {
	if len(years) == 0 {
		years = version.All()
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("拉取以下年份：%v\n", years)
	for _, year := range years {
		if err := fetchYear(dir, interval, year); err != nil {
			return err
		}
	}
	return nil
}

func fetchYear(dir string, interval time.Duration, year int) error {
	if !version.IsValid(year) {
		return version.ErrInvalidYear
	}

	y := strconv.Itoa(year)

	dir = filepath.Join(dir, y) // 带年份的目录
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	base := baseURL + y + "/" // 带年份地址的  URL
	provinces, err := collectProvinces(dir, base)
	if err != nil {
		return err
	}

	for _, province := range provinces {
		if province.ignore {
			colors.Println(colors.Normal, colors.Green, colors.Default, province.text, "\t已完成")
			continue
		}

		fs := newProvince(dir, province)
		if !fs.collect(base) {
			break
		}

		time.Sleep(interval)
	}

	return nil
}
