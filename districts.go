// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import "github.com/issue9/cnregion/v2/id"

// Districts 按行政大区划分
//
// NOTE: 大区划分并不统一，按照各个省份的第一个数字进行划分。
func (db *DB) Districts() []*Region { return db.districts }

func (db *DB) initDistricts() {
	db.districts = make([]*Region, 0, len(districtsMap))

	for index, name := range districtsMap {
		items := make([]*Region, 0, 10)
		for _, p := range db.Provinces() {
			if p.ID()[0] == index {
				items = append(items, p)
			}
		}

		db.districts = append(db.districts, &Region{
			id:       string(index),
			fullID:   id.Fill(string(index), id.Village),
			name:     name,
			fullName: name,
			items:    items,
		})
	}

}

var districtsMap = map[byte]string{
	'1': "华北地区",
	'2': "东北地区",
	'3': "华东地区",
	'4': "中南地区",
	'5': "西南地区",
	'6': "西北地区",
}
