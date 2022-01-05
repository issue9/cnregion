// SPDX-License-Identifier: MIT

//go:build cnregion
// +build cnregion

package cnregion

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/issue9/cnregion/data"
)

func TestEmbed(t *testing.T) {
	a := assert.New(t, false)

	v, err := Load(data.Data, ">", 2020)
	a.NotError(err).NotNil(v)
	r := v.Find("330305000000")
	a.NotNil(r).
		Equal(r.ID(), "05").
		Equal(r.FullID(), "330305000000").
		Equal(r.Name(), "洞头区").
		Equal(r.FullName(), "浙江省>温州市>洞头区")
}
