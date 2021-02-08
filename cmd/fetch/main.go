// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/cmdopt"
)

var (
	fetchDataDir string
	fetchYears   string
)

func main() {
	opt := &cmdopt.CmdOpt{
		Output:        os.Stdout,
		ErrorHandling: flag.ContinueOnError,
		CommandsTitle: "子命令",
		OptionsTitle:  "选项",
	}

	opt.Help("help", "显示当前命令\n")

	fetchFS := opt.New("fetch", "拉取数据\n", doFetch)
	fetchFS.StringVar(&fetchDataDir, "data", "./data", "指定数据的保存目录")
	fetchFS.StringVar(&fetchYears, "years", "", "指定年份，空值表示所有年份。")

	opt.Exec(os.Args[1:])
}

func doFetch(w io.Writer) error {
	if fetchYears == "" {
		return fetch(fetchDataDir)
	}

	yearList := strings.Split(fetchYears, ",")
	years := make([]int, 0, len(yearList))
	for _, y := range yearList {
		year, err := strconv.Atoi(strings.TrimSpace(y))
		if err != nil {
			return err
		}
		years = append(years, year)
	}

	return fetch(fetchDataDir, years...)
}
