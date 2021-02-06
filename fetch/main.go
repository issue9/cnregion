// SPDX-License-Identifier: MIT

// +build ignore

package main

import "github.com/issue9/cnregion/fetch"

func main() {
	if err := fetch.Fetch("./data"); err != nil {
		panic(err)
	}
}
