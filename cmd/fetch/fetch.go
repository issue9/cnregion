// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"io"
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

const baseURL = "https://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/"

var digit = regexp.MustCompile("[0-9]+")

var errNoData = errors.New("no data")

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
		href := strings.TrimSuffix(e.Attr("href"), ".html")
		ignore := exists(filepath.Join(dir, href+".txt"))
		provinces = append(provinces, &item{
			href:   href + ".html",
			text:   e.Text,
			ignore: ignore,
			id:     href + strings.Repeat("0", 10),
		})

		state := colorsSprint(colors.Red, "\t未完成")
		if ignore {
			state = colorsSprint(colors.Green, "\t已完成")
		}
		fmt.Println(href, e.Text, state)
	})

	if err := c.Visit(base); err != nil {
		return err
	}
	c.Wait()
	if len(provinces) == 0 {
		return fmt.Errorf("未获取到 %s 年的省级数据", y)
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %d 年份的省级数据完成，总共 %d 条\n", year, len(provinces)))

	f, err := os.Create(dir + "/../" + y + "-error.log")
	if err != nil {
		return err
	}
	f.WriteString("此文件记录错误信息\n")

	for _, province := range provinces {
		if province.ignore {
			fmt.Println(colorsSprint(colors.Green, province.text, "\t已完成"))
			continue
		}

		if err := fetchProvince(f, dir, base, province); err != nil {
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
func fetchProvince(f io.Writer, dir, base string, p *item) error {
	fs := newProvinceFile(filepath.Join(dir, strings.TrimSuffix(p.href, ".html")+".txt"))
	fs.append(p.text, p.id) // 加入省级标记

	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	cities := make([]*item, 0, 500)
	c.OnHTML(".citytable .citytr", func(e *colly.HTMLElement) {
		cities = append(cities, getItem(e))
	})

	if err := c.Visit(base + p.href); err != nil {
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

		fs.append(city.text, city.id)
		if city.href == "" {
			continue
		}
		err = fetchCity(f, fs, base, city)
		switch {
		case errors.Is(err, errNoData):
			if err1 := fetchCounty(f, fs, base, city); err1 != nil { // 广东省 东莞
				if errors.Is(err1, errNoData) {
					err1 = fmt.Errorf("未获取到 %s:%s 的县/乡镇数据", city.id, city.text)
				}
				return err1
			}
		case err != nil:
			return err
		}
	}

	return fs.dump()
}

func fetchCity(f io.Writer, fs *provinceFile, base string, p *item) error {
	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	counties := make([]*item, 0, 500)
	c.OnHTML(".countytable .countytr", func(e *colly.HTMLElement) {
		counties = append(counties, getItem(e))
	})

	if err := c.Visit(base + p.href); err != nil {
		return err
	}
	c.Wait()

	if len(counties) == 0 {
		return errNoData
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %s 的县级数据完成，总共 %d 条\n", p.text, len(counties)))

	for _, county := range counties {
		if digit.MatchString(county.text) {
			continue
		}

		fs.append(county.text, county.id)
		if county.href == "" {
			continue
		}
		if err := fetchCounty(f, fs, base+firstID(p.href)+"/", county); err != nil {
			return err
		}
	}
	return nil
}

func fetchCounty(f io.Writer, fs *provinceFile, base string, p *item) error {
	c, err := buildCollector(base)
	if err != nil {
		return err
	}

	towns := make([]*item, 0, 500)
	c.OnHTML(".towntable .towntr", func(e *colly.HTMLElement) {
		towns = append(towns, getItem(e))
	})

	c.OnHTML(".countytable .towntr", func(e *colly.HTMLElement) { // 2021 之后的东莞等
		towns = append(towns, getItem(e))
	})

	// http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2014/46/4602.html
	c.OnHTML(".countytable .countytr", func(e *colly.HTMLElement) {
		towns = append(towns, getItem(e))
	})

	if err := c.Visit(base + p.href); err != nil {
		return err
	}
	c.Wait()

	if len(towns) == 0 {
		// 2014 460201
		io.WriteString(f, fmt.Sprintf("%s 返回乡镇数据为空，请确认该内容是否正常\n\n", base+p.href))
		return nil
	}
	fmt.Println(colorsSprintf(colors.Green, "拉取 %s 的乡镇数据完成，总共 %d 条\n", p.text, len(towns)))

	for _, town := range towns {
		if digit.MatchString(town.text) {
			continue
		}

		fs.append(town.text, town.id)
		if town.href == "" {
			continue
		}
		if err := fetchTown(f, fs, base+firstID(p.href)+"/", town); err != nil {
			return err
		}
	}
	return nil
}

func fetchTown(f io.Writer, fs *provinceFile, base string, p *item) error {
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
		fs.append(text, id)
	})

	if err := c.Visit(base + p.href); err != nil {
		return err
	}
	c.Wait()

	if count == 0 {
		// 街道可以为空，比如：
		// http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2015/34/01/11/340111009.html
		io.WriteString(f, fmt.Sprintf("%s 返回空数据，请确认该内容是否正常\n\n", base+p.href))
		return nil
	}
	fmt.Print(colorsSprintf(colors.Green, "拉取 %s 的街道数据完成，总共 %d 条\n", p.text, count))
	return nil
}

func getItem(e *colly.HTMLElement) *item {
	p := &item{}
	e.ForEach("td", func(i int, elem *colly.HTMLElement) {
		if i == 0 {
			elem.ForEach("a", func(j int, child *colly.HTMLElement) {
				if j > 1 {
					return
				}

				p.href = child.Attr("href")
				p.id = child.Text
			})

			if p.id == "" {
				p.href = ""
				p.id = elem.Text
			}
		} else if i == 1 {
			elem.ForEach("a", func(j int, child *colly.HTMLElement) {
				if j > 1 {
					return
				}

				p.text = child.Text
			})
			if p.text == "" {
				p.text = elem.Text
			}
		}
	})

	return p
}
