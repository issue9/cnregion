// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/issue9/errwrap"
	"github.com/issue9/sliceutil"
	"github.com/issue9/term/v3/colors"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"

// 以省为单位的文件内容管理
type provinceFile struct {
	lock  *sync.Mutex
	items []*item
	path  string
}

type item struct {
	id     string // 区域 ID
	href   string // href 的属性值，仅存在于中间过程
	text   string
	ignore bool // 忽略此条数据
}

func newProvinceFile(path string) *provinceFile {
	return &provinceFile{
		path:  path,
		lock:  &sync.Mutex{},
		items: make([]*item, 0, 50000),
	}
}

func (fs *provinceFile) append(text, id string) {
	if id == text {
		return
	}

	if id == "" || text == "" {
		panic(fmt.Sprintf("数据不能为空%s:%s", id, text))
	}

	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.items = append(fs.items, &item{text: text, id: id})
}

func (fs *provinceFile) dump() error {
	fmt.Printf("去重(%d)...\n", len(fs.items))
	fs.items = sliceutil.Unique(fs.items, func(i, j *item) bool { return i.id == j.id })

	fmt.Printf("排序(%d)...\n", len(fs.items))
	sort.SliceStable(fs.items, func(i, j int) bool { return fs.items[i].id < fs.items[j].id })

	fmt.Println(colorsSprintf(colors.Green, "准备将 %d 条数据写入 %s\n", len(fs.items), fs.path))

	buf := errwrap.Buffer{}
	for _, item := range fs.items {
		buf.Printf("%s\t%s\n", item.id, item.text)
	}
	if buf.Err != nil {
		return buf.Err
	}

	if err := os.WriteFile(fs.path, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}

	colors.Printf(colors.Normal, colors.Green, colors.Default, "写入 %s 完成\n\n", fs.path)
	return nil
}

func buildCollector(base string) (*colly.Collector, error) {
	expr := base + "[0-9/]*.html"
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(base),
			regexp.MustCompile(expr),
		),
		colly.UserAgent(userAgent),
		colly.DetectCharset(),
		colly.AllowURLRevisit(),
		colly.CacheDir("./caches"),
	)

	rule := &colly.LimitRule{Parallelism: 100, DomainGlob: "*", Delay: time.Second}
	if err := c.Limit(rule); err != nil {
		return nil, err
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("抓取 %s\n", r.URL)
	})

	c.OnError(func(resp *colly.Response, err error) {
		colors.Printf(colors.Normal, colors.Red, colors.Default, "ERROR: %s 并返回状态码 %d\n", err, resp.StatusCode)

		// 重试
		if err := c.Visit(resp.Request.URL.String()); err != nil {
			colors.Printf(colors.Normal, colors.Red, colors.Default, "ERROR: %s at visit %s\n", err, resp.Request.URL.String())
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if len(r.Body) == 0 {
			colors.Printf(colors.Normal, colors.Red, colors.Default, "页面 %s 没有数据\n", r.Request.URL.String())
		}
	})

	return c, nil
}

func firstID(href string) string {
	href = strings.TrimSuffix(href, ".html")
	index := strings.IndexByte(href, '/')
	if index <= 0 {
		return strings.TrimSuffix(href, ".html")
	}
	return href[:index]
}
