module github.com/zaidfadhil/cerra/redis

go 1.24

require (
	github.com/redis/go-redis/v9 v9.17.2
	github.com/zaidfadhil/cerra v0.1.5
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace github.com/zaidfadhil/cerra => ../
