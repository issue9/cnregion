// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

var (
	_ Region = &dbRegion{}
	_ Region = &districtRegion{}
)
