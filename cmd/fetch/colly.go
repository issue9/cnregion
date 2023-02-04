// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/issue9/errwrap"
	"github.com/issue9/sliceutil"
	"github.com/issue9/term/v3/colors"

	"github.com/issue9/cnregion/id"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"

type province struct {
	lock  *sync.Mutex
	items []*item
	dir   string
	id    string
	text  string
}

type item struct {
	id, text string
	ignore   bool // 忽略此条数据
}

func newProvince(dir string, i *item) *province {
	return &province{
		lock:  &sync.Mutex{},
		items: make([]*item, 0, 5000),
		dir:   dir,
		id:    i.id,
		text:  i.text,
	}
}

func (fs *province) append(id, text string) {
	id = trimID(id)
	if id == text {
		return
	}

	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.items = append(fs.items, &item{id: id, text: text})
}

func (fs *province) dump() (ok bool) {
	fs.append(fs.id, fs.text) // 加入省级标记

	fs.items = sliceutil.Unique(fs.items, func(i, j *item) bool { return i.id == j.id })

	sort.SliceStable(fs.items, func(i, j int) bool { return fs.items[i].id < fs.items[j].id })

	path := filepath.Join(fs.dir, fs.id+".txt")
	colors.Printf(colors.Normal, colors.Green, colors.Default, "完成收集 %s 的数据 %d 条，将写入到 %s\n", fs.text, len(fs.items), path)

	buf := errwrap.Buffer{}
	for _, item := range fs.items {
		buf.Printf("%s\t%s\n", item.id, item.text)
	}
	if buf.Err != nil {
		colors.Println(colors.Normal, colors.Red, colors.Default, buf.Err)
		return false
	}

	if err := os.WriteFile(path, buf.Bytes(), os.ModePerm); err != nil {
		colors.Println(colors.Normal, colors.Red, colors.Default, err)
	} else {
		colors.Printf(colors.Normal, colors.Green, colors.Default, "写入 %s 完成\n", path)
	}
	return true
}

// depth 是否访问子路径
func buildColly(base string) (*colly.Collector, error) {
	expr := base + "[0-9/]+.html"
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(base),
			regexp.MustCompile(expr),
		),
		colly.UserAgent(userAgent),
		colly.DetectCharset(),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	rule := &colly.LimitRule{DomainGlob: "*", Parallelism: 50, RandomDelay: 10 * time.Second}
	if err := c.Limit(rule); err != nil {
		return nil, err
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("抓取 %s\n", r.URL)
	})

	c.OnError(func(resp *colly.Response, err error) {
		fmt.Printf("ERROR: %s 并返回状态码 %d\n", err, resp.StatusCode)

		// 重试
		if err := c.Visit(resp.Request.URL.String()); err != nil {
			fmt.Printf("ERROR:%s at visit %s", err, resp.Request.URL.String())
		}
	})

	return c, nil
}

var digit = regexp.MustCompile("[0-9]+")

// base 需要以 / 结尾
func (fs *province) collect(base string) (ok bool) {
	fmt.Printf("开始收集 %s 的数据\n", fs.text)

	base = base + fs.id
	c, err := buildColly(base)
	if err != nil {
		colors.Println(colors.Normal, colors.Red, colors.Default, err)
		return false
	}

	visit := func(e *colly.HTMLElement) {
		if digit.MatchString(e.Text) {
			return
		}

		if err := e.Request.Visit(e.Attr("href")); err != nil {
			fmt.Printf("ERROR: %s @ %s\n", err, e.Text)
		}
	}

	// 市
	c.OnHTML(".citytable .citytr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
		visit(e)
	})

	// 县
	c.OnHTML(".countytable .countytr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
		visit(e)
	})

	// 乡镇
	c.OnHTML(".towntable .towntr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
		visit(e)
	})

	// 街道、村庄
	c.OnHTML(".villagetable .villagetr", func(e *colly.HTMLElement) {
		var id, text string
		e.ForEach("td", func(i int, elem *colly.HTMLElement) {
			if i == 0 {
				id = elem.Text
			} else if i == 2 {
				text = elem.Text
			}
		})

		fs.append(id, text)
		visit(e)
	})

	if err := c.Visit(base + ".html"); err != nil {
		colors.Println(colors.Normal, colors.Red, colors.Default, err)
		return false
	}

	c.Wait()

	return fs.dump()
}

// 截取数字部分，不够填补后缀 0。
func trimID(regionID string) string {
	regionID = strings.TrimSuffix(regionID, ".html")
	index := strings.LastIndexByte(regionID, '/')
	if index >= 0 {
		regionID = regionID[index+1:]
	}

	l := len(regionID)
	if l < id.Length(id.Village) {
		regionID += strings.Repeat("0", id.Length(id.Village)-l)
	}

	return regionID
}

// 收集省份列表
//
// dir 保存地址，最后一个是年份
// base URL 基地址，需要以 / 结尾
func collectProvinces(dir string, base string) ([]*item, error) {
	c, err := buildColly(base)
	if err != nil {
		return nil, err
	}

	items := make([]*item, 0, 50)
	c.OnHTML(".provincetable .provincetr td a", func(e *colly.HTMLElement) {
		items = append(items, &item{
			id:   strings.TrimSuffix(e.Attr("href"), ".html"),
			text: e.Text,
		})
	})

	if err := c.Visit(base); err != nil {
		return nil, err
	}
	c.Wait()

	// 检测目录
	for _, i := range items {
		fmt.Print(i.id, "\t", i.text, "\t")

		i.ignore = exists(dir, i.id+".txt")
		if i.ignore {
			colors.Println(colors.Normal, colors.Green, colors.Default, "已完成")
		} else {
			colors.Println(colors.Normal, colors.Red, colors.Default, "未完成")
		}
	}

	return items, nil
}

func exists(dir string, id string) bool {
	_, err := os.Stat(filepath.Join(dir, id))
	return err == nil || !errors.Is(err, os.ErrNotExist)
}
