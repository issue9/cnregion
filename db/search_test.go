// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/cnregion/id"
)

func TestDB_Search(t *testing.T) {
	a := assert.New(t)

	rs := obj.Search("合肥", nil)
	a.Equal(1, len(rs)).Equal(rs[0].Name, "合肥")

	rs = obj.Search("合肥", &Options{})
	a.Equal(1, len(rs)).Equal(rs[0].Name, "合肥")

	rs = obj.Search("合肥", &Options{Parent: "340000000000"})
	a.Equal(1, len(rs)).Equal(rs[0].Name, "合肥")

	// 限定 level 只能是省
	rs = obj.Search("合肥", &Options{Parent: "340000000000", Level: id.Province})
	a.Equal(0, len(rs))

	// parent = 浙江
	rs = obj.Search("合肥", &Options{Parent: "330000000000"})
	a.Equal(0, len(rs))

	// parent 不存在
	rs = obj.Search("合肥", &Options{Parent: "110000000000"})
	a.Equal(0, len(rs))
}
