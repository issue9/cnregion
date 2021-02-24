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
	FullID() string
	Items() []Region
}

type dbRegion struct {
	r *db.Region
	v *Version
}

type districtRegion struct {
	id, name, fullName, fullID string
	items                      []Region
}

// Find 查找指定 ID 所表示的 Region
func (v *Version) Find(regionID string) Region {
	dr := v.db.Find(id.SplitFilter(regionID)...)
	if dr == nil || !dr.IsSupported(v.version) {
		return nil
	}

	return &dbRegion{r: dr, v: v}
}

func (r *dbRegion) ID() string       { return r.r.ID }
func (r *dbRegion) Name() string     { return r.r.Name }
func (r *dbRegion) FullName() string { return r.r.FullName }
func (r *dbRegion) FullID() string   { return r.r.FullID }

func (r *dbRegion) Items() []Region {
	items := make([]Region, 0, len(r.r.Items))
	for _, item := range r.r.Items {
		if item.IsSupported(r.v.version) {
			items = append(items, &dbRegion{r: item, v: r.v})
		}
	}
	return items
}

func (r *districtRegion) ID() string       { return r.id }
func (r *districtRegion) Name() string     { return r.name }
func (r *districtRegion) FullName() string { return r.fullName }
func (r *districtRegion) FullID() string   { return r.fullID }
func (r *districtRegion) Items() []Region  { return r.items }
