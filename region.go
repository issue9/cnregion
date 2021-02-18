// SPDX-License-Identifier: MIT

package cnregion

import (
	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
)

// Region 表示某个区域的相关信息
type Region struct {
	*db.Region
}

// Find 查找指定 ID 所表示的 Region
func (v *Version) Find(regionID string) *Region {
	province, city, county, town, village := id.Split(regionID)

	// 过滤掉零值
	items := []string{province, city, county, town, village}
	for index, item := range items {
		if id.IsZero(item) {
			items = items[:index]
			break
		}
	}

	dr := v.db.Find(items...)
	if dr == nil || !dr.IsSupported(v.version) {
		return nil
	}

	return &Region{Region: dr}
}

// Name 区域名称
func (r *Region) Name() string {
	return r.Region.Name
}

// FullName 全名
func (r *Region) FullName() string {
	return r.Region.FullName
}

// Items 子项
func (r *Region) Items() []*Region {
	items := make([]*Region, 0, len(r.Region.Items))
	for _, item := range r.Region.Items {
		items = append(items, &Region{Region: item})
	}
	return items
}
