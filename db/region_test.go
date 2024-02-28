// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cnregion/id"
)

func TestRegion_IsSupported(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{Items: []*Region{
		{Versions: []int{2020, 2019}, Name: "test", db: obj},
	}, db: obj}

	a.True(obj.root.Items[0].IsSupported(2020))
	a.True(obj.root.Items[0].IsSupported(2019))
	a.False(obj.root.Items[0].IsSupported(2018)) // 不支持
	a.False(obj.root.Items[0].IsSupported(2009)) // 不存在于 db
}

func TestRegion_addItem(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{Items: []*Region{}, db: obj}

	a.ErrorString(obj.root.addItem("33", "浙江", id.Province, 2001), "不支持该年份")

	a.NotError(obj.root.addItem("44", "广东", id.Province, 2020))
	a.Equal(obj.root.Items[0].ID, "44").
		NotNil(obj.root.Items[0].db).
		True(obj.root.Items[0].IsSupported(2020)).
		False(obj.root.Items[0].IsSupported(2019))

	a.ErrorString(obj.root.addItem("44", "广东", id.Province, 2020), "存在相同")
}

func TestRegion_SetSupported(t *testing.T) {
	a := assert.New(t, false)

	obj := &DB{versions: []int{2020, 2019, 2018}}
	obj.root = &Region{Items: []*Region{{db: obj}}, db: obj}

	a.NotError(obj.root.addItem("33", "浙江", id.Province, 2020))
	a.NotError(obj.root.Items[0].setSupported(2020)).
		Equal(obj.root.Items[0].Versions, []int{2020})
	a.NotError(obj.root.Items[0].setSupported(2019)).
		Equal(obj.root.Items[0].Versions, []int{2020, 2019})
	a.ErrorString(obj.root.Items[0].setSupported(2001), "不存在该年份")
}

func TestFindEnd(t *testing.T) {
	a := assert.New(t, false)

	data := []byte("0123{56}")
	a.Equal(findEnd(data), 7)
}
