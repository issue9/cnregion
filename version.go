// SPDX-License-Identifier: MIT

package cnregion

import (
	"io/fs"

	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
)

// Version 用于描述与特定版本相关的区域数据
type Version struct {
	db        *db.DB
	provinces []Region
	districts []Region
}

// New 返回 Version 实例
func New(db *db.DB) (*Version, error) { return &Version{db: db}, nil }

// LoadFS 从 path 加载数据并初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func LoadFS(f fs.FS, path, separator string, version ...int) (*Version, error) {
	d, err := db.LoadFS(f, path, separator, true, version...)
	if err != nil {
		return nil, err
	}
	return New(d)
}

// Load 加载 data 数据初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func Load(data []byte, separator string, version ...int) (*Version, error) {
	d, err := db.Load(data, separator, true, version...)
	if err != nil {
		return nil, err
	}
	return New(d)
}

// LoadFile 从 path 加载数据并初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func LoadFile(path, separator string, version ...int) (*Version, error) {
	d, err := db.LoadFile(path, separator, true, version...)
	if err != nil {
		return nil, err
	}
	return New(d)
}

// SearchOptions 为搜索功能提供的参数
type SearchOptions = db.Options

// Search 简单的搜索功能
func (v *Version) Search(opt *SearchOptions) []Region {
	list := v.db.Search(opt)
	rs := make([]Region, 0, len(list))
	for _, item := range list {
		rs = append(rs, &dbRegion{r: item, v: v})
	}
	return rs
}

// Provinces 所有的顶级行政区域
func (v *Version) Provinces() []Region {
	if len(v.provinces) == 0 {
		root := v.db.Find()
		v.provinces = (&dbRegion{r: root, v: v}).Items()
	}
	return v.provinces
}

// Districts 按行政大区划分
//
// NOTE: 大区划分并不统一，按照各个省份的第一个数字进行划分。
func (v *Version) Districts() []Region {
	if len(v.districts) == 0 {
		dMap := make(map[byte]*districtRegion, len(districtsMap))
		provinces := v.Provinces()

		for index, name := range districtsMap {
			dMap[index] = &districtRegion{
				v:        v,
				id:       string(index),
				name:     name,
				fullName: name,
				fullID:   id.Fill(string(index), id.Village),
			}

			for _, p := range provinces {
				if p.ID()[0] == index {
					dMap[index].items = append(dMap[index].items, p)
				}
			}
		}

		v.districts = make([]Region, 0, len(dMap))
		for _, item := range dMap {
			v.districts = append(v.districts, item)
		}
	}

	return v.districts
}

var districtsMap = map[byte]string{
	'1': "华北地区",
	'2': "东北地区",
	'3': "华东地区",
	'4': "中南地区",
	'5': "西南地区",
	'6': "西北地区",
}
