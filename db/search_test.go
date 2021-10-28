// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/cnregion/id"
)

func TestDB_Search(t *testing.T) {
	a := assert.New(t)

	rs := obj.Search(&Options{Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].Name, "合肥")

	rs = obj.Search(&Options{Parent: "340000000000", Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].Name, "合肥")

	rs = obj.Search(&Options{Parent: "000000000000", Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].Name, "合肥")

	// 限定 level 只能是省以及 parent 为 34 开头
	rs = obj.Search(&Options{Parent: "340000000000", Level: id.Province, Text: "合肥"})
	a.Equal(0, len(rs))

	// 未限定 parent 且 level 正确
	rs = obj.Search(&Options{Level: id.City, Text: "合肥"})
	a.Equal(1, len(rs))

	rs = obj.Search(&Options{Level: id.City, Text: "湖"})
	a.Equal(2, len(rs))

	rs = obj.Search(&Options{Level: id.City, Parent: "340000000000", Text: "湖"})
	a.Equal(2, len(rs))

	// parent = 浙江
	rs = obj.Search(&Options{Parent: "330000000000", Text: "合肥"})
	a.Equal(0, len(rs))

	// parent 不存在
	rs = obj.Search(&Options{Parent: "110000000000", Text: "合肥"})
	a.Equal(0, len(rs))
}
