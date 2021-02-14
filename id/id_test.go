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
