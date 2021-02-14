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
	regionid "github.com/issue9/cnregion/id"
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

func appendDB(d *db.DB, year int, id, name string) error {
	province, city, county, town, village := regionid.Split(id)

	switch {
	case province == "00":
		panic(fmt.Sprintf("无效的 ID %s，省份不能为 00", id))
	case city == "00": // 添加省份 ID
		d.Items = append(d.Items, &db.Region{
			ID:   province,
			Name: name,
		})
	case county == "00": // 添加市
		item := findItem(d.Items, province)
		item.Items = append(item.Items, &db.Region{
			ID:   city,
			Name: name,
		})
	case town == "000": // 添加县
		item := findItem(d.Items, province)
		item = findItem(item.Items, city)
		item.Items = append(item.Items, &db.Region{
			ID:   county,
			Name: name,
		})
	case village == "000": // 添加乡
		item := findItem(d.Items, province)
		item = findItem(item.Items, city)
		item = findItem(item.Items, county)
		item.Items = append(item.Items, &db.Region{
			ID:   town,
			Name: name,
		})
	default:
		item := findItem(d.Items, province)
		item = findItem(item.Items, city)
		item = findItem(item.Items, county)
		item = findItem(item.Items, town)
		item.Items = append(item.Items, &db.Region{
			ID:   village,
			Name: name,
		})
	}

	return nil
}

func findItem(items []*db.Region, id string) *db.Region {
	for _, item := range items {
		if item.ID == id {
			return item
		}
	}
	panic(fmt.Sprintf("未在当前列表中找到 %s", id))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
