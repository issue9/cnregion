// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cnregion/v2/id"
)

func TestRegion_IsSupported(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{items: []*Region{
		{versions: []int{2020, 2019}, name: "test", db: obj},
	}, db: obj}

	a.True(obj.root.items[0].IsSupported(2020))
	a.True(obj.root.items[0].IsSupported(2019))
	a.False(obj.root.items[0].IsSupported(2018)) // 不支持
	a.False(obj.root.items[0].IsSupported(2009)) // 不存在于 db
}

func TestRegion_addItem(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{items: []*Region{}, db: obj}

	a.ErrorString(obj.root.addItem("33", "浙江", id.Province, 2001), "不支持该年份")

	a.NotError(obj.root.addItem("44", "广东", id.Province, 2020))
	a.Equal(obj.root.items[0].id, "44").
		NotNil(obj.root.items[0].db).
		True(obj.root.items[0].IsSupported(2020)).
		False(obj.root.items[0].IsSupported(2019))

	a.ErrorString(obj.root.addItem("44", "广东", id.Province, 2020), "存在相同")
}

func TestRegion_SetSupported(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{items: []*Region{{db: obj}}, db: obj}

	a.NotError(obj.root.addItem("33", "浙江", id.Province, 2020))
	a.NotError(obj.root.items[0].setSupported(2020)).
		Equal(obj.root.items[0].versions, []int{2020})
	a.NotError(obj.root.items[0].setSupported(2019)).
		Equal(obj.root.items[0].versions, []int{2020, 2019})
	a.ErrorString(obj.root.items[0].setSupported(2001), "不存在该年份")
}

func TestFindEnd(t *testing.T) {
	a := assert.New(t, false)

	data := []byte("0123{56}")
	a.Equal(findEnd(data), 7)
}

func TestDB_Provinces(t *testing.T) {
	a := assert.New(t, false)

	v, err := LoadFile("./data/regions.db", ">", true, 2020)
	a.NotError(err).NotNil(v)

	for _, p := range v.Provinces() {
		if p.ID() == "33" {
			a.Equal(p.Name(), "浙江省")
		}
	}
}

func TestRegion_Items(t *testing.T) {
	a := assert.New(t, false)

	// 2020
	var x05, x22 bool
	v, err := LoadFile("./data/regions.db", ">", true, 2020)
	a.NotError(err).NotNil(v)
	r := v.Find("330300000000")
	for _, item := range r.Items() {
		if item.ID() == "05" {
			x05 = true
		}
		if item.ID() == "22" {
			x22 = true
		}
	}
	a.True(x05).False(x22)

	// 2009
	x05 = false
	x22 = false
	v, err = LoadFile("./data/regions.db", ">", true, 2009)
	a.NotError(err).NotNil(v)
	r = v.Find("330300000000")
	for _, item := range r.Items() {
		if item.ID() == "05" {
			x05 = true
		}
		if item.ID() == "22" {
			x22 = true
		}
	}
	a.False(x05).True(x22)

	//2020 + 2009
	x05 = false
	x22 = false
	v, err = LoadFile("./data/regions.db", ">", true, 2009, 2020)
	a.NotError(err).NotNil(v)
	r = v.Find("330300000000")
	for _, item := range r.Items() {
		if item.ID() == "05" {
			x05 = true
		}
		if item.ID() == "22" {
			x22 = true
		}
	}
	a.True(x05).True(x22)
}

func TestRegion_findItem(t *testing.T) {
	a := assert.New(t, false)

	r := obj.root.findItem("34", "01")
	a.NotNil(r).Equal(r.name, "合肥").Equal(r.fullName, "安徽-合肥")

	r = obj.root.findItem("34", "01", "00")
	a.Nil(r)

	r = obj.root.findItem("34")
	a.NotNil(r).Equal(r.name, "安徽").Equal(r.fullName, "安徽")

	r = obj.root.findItem()
	a.NotNil(r).Equal(r.name, "").Equal(r.fullName, "").Equal(2, len(r.items))

	// 不存在于 obj
	a.Nil(obj.root.findItem("99"))
	a.Nil(obj.root.findItem(""))
}
