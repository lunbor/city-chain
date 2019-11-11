module github.com/scryinfo/citychain/chain_server

go 1.13

require (
	github.com/scryinfo/citychain/dots v0.0.0
	github.com/scryinfo/dot v0.1.3-0.20191026032307-4fe8cc8e04c9
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	go.uber.org/zap v1.12.0
)

replace (
	github.com/scryinfo/citychain/api/server v0.0.0 => ../../api/server
	github.com/scryinfo/citychain/dots v0.0.0 => ../../dots
)
