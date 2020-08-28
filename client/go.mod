module github.com/RobinUS2/tsxdb/client

go 1.13

replace github.com/RobinUS2/tsxdb/client => ../client

replace github.com/RobinUS2/tsxdb/rpc => ../rpc

replace github.com/RobinUS2/tsxdb/rpc/types => ../rpc/types

replace github.com/RobinUS2/tsxdb/tools => ../tools

require (
	github.com/RobinUS2/tsxdb/rpc v0.0.0-20200828082527-116be78017fa
	github.com/RobinUS2/tsxdb/tools v0.0.0-20200828082527-116be78017fa
	github.com/karlseguin/ccache/v2 v2.0.6
	github.com/pkg/errors v0.9.1
)
