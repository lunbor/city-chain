module github.com/scryinfo/citychain/dots

go 1.13

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/ethereum/go-ethereum v1.10.17
	github.com/go-kit/kit v0.8.0
	github.com/jinzhu/gorm v1.9.11
	github.com/pkg/errors v0.9.1
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/scryinfo/citychain/api/server v0.0.0
	github.com/scryinfo/dot v0.1.3-0.20191026032307-4fe8cc8e04c9
	github.com/scryinfo/scryg v0.1.3-0.20190608053141-a292b801bfd6
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/ybbus/jsonrpc v2.1.2+incompatible
	go.uber.org/atomic v1.5.0
	go.uber.org/zap v1.12.0
)

replace github.com/scryinfo/citychain/api/server v0.0.0 => ../api/server
