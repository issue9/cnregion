// SPDX-License-Identifier: MIT

package data

import (
	"embed"

	"github.com/issue9/cnregion"
)

//go:embed regions.db
var data embed.FS

// Embed 将 regions.db 的内容嵌入到程序中
//
// 这样可以让程序不依赖外部文件，但同时也会增加编译后程序的大小。
func Embed(separator string, version ...int) (*cnregion.Version, error) {
	return cnregion.LoadFS(data, "regions.db", separator, version...)
}
