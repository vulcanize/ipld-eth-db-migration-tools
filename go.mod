module github.com/vulcanize/migration-tools

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.14
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.10.1
)

replace github.com/ethereum/go-ethereum v1.10.14 => github.com/vulcanize/go-ethereum v1.10.14-statediff-0.0.29
