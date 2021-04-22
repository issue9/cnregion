// SPDX-License-Identifier: MIT

// +build cnregion

package cnregion

import _ "embed" // dbData

//go:embed data/regions.db
var dbData []byte

// Embed 将 data/regions.db 的内容嵌入到程序中
//
// 这样可以让程序不依赖外部文件，但同时也会增加编译后程序的大小。
// data/regions.db 目前大小为 7M 左右。
func Embed(separator string, version ...int) (*Version, error) {
	return Load(dbData, separator, version...)
}
