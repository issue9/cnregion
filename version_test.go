// SPDX-License-Identifier: MIT

package cnregion

import (
	"testing"

	"github.com/issue9/assert"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	v, err := LoadFile("./data/regions.db", ">", 2020)
	a.NotError(err).NotNil(v)
	r := v.Find("330305000000")
	a.NotNil(r).
		Equal(r.ID(), "05").
		Equal(r.Name(), "洞头区").
		Equal(r.FullName(), "浙江省>温州市>洞头区")
	r = v.Find("330322000000") // 洞头县，已改为洞头区
	a.Nil(r)

	v, err = LoadFile("./data/regions.db", ">", 2009)
	a.NotError(err).NotNil(v)
	r = v.Find("330322000000")
	a.NotNil(r).
		Equal(r.ID(), "22").
		Equal(r.Name(), "洞头县").
		Equal(r.FullName(), "浙江省>温州市>洞头县")
	r = v.Find("330305000000")
	a.Nil(r)
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
