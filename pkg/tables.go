// VulcanizeDB
// Copyright © 2022 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package migration_tools

import (
	"fmt"

	"github.com/vulcanize/migration-tools/pkg/eth_access_lists"
	"github.com/vulcanize/migration-tools/pkg/eth_accounts"
	"github.com/vulcanize/migration-tools/pkg/eth_headers"
	"github.com/vulcanize/migration-tools/pkg/eth_logs"
	"github.com/vulcanize/migration-tools/pkg/eth_receipts"
	"github.com/vulcanize/migration-tools/pkg/eth_state"
	"github.com/vulcanize/migration-tools/pkg/eth_storage"
	"github.com/vulcanize/migration-tools/pkg/eth_transactions"
	"github.com/vulcanize/migration-tools/pkg/eth_uncles"
	"github.com/vulcanize/migration-tools/pkg/interfaces"
	"github.com/vulcanize/migration-tools/pkg/public_nodes"
)

// TableName explicitly types table name strings
type TableName string

const (
	PublicNodes           TableName = "public.nodes"
	EthHeaders            TableName = "eth.header_cids"
	EthUncles             TableName = "eth.uncle_cids"
	EthTransactions       TableName = "eth.transaction_cids"
	EthAccessListElements TableName = "eth.access_list_elements"
	EthReceipts           TableName = "eth.receipt_cids"
	EthLogs               TableName = "eth.log_cids"
	EthState              TableName = "eth.state_cids"
	EthAccounts           TableName = "eth.state_accounts"
	EthStorage            TableName = "eth.storage_cids"
)

// NewTableTransformerSet inits and returns a set of Transformers for the provided tables
func NewTableTransformerSet(tables []TableName) map[TableName]interfaces.Transformer {
	tableReaderSet := make(map[TableName]interfaces.Transformer, len(tables))
	for _, tableName := range tables {
		tableReaderSet[tableName] = tableTransformerConstructorMappings[tableName]()
	}
	return tableReaderSet
}

// NewTableV2Model returns an allocation for a DB model of the provided table
func NewTableV2Model(tableName TableName) (interface{}, error) {
	switch tableName {
	case PublicNodes:
		return new(public_nodes.NodeModel), nil
	case EthHeaders:
		return new(eth_headers.HeaderModelV2WithMeta), nil
	case EthUncles:
		return new(eth_uncles.UncleModelV2WithMeta), nil
	case EthTransactions:
		return new(eth_transactions.TransactionModelV2WithMeta), nil
	case EthAccessListElements:
		return new(eth_access_lists.AccessListElementModelV2WithMeta), nil
	case EthReceipts:
		return new(eth_receipts.ReceiptModelV2WithMeta), nil
	case EthLogs:
		return new(eth_logs.LogModelV2WithMeta), nil
	case EthAccounts:
		return new(eth_accounts.AccountModelV2WithMeta), nil
	case EthStorage:
		return new(eth_storage.StorageModelV2WithMeta), nil
	default:
		return nil, fmt.Errorf("unsupported table name: %s", tableName)
	}
}

var tableTransformerConstructorMappings = map[TableName]interfaces.TransformerConstructor{
	PublicNodes:           public_nodes.NewTransformer,
	EthHeaders:            eth_headers.NewTransformer,
	EthUncles:             eth_uncles.NewTransformer,
	EthTransactions:       eth_transactions.NewTransformer,
	EthAccessListElements: eth_access_lists.NewTransformer,
	EthReceipts:           eth_receipts.NewTransformer,
	EthLogs:               eth_logs.NewTransformer,
	EthState:              eth_state.NewTransformer,
	EthAccounts:           eth_accounts.NewTransformer,
	EthStorage:            eth_storage.NewTransformer,
}

var tableReaderStrMappings = map[TableName]ReadPgStr{
	PublicNodes:           PgReadNodesStr,
	EthHeaders:            PgReadEthHeadersStr,
	EthUncles:             PgReadEthUnclesStr,
	EthTransactions:       PgReadEthTransactionsStr,
	EthAccessListElements: PgReadAccessListElementsStr,
	EthReceipts:           PgReadEthReceiptsStr,
	EthLogs:               PgReadEthLogsStr,
	EthState:              PgReadEthStateStr,
	EthAccounts:           PgReadEthAccountsStr,
	EthStorage:            PgReadEthStorageStr,
}

var tableWriterStrMappings = map[TableName]WritePgStr{
	PublicNodes:           PgWriteNodesStr,
	EthHeaders:            PgWriteEthHeadersStr,
	EthUncles:             PgWriteEthUnclesStr,
	EthTransactions:       PgWriteEthTransactionsStr,
	EthAccessListElements: PgWriteAccessListElementsStr,
	EthReceipts:           PgWriteEthReceiptsStr,
	EthLogs:               PgWriteEthLogsStr,
	EthState:              PgWriteEthStateStr,
	EthAccounts:           PgWriteEthAccountsStr,
	EthStorage:            PgWriteEthStorageStr,
}
