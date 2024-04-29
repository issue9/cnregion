// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"strings"

	"github.com/issue9/cnregion/id"
)

// Options 搜索选项
type Options struct {
	// 表示你需要搜索的地名需要包含的内容
	//
	// 不能是多个名称的组合，比如"浙江温州"，直接写"温州"就可以。
	// 也不要提供类似于"居委会"这种无实际意义的地名；
	Text string

	// 上一级的区域 ID
	//
	// 为空表示不限制。
	Parent string

	// 搜索的城市类型
	//
	// 该值取值于 github.com/issue9/cnregion/id.Level 类型。 多个值可以通过或运算叠加。
	// 0 表示所有类型。
	Level id.Level

	// 最大的搜索数量。0 表示不限制数量。
	Max       int
	unlimited bool
}

func (o *Options) isEmpty() bool {
	return o.Text == "" &&
		(o.Parent == "" || o.Parent == "000000000000") &&
		o.Level == 0 &&
		o.Max == 0
}

// Search 简单的搜索功能
func (db *DB) Search(opt *Options) []*Region {
	if opt == nil || opt.isEmpty() {
		panic("参数 opt 不能为空值")
	}

	r := db.root
	if opt.Parent != "" {
		r = db.Find(opt.Parent)
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

	return r.search(opt, list)
}

func (reg *Region) search(opt *Options, list []*Region) []*Region {
	if strings.Contains(reg.name, opt.Text) &&
		(reg.level&opt.Level == reg.level) && reg.level != 0 { // level == 0 只有根元素才有
		list = append(list, reg)
		opt.Max--
	}

	if !opt.unlimited && opt.Max <= 0 {
		return list
	}

	for _, item := range reg.items {
		list = item.search(opt, list)
	}

	return list
}
