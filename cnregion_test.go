// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cnregion/v2/id"
	"github.com/issue9/cnregion/v2/version"
)

func TestLoad(t *testing.T) {
	a := assert.New(t, false)

	o1, err := Load(data, "-", false)
	a.NotError(err).
		Equal(o1.fullNameSeparator, obj.fullNameSeparator).
		Equal(o1.versions, obj.versions).
		Equal(len(o1.root.items), len(obj.root.items)).
		Equal(o1.root.items[0].id, obj.root.items[0].id).
		Equal(o1.root.items[0].fullID, obj.root.items[0].fullID).
		Equal(o1.root.items[0].items[0].id, obj.root.items[0].items[0].id).
		Equal(o1.root.items[1].items[0].id, obj.root.items[1].items[0].id).
		Equal(o1.root.items[1].items[0].fullID, obj.root.items[1].items[0].fullID).
		Equal(o1.root.items[1].items[1].fullID, obj.root.items[1].items[1].fullID).
		NotEqual(o1.root.items[1].items[1].fullID, obj.root.items[1].items[0].fullID)

	d1, err := obj.marshal()
	a.NotError(err).NotNil(d1)
	a.Equal(string(d1), string(data))

	_, err = Load([]byte("100:[2020]:::1:0{}"), "-", false)
	a.Equal(err, ErrIncompatible)

	o1, err = Load(data, "-", false, 2019)
	a.NotError(err).
		Equal(0, len(o1.root.items))
}

func TestDB_LoadDump(t *testing.T) {
	a := assert.New(t, false)

	path := filepath.Join(os.TempDir(), "cnregion_db.dict")
	a.NotError(obj.Dump(path, false))
	d, err := LoadFile(path, "-", false)
	a.NotError(err).NotNil(d)

	path = filepath.Join(os.TempDir(), "cnregion_db_compress.dict")
	a.NotError(obj.Dump(path, true))
	d, err = LoadFile(path, "-", true)
	a.NotError(err).NotNil(d)
}

func TestLoadFS(t *testing.T) {
	a := assert.New(t, false)

	obj, err := LoadFS(os.DirFS("./data"), "regions.db", "-", true)
	a.NotError(err).NotNil(obj)
	a.Equal(obj.versions, version.All()).
		Equal(obj.fullNameSeparator, "-").
		True(len(obj.root.items) > 0).
		Equal(obj.root.items[0].level, id.Province).
		Equal(obj.root.items[0].items[0].level, id.City).
		Equal(obj.root.items[0].items[0].items[0].level, id.County).
		Equal(obj.root.items[1].level, id.Province).
		Equal(obj.root.items[2].items[0].level, id.City)
}
