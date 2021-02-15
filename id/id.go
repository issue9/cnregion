// SPDX-License-Identifier: MIT

// Package id 针对 ID 的一些操作函数
package id

import (
	"fmt"
	"strings"
)

// Length ID 的长度
const Length = 12

// Split 将一个区域 ID 按区域进行划分
func Split(id string) (province, city, county, town, village string) {
	if len(id) != Length {
		panic(fmt.Sprintf("id 的长度只能为 %d，当前为 %s", Length, id))
	}

	return id[:2], id[2:4], id[4:6], id[6:9], id[9:Length]
}

// Fill 为 id 填充后缀的 0
func Fill(id string) string {
	rem := Length - len(id)
	switch {
	case rem == 0:
		return id
	case rem > Length || rem < 2:
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
