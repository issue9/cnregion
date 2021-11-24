// SPDX-License-Identifier: MIT

package version

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestAll(t *testing.T) {
	a := assert.New(t, false)

	all := All()
	// 保证从大到小
	a.Equal(all[0], latest).
		Equal(all[len(all)-1], start)
}

func TestBeginWith(t *testing.T) {
	a := assert.New(t, false)

	list := BeginWith(latest)
	a.Equal(1, len(list)).Equal(list[0], latest)

	a.Panic(func() {
		BeginWith(start - 1)
	})
}
