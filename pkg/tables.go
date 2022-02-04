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
	"strings"

	"github.com/vulcanize/migration-tools/pkg/csv"
	"github.com/vulcanize/migration-tools/pkg/sql"

	"github.com/vulcanize/migration-tools/pkg/eth_logs/repair"

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
	PublicNodes           TableName = "nodes"
	EthHeaders            TableName = "header_cids"
	EthUncles             TableName = "uncle_cids"
	EthTransactions       TableName = "transaction_cids"
	EthAccessListElements TableName = "access_list_elements"
	EthReceipts           TableName = "receipt_cids"
	EthLogs               TableName = "log_cids"
	EthLogsRepair         TableName = "log_cids_repair"
	EthState              TableName = "state_cids"
	EthAccounts           TableName = "state_accounts"
	EthStorage            TableName = "storage_cids"
	Unknown               TableName = "unknown"
)

// NewTableNameFromString returns the TableName from the provided string
func NewTableNameFromString(tableNameStr string) (TableName, error) {
	switch strings.ToLower(tableNameStr) {
	case "public.nodes", "nodes":
		return PublicNodes, nil
	case "eth.header_cids", "header_cids", "headers":
		return EthHeaders, nil
	case "eth.uncle_cids", "uncle_cids", "uncles":
		return EthUncles, nil
	case "eth.transaction_cids", "transaction_cids", "transactions", "txs", "trxs":
		return EthTransactions, nil
	case "eth.access_list_elements", "access_list_elements", "access_list":
		return EthAccessListElements, nil
	case "eth.receipt_cids", "receipt_cids", "receipts", "rcts":
		return EthReceipts, nil
	case "eth.log_cids", "log_cids", "logs":
		return EthLogs, nil
	case "eth.state_cids", "state_cids", "state":
		return EthState, nil
	case "eth.state_accounts", "state_accounts", "accounts":
		return EthAccounts, nil
	case "eth.storage_cids", "storage_cids", "storage":
		return EthStorage, nil
	case "eth.log_cids.repair", "log_repair", "log_cids_repair":
		return EthLogsRepair, nil
	default:
		return Unknown, fmt.Errorf("unrecognized table name: %s", tableNameStr)
	}
}

// NewTableTransformerSet inits and returns a set of Transformers for the provided tables
func NewTableTransformerSet(tables []TableName) map[TableName]interfaces.Transformer {
	tableReaderSet := make(map[TableName]interfaces.Transformer, len(tables))
	for _, tableName := range tables {
		tableReaderSet[tableName] = tableTransformerConstructorMappings[tableName]()
	}
	return tableReaderSet
}

// NewTableTransformer inits and returns a Transformers for the provided tables
func NewTableTransformer(table TableName) interfaces.Transformer {
	return tableTransformerConstructorMappings[table]()
}

// NewTableReadModels returns an allocation for the read DB models of the provided table
func NewTableReadModels(tableName TableName) (interface{}, error) {
	switch tableName {
	case PublicNodes:
		return new([]public_nodes.NodeModel), nil
	case EthHeaders:
		return new([]eth_headers.HeaderModelV2WithMeta), nil
	case EthUncles:
		return new([]eth_uncles.UncleModelV2WithMeta), nil
	case EthTransactions:
		return new([]eth_transactions.TransactionModelV2WithMeta), nil
	case EthAccessListElements:
		return new([]eth_access_lists.AccessListElementModelV2WithMeta), nil
	case EthReceipts:
		return new([]eth_receipts.ReceiptModelV2WithMeta), nil
	case EthLogs:
		return new([]eth_logs.LogModelV2WithMeta), nil
	case EthLogsRepair:
		return new([]eth_logs.LogModelV3), nil
	case EthAccounts:
		return new([]eth_accounts.AccountModelV2WithMeta), nil
	case EthStorage:
		return new([]eth_storage.StorageModelV2WithMeta), nil
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
	EthLogsRepair:         repair.NewTransformer,
	EthState:              eth_state.NewTransformer,
	EthAccounts:           eth_accounts.NewTransformer,
	EthStorage:            eth_storage.NewTransformer,
}

var tableReaderStrMappings = map[TableName]sql.ReadPgStr{
	PublicNodes:           sql.PgReadNodesStr,
	EthHeaders:            sql.PgReadEthHeadersStr,
	EthUncles:             sql.PgReadEthUnclesStr,
	EthTransactions:       sql.PgReadEthTransactionsStr,
	EthAccessListElements: sql.PgReadAccessListElementsStr,
	EthReceipts:           sql.PgReadEthReceiptsStr,
	EthLogs:               sql.PgReadEthLogsStr,
	EthLogsRepair:         sql.PgReadBrokenLogsStr,
	EthState:              sql.PgReadEthStateStr,
	EthAccounts:           sql.PgReadEthAccountsStr,
	EthStorage:            sql.PgReadEthStorageStr,
}

var tableWriterStrMappings = map[TableName]sql.WritePgStr{
	PublicNodes:           sql.PgWriteNodesStr,
	EthHeaders:            sql.PgWriteEthHeadersStr,
	EthUncles:             sql.PgWriteEthUnclesStr,
	EthTransactions:       sql.PgWriteEthTransactionsStr,
	EthAccessListElements: sql.PgWriteAccessListElementsStr,
	EthReceipts:           sql.PgWriteEthReceiptsStr,
	EthLogs:               sql.PgWriteEthLogsStr,
	EthLogsRepair:         sql.PgWriteIPLDsStr,
	EthState:              sql.PgWriteEthStateStr,
	EthAccounts:           sql.PgWriteEthAccountsStr,
	EthStorage:            sql.PgWriteEthStorageStr,
}

var csvWriterStrMappings = map[TableName]csv.WriteCSVStr{
	PublicNodes:           csv.CSVWriteNodesStr,
	EthHeaders:            csv.CSVWriteEthHeadersStr,
	EthUncles:             csv.CSVWriteEthUnclesStr,
	EthTransactions:       csv.CSVWriteEthTransactionsStr,
	EthAccessListElements: csv.CSVWriteAccessListElementsStr,
	EthReceipts:           csv.CSVWriteEthReceiptsStr,
	EthLogs:               csv.CSVWriteEthLogsStr,
	EthLogsRepair:         csv.CSVWriteIPLDsStr,
	EthState:              csv.CSVWriteEthStateStr,
	EthAccounts:           csv.CSVWriteEthAccountsStr,
	EthStorage:            csv.CSVWriteEthStorageStr,
}
