// SPDX-License-Identifier: MIT

package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/issue9/errwrap"
)

type files map[string][]*item

type item struct {
	id, text string
}

func (fs files) append(id, text string) {
	id = trimID(id)
	if id == text {
		return
	}

	path := id[:2]
	if _, found := fs[path]; !found {
		fs[path] = make([]*item, 0, 100)
	}
	fs[path] = append(fs[path], &item{id: id, text: text})
}

func (fs files) dump(dir string) error {
	for id, items := range fs {
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

func collect(dir string, buf *errwrap.Buffer, base string) error {
	expr := base + "/[0-9]*.html"
	c := colly.NewCollector(colly.URLFilters(
		regexp.MustCompile(base),
		regexp.MustCompile(expr),
	), colly.DetectCharset())

	fs := files{}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
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

	if err := c.Visit(base); err != nil {
		return err
	}

	c.Wait()

	return fs.dump(dir)
}

func trimID(id string) string {
	id = strings.TrimSuffix(id, ".html")
	index := strings.IndexByte(id, '/')
	if index >= 0 {
		id = id[index+1:]
	}

	l := len(id)
	if l < 12 {
		id += strings.Repeat("0", 12-l)
	}

	return id
}
