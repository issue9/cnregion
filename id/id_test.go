// SPDX-License-Identifier: MIT

package id

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSplit(t *testing.T) {
	a := assert.New(t)

	province, city, county, town, village := Split("330203103233")
	a.Equal(province, "33").
		Equal(city, "02").
		Equal(county, "03").
		Equal(town, "103").
		Equal(village, "233")

	a.Panic(func() {
		Split("3303")
	})
}

func TestSplitFilter(t *testing.T) {
	a := assert.New(t)

	list := SplitFilter("330203103000")
	a.Equal(4, len(list)).
		Equal(list[0], "33").
		Equal(list[1], "02").
		Equal(list[2], "03").
		Equal(list[3], "103")

	// 碰到第一个零值，即结果后续的判断
	list = SplitFilter("330003103000")
	a.Equal(1, len(list)).Equal(list[0], "33")
}

func TestParent(t *testing.T) {
	a := assert.New(t)

	a.Equal(Parent("330300000000"), "33")
	a.Equal(Parent("330302111000"), "330302")
}

func TestPrefix(t *testing.T) {
	a := assert.New(t)

	a.Equal(Prefix("330301001001"), "330301001001")
	a.Equal(Prefix("330300000000"), "3303")
	a.Equal(Prefix("330302000000"), "330302")
}

func TestFill(t *testing.T) {
	a := assert.New(t)

	a.Equal(Fill("34", Village), "340000000000")
	a.Equal(Fill("3", Village), "300000000000")
	a.Equal(Fill("34", Province), "34")
	a.Equal(Fill("34", City), "3400")
	a.Equal(Fill("341234666777", Village), "341234666777")
	a.Panic(func() {
		Fill("34112233444332", Village)
	})
}

func TestIsZero(t *testing.T) {
	a := assert.New(t)

	a.True(isZero("000"))
	a.False(isZero("00x"))
	a.True(isZero(""))
}
