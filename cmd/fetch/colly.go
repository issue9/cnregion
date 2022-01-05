// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/issue9/errwrap"

	"github.com/issue9/cnregion/id"
)

type files struct {
	lock  *sync.Mutex
	items map[string][]*item
}

type item struct {
	id, text string
}

func newFiles() *files {
	return &files{
		lock:  &sync.Mutex{},
		items: make(map[string][]*item, 40),
	}
}

func (fs files) append(id, text string) {
	id = trimID(id)
	if id == text {
		return
	}

	fs.lock.Lock()
	defer fs.lock.Unlock()

	path := id[:2]
	if _, found := fs.items[path]; !found {
		fs.items[path] = make([]*item, 0, 100)
	}
	fs.items[path] = append(fs.items[path], &item{id: id, text: text})
}

func (fs files) dump(dir string) error {
	for id, items := range fs.items {
		sort.SliceStable(items, func(i, j int) bool { return items[i].id < items[j].id })

		buf := errwrap.Buffer{Buffer: bytes.Buffer{}}
		for _, item := range items {
			buf.Printf("%s\t%s\n", item.id, item.text)
		}
		if buf.Err != nil {
			return buf.Err
		}

		path := filepath.Join(dir, id+".txt")
		if err := ioutil.WriteFile(path, buf.Bytes(), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func collect(dir string, base string) error {
	expr := base + "/[0-9]+.html"
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(base),
			regexp.MustCompile(expr),
		),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
		colly.DetectCharset(),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	rule := &colly.LimitRule{DomainGlob: "*", Parallelism: 50, RandomDelay: 10 * time.Second}
	if err := c.Limit(rule); err != nil {
		return err
	}

	fs := newFiles()

	digit := regexp.MustCompile("[0-9]+")
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if digit.MatchString(e.Text) {
			return
		}

		if err := e.Request.Visit(e.Attr("href")); err != nil {
			fmt.Printf("ERROR: %s @ %s\n", err, e.Text)
		}
	})

	// 省
	c.OnHTML(".provincetable .provincetr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
	})

	// 市
	c.OnHTML(".citytable .citytr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
	})

	// 县
	c.OnHTML(".countytable .countytr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
	})

	// 乡镇
	c.OnHTML(".towntable .towntr td a", func(e *colly.HTMLElement) {
		fs.append(e.Attr("href"), e.Text)
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
		id = trimID(id)
		fs.append(id, text)
	})

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

	if err := c.Visit(base); err != nil {
		return err
	}

	c.Wait()

	return fs.dump(dir)
}

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
