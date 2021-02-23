// SPDX-License-Identifier: MIT

package db

import (
	"strings"

	"github.com/issue9/cnregion/id"
)

// Options 搜索的附加选项
type Options struct {
	// 上一级的区域 ID
	Parent string

	// 搜索的城市类型，该值取值于 github.com/issue9/cnregion/id.Level 类型。
	// 多个值可以通过或运算叠加。
	Level id.Level

	// 最大的搜索数量。0 表示不限制数量。
	Max       int
	unlimited bool
}

// Search 简单的搜索功能
//
// text 表示你需要搜索的地名，不能是多个名称的组合，比如浙江温州，
// 直接搜温州就可以。也不要提供类似于居委会这种无实际意义的地名；
func (db *DB) Search(text string, opt *Options) []*Region {
	if opt == nil {
		opt = &Options{}
	}

	r := db.region
	if opt.Parent != "" {
		r = db.Find(id.SplitFilter(opt.Parent)...)
	}
	if r == nil { // 不存在 opt.Parent 指定的数据
		return nil
	}

	if opt.Level == 0 {
		opt.Level = id.AllLevel
	}

	opt.unlimited = opt.Max == 0
	size := 100
	if !opt.unlimited {
		size = opt.Max
	}
	list := make([]*Region, 0, size)

	return r.search(text, opt, list)
}

func (reg *Region) search(text string, opt *Options, list []*Region) []*Region {
	if strings.Contains(reg.Name, text) && (reg.level&opt.Level == reg.level) {
		list = append(list, reg)
		opt.Max--
	}

	if !opt.unlimited && opt.Max <= 0 {
		return list
	}

	for _, item := range reg.Items {
		list = item.search(text, opt, list)
	}

	return list
}
