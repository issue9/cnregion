// SPDX-License-Identifier: MIT

package cnregion

import (
	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
)

// Region 表示某个区域的相关信息
type Region interface {
	ID() string
	Name() string
	FullName() string
	Items() []Region
}

type dbRegion struct {
	r *db.Region
}

type districtRegion struct {
	id, name, fullName string
	items              []Region
}

// Find 查找指定 ID 所表示的 Region
func (v *Version) Find(regionID string) Region {
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

	return &dbRegion{r: dr}
}

func (r *dbRegion) ID() string       { return r.r.ID }
func (r *dbRegion) Name() string     { return r.r.Name }
func (r *dbRegion) FullName() string { return r.r.FullName }

func (r *dbRegion) Items() []Region {
	items := make([]Region, 0, len(r.r.Items))
	for _, item := range r.r.Items {
		items = append(items, &dbRegion{r: item})
	}
	return items
}

func (r *districtRegion) ID() string       { return r.id }
func (r *districtRegion) Name() string     { return r.name }
func (r *districtRegion) FullName() string { return r.fullName }
func (r *districtRegion) Items() []Region  { return r.items }
