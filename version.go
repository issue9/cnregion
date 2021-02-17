// SPDX-License-Identifier: MIT

package cnregion

import (
	"fmt"

	"github.com/issue9/cnregion/db"
)

// Version 用于描述与特定版本相关的区域数据
type Version struct {
	version int
	db      *db.DB
}

// New 返回 Version 实例
//
// version 表示需要的数据版本，即四位数的年份。
func New(db *db.DB, version int) *Version {
	if -1 == db.VersionIndex(version) {
		panic(fmt.Sprintf("版本号 %d 并不存在于 db", version))
	}

	return &Version{
		version: version,
		db:      db,
	}
}

// Load 从 path 加载数据并初始化 Version 实例
func Load(path, separator string, version int) (*Version, error) {
	d, err := db.Load(path, separator)
	if err != nil {
		return nil, err
	}

	return New(d, version), nil
}
