// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/cmdopt"
)

var (
	fetchDataDir string
	fetchYears   string

	buildDataDir string
	buildOutput  string
	buildYears   string
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

	buildFS := opt.New("build", "生成数据\n", doBuild)
	buildFS.StringVar(&buildDataDir, "data", "", "指定数据目录")
	buildFS.StringVar(&buildOutput, "output", "", "指定输出文件路径")
	buildFS.StringVar(&buildYears, "years", "", "指定年份，空值表示所有年份。")

	if err := opt.Exec(os.Args[1:]); err != nil {
		fmt.Fprintln(opt.Output, err)
		os.Exit(2)
	}
}

func doFetch(w io.Writer) error {
	years, err := getYears(fetchYears)
	if err != nil {
		return err
	}

	return fetch(fetchDataDir, years...)
}

func doBuild(w io.Writer) error {
	years, err := getYears(buildYears)
	if err != nil {
		return err
	}

	return build(buildDataDir, buildOutput, years...)
}

func getYears(years string) ([]int, error) {
	if years == "" {
		return nil, nil
	}

	yearList := strings.Split(years, ",")
	ys := make([]int, 0, len(yearList))
	for _, y := range yearList {
		year, err := strconv.Atoi(strings.TrimSpace(y))
		if err != nil {
			return nil, err
		}
		ys = append(ys, year)
	}

	return ys, nil
}
