// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"

	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/version"
)

const baseURL = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"

// 拉取指定年份的数据
//
// years 为指定的一个或多个年份，如果为空，则表示所有的年份。
// 年份时间为 http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/
// 上存在的时间，从 2009 开始，到当前年份的上一年。
func fetch(dir string, years ...int) error {
	if len(years) == 0 {
		years = version.All()
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, year := range years {
		if err := fetchYear(dir, year); err != nil {
			return err
		}
	}
	return nil
}

func fetchYear(dir string, year int) error {
	if year < version.Start || year > version.Last {
		return version.ErrInvalidYear
	}

	buf := &errwrap.Buffer{Buffer: bytes.Buffer{}}

	y := strconv.Itoa(year)

	dir = filepath.Join(dir, y)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return collect(dir, buf, baseURL+y)
}
