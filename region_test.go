// SPDX-License-Identifier: MIT

package cnregion

var (
	_ Region = &dbRegion{}
	_ Region = &districtRegion{}
)
