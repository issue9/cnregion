// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"
)

func TestRegion_IsSupported(t *testing.T) {
	a := assert.New(t)

	a.True(obj.Items[0].IsSupported(obj, 2020))
	a.False(obj.Items[0].IsSupported(obj, 2019))

	a.True(obj.Items[1].Items[0].IsSupported(obj, 2020))
	a.True(obj.Items[1].Items[0].IsSupported(obj, 2019))
	a.False(obj.Items[1].Items[0].IsSupported(obj, 2018)) // 2018 不存在于 obj.Versions
}

func TestRegion_AddItem(t *testing.T) {
	a := assert.New(t)

	obj := &DB{Versions: []int{2020, 2019, 2018}, Region: &Region{Items: []*Region{}}}
	a.ErrorString(obj.AddItem(obj, "33", "浙江", 2001), "不支持该年份")
	a.NotError(obj.AddItem(obj, "33", "浙江", 2020))
	a.ErrorString(obj.AddItem(obj, "33", "浙江", 2020), "存在相同")
}

func TestRegion_SetSupported(t *testing.T) {
	a := assert.New(t)

	obj := &DB{Versions: []int{2020, 2019, 2018}, Region: &Region{Items: []*Region{}}}
	a.NotError(obj.AddItem(obj, "33", "浙江", 2020))
	a.NotError(obj.Items[0].SetSupported(obj, 2020))
	a.NotError(obj.Items[0].SetSupported(obj, 2019))
	a.ErrorString(obj.Items[0].SetSupported(obj, 2001), "不存在该年份")
}

func TestFindEnd(t *testing.T) {
	a := assert.New(t)

	data := []byte("0123{56}")
	a.Equal(findEnd(data), 7)
}
