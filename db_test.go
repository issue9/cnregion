// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cnregion/v2/id"
	"github.com/issue9/cnregion/v2/version"
)

var data = []byte(`1:[2020,2019]:::1:2{33:浙江:1:1{01:温州:3:0{}}34:安徽:1:3{01:合肥:3:0{}02:芜湖:1:0{}03:芜湖-2:1:0{}}}`)

var obj = &DB{
	versions:          []int{2020, 2019},
	fullNameSeparator: "-",
	root: &Region{
		name:     "",
		versions: []int{2020},
		items: []*Region{
			{
				id:       "33",
				name:     "浙江",
				versions: []int{2020},
				fullName: "浙江",
				fullID:   "330000000000",
				level:    id.Province,
				items: []*Region{
					{
						id:       "01",
						name:     "温州",
						versions: []int{2020, 2019},
						fullName: "浙江-温州",
						fullID:   "330100000000",
						level:    id.City,
					},
				},
			},
			{
				id:       "34",
				name:     "安徽",
				fullName: "安徽",
				fullID:   "340000000000",
				versions: []int{2020},
				level:    id.Province,
				items: []*Region{
					{
						id:       "01",
						name:     "合肥",
						versions: []int{2020, 2019},
						fullName: "安徽-合肥",
						fullID:   "340100000000",
						level:    id.City,
					},
					{
						id:       "02",
						name:     "芜湖",
						versions: []int{2020},
						fullName: "安徽-芜湖",
						fullID:   "340200000000",
						level:    id.City,
					},
					{
						id:       "03",
						name:     "芜湖-2",
						versions: []int{2020},
						fullName: "安徽-芜湖-2",
						fullID:   "340300000000",
						level:    id.City,
					},
				},
			},
		},
	},
}

func init() {
	setRegionDB(obj.root, obj)
}

func setRegionDB(r *Region, db *DB) {
	r.db = db
	for _, i := range r.items {
		setRegionDB(i, db)
	}
}

func TestDB_Find(t *testing.T) {
	a := assert.New(t, false)

	// 2020
	db, err := LoadFile("./data/regions.db", ">", true, 2020)
	a.NotError(err).NotNil(db)
	r := db.Find("330305000000")
	a.NotNil(r).
		Equal(r.ID(), "05").
		Equal(r.FullID(), "330305000000").
		Equal(r.Name(), "洞头区").
		Equal(r.FullName(), "浙江省>温州市>洞头区").
		Equal(r.Versions(), []int{2020})
	r = db.Find("330322000000") // 洞头县，已改为洞头区
	a.Nil(r)

	// 2009
	db, err = LoadFile("./data/regions.db", ">", true, 2009)
	a.NotError(err).NotNil(db)
	r = db.Find("330322000000")
	a.NotNil(r).
		Equal(r.ID(), "22").
		Equal(r.FullID(), "330322000000").
		Equal(r.Name(), "洞头县").
		Equal(r.FullName(), "浙江省>温州市>洞头县").
		Equal(r.Versions(), []int{2009})
	r = db.Find("330305000000")
	a.Nil(r)

	// 所有年份的数据
	db, err = LoadFile("./data/regions.db", ">", true, version.Range(2009, 2020)...)
	a.NotError(err).NotNil(db)
	r = db.Find("330322000000")
	a.NotNil(r).
		Equal(r.ID(), "22").
		Equal(r.Versions(), []int{2014, 2013, 2012, 2011, 2010, 2009})
	r = db.Find("330305000000")
	a.NotNil(r).
		Equal(r.ID(), "05").
		Contains(r.Versions(), []int{2018, 2017, 2016, 2015})
}

func TestDB_versionIndex(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(0, obj.versionIndex(2020))
	a.Equal(1, obj.versionIndex(2019))
	a.Equal(-1, obj.versionIndex(1990))
}
