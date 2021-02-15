// SPDX-License-Identifier: MIT

// Package id 针对 ID 的一些操作函数
package id

import (
	"fmt"
	"strings"
)

// Split 将一个区域 ID 按区域进行划分
func Split(id string) (province, city, county, town, village string) {
	if len(id) != 12 {
		panic(fmt.Sprintf("id 的长度只能为 12，当前为 %s", id))
	}

	return id[:2], id[2:4], id[4:6], id[6:9], id[9:12]
}

// Fill 为 id 填充后缀的 0
func Fill(id string) string {
	rem := 12 - len(id)
	switch {
	case rem == 0:
		return id
	case rem > 12 || rem < 2:
		panic(fmt.Sprintf("无效的 id %s，无法为其填充 0", id))
	default:
		return id + strings.Repeat("0", rem)
	}
}

// IsZero 判断一组字符串是否都由 0 组成
func IsZero(id string) bool {
	for _, r := range id {
		if r != '0' {
			return false
		}
	}
	return true
}
