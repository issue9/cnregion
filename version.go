// SPDX-License-Identifier: MIT

package cnregion

import (
	_ "embed" // dbData
	"fmt"
	"io/fs"

	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
	"github.com/issue9/cnregion/version"
)

//go:embed data/regions.db
var dbData []byte

// Version 用于描述与特定版本相关的区域数据
type Version struct {
	versions  []int
	db        *db.DB
	provinces []Region
	districts []Region
}

// New 返回 Version 实例
//
// versions 表示需要加载的数据版本，即四位数的年份，可以同时指定多个版本。
// 有关数据版本的具体说明，可以参考 github.com/issue9/cnregion/version 包中的相关说明。
func New(db *db.DB, versions ...int) (*Version, error) {
	if len(versions) == 0 {
		versions = version.All()
	}

	for _, v := range versions {
		if -1 == db.VersionIndex(v) {
			return nil, fmt.Errorf("版本号 %d 并不存在于 db", v)
		}
	}

	return &Version{
		versions: versions,
		db:       db,
	}, nil
}

// Embed 将 data/regions.db 的内容嵌入到程序中
//
// 这样可以让程序不依赖外部文件，但同时也会增加编译后程序的大小。
// data/regions.db 目前大小为 7M 左右。
func Embed(separator string, version ...int) (*Version, error) {
	return Load(dbData, separator, version...)
}

// LoadFS 从 path 加载数据并初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func LoadFS(f fs.FS, path, separator string, version ...int) (*Version, error) {
	d, err := db.LoadFS(f, path, separator, true)
	if err != nil {
		return nil, err
	}

	return New(d, version...)
}

// Load 加载 data 数据初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func Load(data []byte, separator string, version ...int) (*Version, error) {
	d, err := db.Load(data, separator, true)
	if err != nil {
		return nil, err
	}

	return New(d, version...)
}

// LoadFile 从 path 加载数据并初始化 Version 实例
//
// separator 表示地名全名显示中，上下级之间的分隔符，比如"浙江-温州"，可以为空。
func LoadFile(path, separator string, version ...int) (*Version, error) {
	d, err := db.LoadFile(path, separator, true)
	if err != nil {
		return nil, err
	}

	return New(d, version...)
}

func (v *Version) isSupported(r *db.Region) bool {
	for _, vv := range v.versions {
		if r.IsSupported(vv) {
			return true
		}
	}
	return false
}

// SearchOptions 为搜索功能提供的参数
type SearchOptions = db.Options

// Search 简单的搜索功能
//
// text 表示你需要搜索的地名，不能是多个名称的组合，比如浙江温州，
// 直接搜温州就可以。也不要提供类似于居委会这种无实际意义的地名；
func (v *Version) Search(text string, opt *SearchOptions) []Region {
	list := v.db.Search(text, opt)
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

		for k, v := range districtsMap {
			dMap[k] = &districtRegion{
				id:       string(k),
				name:     v,
				fullName: v,
				fullID:   id.Fill(string(k), id.Village),
			}

			for _, p := range provinces {
				if p.ID()[0] == k {
					dMap[k].items = append(dMap[k].items, p)
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
