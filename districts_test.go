// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestDB_Districts(t *testing.T) {
	a := assert.New(t, false)

	db, err := LoadFile("./data/regions.db", ">", true, 2020)
	a.NotError(err).NotNil(db)
	a.Length(db.Districts(), len(districtsMap))

	for _, d := range db.Districts() {
		if d.ID() == "1" {
			a.Equal(d.Name(), "华北地区").Equal(5, len(d.Items()))
		}
	}
}
