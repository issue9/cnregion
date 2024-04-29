// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package cnregion 中国区域划分代码
//
// 中国行政区域五级划分代码，包含了省、市、县、乡和村五个级别。
// [数据规则]以及[数据来源]。
//
// [数据规则]: http://www.stats.gov.cn/tjsj/tjbz/200911/t20091125_8667.html
// [数据来源]: http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/
package cnregion

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
)

// LoadFS 从数据文件加载数据
func LoadFS(f fs.FS, file, separator string, compress bool, version ...int) (*DB, error) {
	data, err := fs.ReadFile(f, file)
	if err != nil {
		return nil, err
	}
	return Load(data, separator, compress, version...)
}

// Load 将数据内容加载至 DB 对象
//
// version 仅加载指定年份的数据，如果为空，则加载所有数据；
func Load(data []byte, separator string, compress bool, version ...int) (*DB, error) {
	if compress {
		rd, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		data, err = io.ReadAll(rd)
		if err != nil {
			return nil, err
		}
	}

	db := &DB{
		fullNameSeparator: separator,
		filters:           version,
	}
	if err := db.unmarshal(data); err != nil {
		return nil, err
	}

	db.initDistricts()

	return db, nil
}

// LoadFile 从数据文件加载数据
func LoadFile(file, separator string, compress bool, version ...int) (*DB, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Load(data, separator, compress, version...)
}

// Dump 输出到文件
func (db *DB) Dump(file string, compress bool) error {
	data, err := db.marshal()
	if err != nil {
		return err
	}

	if compress {
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err = w.Write(data); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}

		data = buf.Bytes()
	}

	return os.WriteFile(file, data, os.ModePerm)
}
