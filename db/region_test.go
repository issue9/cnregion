// SPDX-License-Identifier: MIT

package db

import (
	"testing"

	"github.com/issue9/assert"
)

func TestFindEnd(t *testing.T) {
	a := assert.New(t)

	data := []byte("0123{56}")
	a.Equal(findEnd(data), 7)
}
