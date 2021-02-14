// SPDX-License-Identifier: MIT

package version

import (
	"testing"

	"github.com/issue9/assert"
)

func TestAll(t *testing.T) {
	a := assert.New(t)

	all := All()
	// 保证从大到小
	a.Equal(all[0], Latest).
		Equal(all[len(all)-1], Start)
}
