// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v3/colors"
)

func main() {
	const usage = `fetch
	commands:
	{{commands}}

	flag:
	{{flags}}
	`
	opt := cmdopt.New(os.Stdout, flag.ContinueOnError, usage, nil, func(s string) string { return fmt.Sprintf("not found %s", s) })
	cmdopt.Help(opt, "help", "显示当前命令\n", "显示当前命令\n")

	opt.New("fetch", "拉取数据\n", "拉取数据\n", doFetch)

	opt.New("build", "生成数据\n", "生成数据\n", doBuild)

	if err := opt.Exec(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(2)
	}
}

func doFetch(fs *flag.FlagSet) cmdopt.DoFunc {
	var (
		fetchDataDir  string
		fetchYears    string
		fetchInterval string
	)
	fs.StringVar(&fetchDataDir, "data", "./data", "指定数据的保存目录")
	fs.StringVar(&fetchYears, "years", "", "指定年份，空值表示所有年份。格式 y1,y2。")
	fs.StringVar(&fetchInterval, "internal", "1m", "每拉取一个省份数据后的间隔时间。")

	return func(w io.Writer) error {
		years, err := getYears(fetchYears)
		if err != nil {
			return err
		}

		interval, err := time.ParseDuration(fetchInterval)
		if err != nil {
			return err
		}

		return fetch(fetchDataDir, interval, years...)
	}
}

func doBuild(fs *flag.FlagSet) cmdopt.DoFunc {
	var (
		buildDataDir string
		buildOutput  string
		buildYears   string
	)
	fs.StringVar(&buildDataDir, "data", "", "指定数据目录")
	fs.StringVar(&buildOutput, "output", "", "指定输出文件路径")
	fs.StringVar(&buildYears, "years", "", "指定年份，空值表示所有年份。格式 y1,y2。")

	return func(io.Writer) error {
		years, err := getYears(buildYears)
		if err != nil {
			return err
		}

		return build(buildDataDir, buildOutput, years...)
	}
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

func colorsSprintf(fore colors.Color, format string, v ...any) string {
	return colors.Sprintf(colors.Normal, fore, colors.Default, format, v...)
}

func colorsSprint(fore colors.Color, v ...any) string {
	return colors.Sprint(colors.Normal, fore, colors.Default, v...)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !errors.Is(err, os.ErrNotExist)
}
