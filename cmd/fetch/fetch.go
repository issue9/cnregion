// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/issue9/cnregion/version"
	"github.com/issue9/term/v3/colors"
)

const baseURL = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"

var digit = regexp.MustCompile("[0-9]+")

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

	fmt.Printf("拉取以下年份：%v\n", colorsSprint(colors.Green, years))
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

	fmt.Printf("准备拉取 %s 的数据\n", colorsSprint(colors.Green, year))

	y := strconv.Itoa(year)
	dir = filepath.Join(dir, y) // 带年份的目录
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	base := baseURL + y + "/" // 带年份地址的  URL
	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	provinces := make([]*item, 0, 50)
	c.OnHTML(".provincetable .provincetr td a", func(e *colly.HTMLElement) {
		id := strings.TrimSuffix(e.Attr("href"), ".html")
		ignore := exists(filepath.Join(dir, id+".txt"))
		provinces = append(provinces, &item{
			id:     id,
			text:   e.Text,
			ignore: ignore,
		})

		state := colorsSprint(colors.Red, "\t未完成")
		if ignore {
			state = colorsSprint(colors.Green, "\t已完成")
		}
		fmt.Println(id, e.Text, state)
	})

	if err := c.Visit(base); err != nil {
		return err
	}
	c.Wait()
	fmt.Println(colorsSprintf(colors.Green, "拉取 %d 年份的数据完成，总共 %d 条\n", year, len(provinces)))

	f, err := os.Create(dir + "/../error.log")
	if err != nil {
		return err
	}
	f.WriteString("此文件记录错误信息\n")

	for _, province := range provinces {
		if province.ignore {
			fmt.Println(colorsSprint(colors.Green, province.text, "\t已完成"))
			continue
		}

		if err := fetchProvince(dir, base, province); err != nil {
			// 出错就忽略这个省份的输出，继续下一个省的。
			fmt.Println(colorsSprint(colors.Red, err))
			f.WriteString(y)
			f.WriteString("\t")
			f.WriteString(err.Error())
			f.WriteString("\n\n")
		}
		time.Sleep(interval)
	}

	return f.Close()
}

// base 格式： https://example.com/2022/ 到年份为止的数据
func fetchProvince(dir, base string, p *item) error {
	fs := newProvinceFile(filepath.Join(dir, p.id+".txt"))
	fs.append(p.id, p.text) // 加入省级标记

	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	cities := make([]*item, 0, 500)
	c.OnHTML(".citytable .citytr td a", func(e *colly.HTMLElement) {
		href := strings.TrimSuffix(e.Attr("href"), ".html")
		cities = append(cities, &item{id: href, text: e.Text})
	})

	if err := c.Visit(base + p.id + ".html"); err != nil {
		return err
	}
	c.Wait()

	if len(cities) == 0 {
		return fmt.Errorf("未获取到 %s:%s 的市级数据", p.id, p.text)
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %s 的市级数据完成，总共 %d 条\n", p.text, len(cities)))

	for _, city := range cities {
		if digit.MatchString(city.text) {
			continue
		}

		if err := fetchCity(fs, base, city); err != nil {
			return err
		}
	}

	return fs.dump()
}

func fetchCity(fs *provinceFile, base string, p *item) error {
	fs.append(p.id, p.text)

	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	counties := make([]*item, 0, 500)
	c.OnHTML(".countytable .countytr td a", func(e *colly.HTMLElement) {
		href := strings.TrimSuffix(e.Attr("href"), ".html")
		counties = append(counties, &item{id: href, text: e.Text})
	})

	if err := c.Visit(base + p.id + ".html"); err != nil {
		return err
	}
	c.Wait()

	if len(counties) == 0 {
		return fmt.Errorf("未获取到 %s:%s 的县级数据", p.id, p.text)
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %s 的县级数据完成，总共 %d 条\n", p.text, len(counties)))

	for _, county := range counties {
		if digit.MatchString(county.text) {
			continue
		}

		if err := fetchCounty(fs, base+firstID(p.id)+"/", county); err != nil {
			return err
		}
	}
	return nil
}

func fetchCounty(fs *provinceFile, base string, p *item) error {
	fs.append(p.id, p.text)

	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	towns := make([]*item, 0, 500)
	c.OnHTML(".towntable .towntr td a", func(e *colly.HTMLElement) {
		href := strings.TrimSuffix(e.Attr("href"), ".html")
		towns = append(towns, &item{id: href, text: e.Text})
	})

	if err := c.Visit(base + p.id + ".html"); err != nil {
		return err
	}
	c.Wait()

	if len(towns) == 0 {
		return fmt.Errorf("未获取到 %s:%s 的乡镇数据", p.id, p.text)
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %s 的乡镇数据完成，总共 %d 条\n", p.text, len(towns)))

	for _, town := range towns {
		if digit.MatchString(town.text) {
			continue
		}

		if err := fetchTown(fs, base+firstID(p.id)+"/", town); err != nil {
			return err
		}
	}
	return nil
}

func fetchTown(fs *provinceFile, base string, p *item) error {
	fs.append(p.id, p.text)

	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	var count int
	c.OnHTML(".villagetable .villagetr", func(e *colly.HTMLElement) {
		var id, text string
		e.ForEach("td", func(i int, elem *colly.HTMLElement) {
			if i == 0 {
				id = elem.Text
			} else if i == 2 {
				text = elem.Text
			}
		})
		count++
		fs.append(id, text)
	})

	if err := c.Visit(base + p.id + ".html"); err != nil {
		return err
	}
	c.Wait()

	if count == 0 {
		return fmt.Errorf("未获取到 %s:%s 的街道数据", p.id, p.text)
	}
	fmt.Print(colorsSprintf(colors.Green, "拉取 %s 的街道数据完成，总共 %d 条\n", p.text, count))
	return nil
}
