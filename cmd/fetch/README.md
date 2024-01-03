# fetch

时间比较漫长，一年份的数据估计在 0.5 天左右。如果某个省的数据出错，会自动忽略该省的所有数据，下次再运行即可。

调整 colly limit 或是多协程运行 fetchTown 可以一定程序上提交效率。

如果出错，可以在执行完一轮之后重新再执行一次，会自动拉取有错误的数据。

拉取数据：
`
fetch fetch -years=2003,2004
`

生成数据：
`
fetch build -output=../data -data=./data
`
