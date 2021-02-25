// SPDX-License-Identifier: MIT

package cnregion

import (
	"testing"

	"github.com/issue9/assert"
)

func TestVersion_Find(t *testing.T) {
	a := assert.New(t)

	// 2020
	v, err := LoadFile("./data/regions.db", ">", 2020)
	a.NotError(err).NotNil(v)
	r := v.Find("330305000000")
	a.NotNil(r).
		Equal(r.ID(), "05").
		Equal(r.FullID(), "330305000000").
		Equal(r.Name(), "洞头区").
		Equal(r.FullName(), "浙江省>温州市>洞头区")
	r = v.Find("330322000000") // 洞头县，已改为洞头区
	a.Nil(r)

	// 2009
	v, err = LoadFile("./data/regions.db", ">", 2009)
	a.NotError(err).NotNil(v)
	r = v.Find("330322000000")
	a.NotNil(r).
		Equal(r.ID(), "22").
		Equal(r.FullID(), "330322000000").
		Equal(r.Name(), "洞头县").
		Equal(r.FullName(), "浙江省>温州市>洞头县")
	r = v.Find("330305000000")
	a.Nil(r)

	// 所有年份的数据
	v, err = LoadFile("./data/regions.db", ">")
	a.NotError(err).NotNil(v)
	r = v.Find("330322000000")
	a.NotNil(r).Equal(r.ID(), "22")
	r = v.Find("330305000000")
	a.NotNil(r).Equal(r.ID(), "05")
}

func TestRegion_Items(t *testing.T) {
	a := assert.New(t)

	// 2020
	var x05, x22 bool
	v, err := LoadFile("./data/regions.db", ">", 2020)
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
	v, err = LoadFile("./data/regions.db", ">", 2009)
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
	v, err = LoadFile("./data/regions.db", ">", 2009, 2020)
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

func TestVersion_Provinces(t *testing.T) {
	a := assert.New(t)

	v, err := LoadFile("./data/regions.db", ">", 2020)
	a.NotError(err).NotNil(v)
	a.Equal(0, len(v.provinces))
	provinces := v.Provinces()
	a.Equal(31, len(provinces))
	a.Equal(31, len(provinces)) // 第二次读了缓存内容

	for _, p := range provinces {
		if p.ID() == "33" {
			a.Equal(p.Name(), "浙江省")
		}
	}
}

func TestVersion_Districts(t *testing.T) {
	a := assert.New(t)

	v, err := LoadFile("./data/regions.db", ">", 2020)
	a.NotError(err).NotNil(v)
	a.Equal(0, len(v.districts))
	districts := v.Districts()
	a.Equal(6, len(districts))
	a.Equal(6, len(districts))

	for _, d := range districts {
		if d.ID() == "1" {
			a.Equal(d.Name(), "华北地区").Equal(5, len(d.Items()))
		}
	}
}
