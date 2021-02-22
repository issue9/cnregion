// SPDX-License-Identifier: MIT

package cnregion

import (
	"fmt"

	"github.com/issue9/cnregion/db"
)

// Version 用于描述与特定版本相关的区域数据
type Version struct {
	version   int
	db        *db.DB
	provinces []Region
	districts []Region
}

// New 返回 Version 实例
//
// version 表示需要的数据版本，即四位数的年份。
func New(db *db.DB, version int) *Version {
	if -1 == db.VersionIndex(version) {
		panic(fmt.Sprintf("版本号 %d 并不存在于 db", version))
	}

	return &Version{
		version: version,
		db:      db,
	}
}

// Load 加载 data 数据初始化 Version 实例
func Load(data []byte, separator string, version int) (*Version, error) {
	d, err := db.Load(data, separator, true)
	if err != nil {
		return nil, err
	}

	return New(d, version), nil
}

// LoadFile 从 path 加载数据并初始化 Version 实例
func LoadFile(path, separator string, version int) (*Version, error) {
	d, err := db.LoadFile(path, separator, true)
	if err != nil {
		return nil, err
	}

	return New(d, version), nil
}

// Provinces 所有的顶级行政区域
func (v *Version) Provinces() []Region {
	if len(v.provinces) == 0 {
		root := v.db.Find()
		v.provinces = (&dbRegion{r: root}).Items()
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
