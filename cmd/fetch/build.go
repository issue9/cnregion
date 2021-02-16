// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
	"github.com/issue9/cnregion/version"
)

func build(dataDir string, output string, years ...int) error {
	if len(years) == 0 {
		years = version.All()
	}

	var d *db.DB

	if fileExists(output) {
		dd, err := db.Load(output)
		if err != nil {
			return err
		}
		d = dd
	} else {
		d = &db.DB{Region: &db.Region{}}
	}

	for _, year := range years {
		if err := buildYear(d, dataDir, year); err != nil {
			return err
		}
	}

	return d.Dump(output)
}

func buildYear(d *db.DB, dataDir string, year int) error {
	for _, v := range d.Versions { // 检测 year 是否已经存在？
		if v == year {
			fmt.Printf("年份 %d 已经存在，忽略该数据\n", year)
			return nil
		}
	}

	fmt.Printf("添加 %d 的数据\n", year)

	d.Versions = append(d.Versions, year)

	y := strconv.Itoa(year)
	dataDir = filepath.Join(dataDir, y)

	return filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		s := bufio.NewScanner(bytes.NewBuffer(data))
		s.Split(bufio.ScanLines)
		for s.Scan() {
			txt := s.Text()
			vals := strings.Split(txt, "\t")
			if len(vals) != 2 {
				return fmt.Errorf("无效的格式，位于 %s:%s", path, txt)
			}
			id, name := vals[0], vals[1]

			if err := appendDB(d, year, id, name); err != nil {
				return err
			}
		}

		return nil
	})
}

func appendDB(d *db.DB, year int, regionID, name string) error {
	province, city, county, town, village := id.Split(regionID)
	list := filterZero(province, city, county, town, village)
	item := d.Find(list...)

	if item == nil {
		item = d.Find(list[:len(list)-1]...) // 上一级
		return item.AddItem(d, list[len(list)-1], name, year)
	}

	return item.SetSupported(d, year)
}

func filterZero(regionID ...string) []string {
	for index, i := range regionID { // 过滤掉数组中的零值
		if id.IsZero(i) {
			regionID = regionID[:index]
			break
		}
	}
	return regionID
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
