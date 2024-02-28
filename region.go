// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"github.com/issue9/cnregion/db"
	"github.com/issue9/cnregion/id"
)

// Region 表示某个区域的相关信息
type Region interface {
	ID() string       // 区域的 ID，不包括后缀 0 和上一级的 ID
	FullID() string   // 区域的 ID，包括后缀的 0 以及上一级的 ID，长度为 12
	Name() string     // 区域的名称
	FullName() string // 区域的全称，包括上一级的名称
	Items() []Region  // 子项
	Versions() []int  // 支持的年份版本
}

type dbRegion struct {
	r *db.Region
	v *Version
}

type districtRegion struct {
	v                          *Version
	id, name, fullName, fullID string
	items                      []Region
}

// Find 查找指定 ID 所表示的 Region
func (v *Version) Find(regionID string) Region {
	if dr := v.db.Find(id.SplitFilter(regionID)...); dr != nil {
		return &dbRegion{r: dr, v: v}
	}
	return nil
}

func (r *dbRegion) ID() string       { return r.r.ID }
func (r *dbRegion) Name() string     { return r.r.Name }
func (r *dbRegion) FullName() string { return r.r.FullName }
func (r *dbRegion) FullID() string   { return r.r.FullID }
func (r *dbRegion) Versions() []int  { return r.r.Versions }

func (r *dbRegion) Items() []Region {
	items := make([]Region, 0, len(r.r.Items))
	for _, item := range r.r.Items {
		items = append(items, &dbRegion{r: item, v: r.v})
	}
	return items
}

func (r *districtRegion) ID() string       { return r.id }
func (r *districtRegion) Name() string     { return r.name }
func (r *districtRegion) FullName() string { return r.fullName }
func (r *districtRegion) FullID() string   { return r.fullID }
func (r *districtRegion) Items() []Region  { return r.items }
func (r *districtRegion) Versions() []int  { return r.v.db.Versions() }
