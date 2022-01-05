// SPDX-License-Identifier: MIT

//go:build cnregion
// +build cnregion

package cnregion

// Embed 将 data/regions.db 的内容嵌入到程序中
//
// 这样可以让程序不依赖外部文件，但同时也会增加编译后程序的大小。
//
// Deprecated: 请使用 data.Data 代替
func Embed(separator string, version ...int) (*Version, error) {
	return Load(data.Data, separator, version...)
}
