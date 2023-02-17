# cnregion

[![Test](https://github.com/issue9/cnregion/workflows/Test/badge.svg)](https://github.com/issue9/cnregion/actions?query=workflow%3ATest)
[![Go version](https://img.shields.io/github/go-mod/go-version/issue9/cnregion)](https://golang.org)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/cnregion)](https://pkg.go.dev/github.com/issue9/cnregion)
[![codecov](https://codecov.io/gh/issue9/cnregion/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/cnregion)
![License](https://img.shields.io/github/license/issue9/cnregion)

历年统计用区域和城乡划分代码，数据来源于 <http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/>。
符合国家标准 GB/T 2260 与 GB/T 10114。

关于版本号，主版本号代码不兼容性更改，次版本号代码最后一次生成的数据年份，BUG 修正和兼容性的功能增加则增加修订版本号。

```go
v, err := cnregion.LoadFile("./data/regions.db", "-", 2020)

p := v.Provinces() // 返回所有省列表
cities := p[0].Items() // 返回该省下的所有市
counties := cities[0].Items() // 返回该市下的所有县
towns := counties[0].Items() // 返回所有镇
villages := towns[0].Items() // 所有村和街道信息

d := v.Districts() // 按以前的行政大区进行划分
provinces := d[0].Items() // 该大区下的所有省份

list := v.Search(&SearchOptions{Text: ""温州"}) // 按索地名中带温州的区域列表
```

对采集的数据进行了一定的加工，以减少文件的体积，文件保存在 `./data/regions.db` 中。

## 安装

```shell
go get github.com/issue9/cnregion
```

## 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
