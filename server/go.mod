module github.com/RobinUS2/tsxdb/server

go 1.12

replace github.com/RobinUS2/tsxdb/client => ../client

replace github.com/RobinUS2/tsxdb/rpc => ../rpc

replace github.com/RobinUS2/tsxdb/server => ../server

replace github.com/RobinUS2/tsxdb/tools => ../tools

require (
	github.com/RobinUS2/tsxdb/client v0.0.0-00010101000000-000000000000
	github.com/RobinUS2/tsxdb/rpc v0.0.0-00010101000000-000000000000
	github.com/RobinUS2/tsxdb/tools v0.0.0-00010101000000-000000000000
	github.com/bsm/redis-lock v8.0.0+incompatible
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/pkg/errors v0.8.1
	k8s.io/apimachinery v0.0.0-20190511063452-5b67e417bf61
)
