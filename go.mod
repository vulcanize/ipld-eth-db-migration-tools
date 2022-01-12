module github.com/vulcanize/migration-tools

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.15
	github.com/ipfs/go-ipfs-blockstore v1.0.1
	github.com/ipfs/go-ipfs-ds-help v1.0.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.10.2
	github.com/multiformats/go-multihash v0.0.14
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.10.1
)

replace github.com/ethereum/go-ethereum v1.10.15 => github.com/vulcanize/go-ethereum v1.10.15-statediff-2.0.0
