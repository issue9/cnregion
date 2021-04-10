module github.com/issue9/cnregion/cmd/fetch

go 1.16

require (
	github.com/gocolly/colly/v2 v2.1.0
	github.com/issue9/cmdopt v0.7.0
	github.com/issue9/errwrap v0.2.0
	github.com/issue9/cnregion v0.1.0
)

replace (
	github.com/issue9/cnregion => ../../
)
