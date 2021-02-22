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

func TestFill(t *testing.T) {
	a := assert.New(t)

	a.Equal(Fill("34", Village), "340000000000")
	a.Equal(Fill("34", Province), "34")
	a.Equal(Fill("34", City), "3400")
	a.Equal(Fill("341234666777", Village), "341234666777")
	a.Panic(func() {
		Fill("34112233444332", Village)
	})
}

func TestIsZero(t *testing.T) {
	a := assert.New(t)

	a.True(IsZero("000"))
	a.False(IsZero("00x"))
	a.True(IsZero(""))
}
