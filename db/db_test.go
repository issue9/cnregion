// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"
)

var data = []byte(`[2020,2019]::中国:1:2{33:浙江:1:0{}34:安徽:1:2{01:合肥:1:0{}02:芜湖:1:0{}}}`)

var obj = &DB{
	Versions: []int{2020, 2019},
	Region: &Region{
		Name:      "中国",
		Supported: 1,
		Items: []*Region{
			{
				ID:        "33",
				Name:      "浙江",
				Supported: 1,
			},
			{
				ID:        "34",
				Name:      "安徽",
				Supported: 1,
				Items: []*Region{
					{
						ID:        "01",
						Name:      "合肥",
						Supported: 1,
					},
					{
						ID:        "02",
						Name:      "芜湖",
						Supported: 1,
					},
				},
			},
		},
	},
}

func TestMarshal(t *testing.T) {
	a := assert.New(t)

	o1, err := Unmarshal(data)
	a.NotError(err).Equal(o1, obj)

	d1, err := Marshal(obj)
	a.NotError(err).NotNil(d1)
	a.Equal(string(d1), string(data))
}
