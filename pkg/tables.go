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
	"github.com/jmoiron/sqlx"
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
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
)

// TableName explicitly types table name strings
type TableName string

const (
	PublicBlocks          TableName = "public.blocks"
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

var tableSet = []TableName{
	PublicBlocks,
	EthHeaders,
	EthUncles,
	EthTransactions,
	EthAccessListElements,
	EthReceipts,
	EthLogs,
	EthState,
	EthAccounts,
	EthStorage,
}

var numTables = len(tableSet)

// NewTableReaderSet inits and returns a set of Readers for the provided tables
func NewTableReaderSet(tables []TableName, db *sqlx.DB) map[TableName]interfaces.Reader {
	tableReaderSet := make(map[TableName]interfaces.Reader, len(tables))
	for _, tableName := range tables {
		tableReaderSet[tableName] = tableReaderConstructorMappings[tableName](db)
	}
	return tableReaderSet
}

// NewTableWriterSet inits and returns a set of Writers for the provided tables
func NewTableWriterSet(tables []TableName, db *sqlx.DB) map[TableName]interfaces.Writer {
	tableReaderSet := make(map[TableName]interfaces.Writer, len(tables))
	for _, tableName := range tables {
		tableReaderSet[tableName] = tableWriterConstructorMappings[tableName](db)
	}
	return tableReaderSet
}

// NewTableTransformerSet inits and returns a set of Transformers for the provided tables
func NewTableTransformerSet(tables []TableName) map[TableName]interfaces.Transformer {
	tableReaderSet := make(map[TableName]interfaces.Transformer, len(tables))
	for _, tableName := range tables {
		tableReaderSet[tableName] = tableTransformerConstructorMappings[tableName]()
	}
	return tableReaderSet
}

var tableReaderConstructorMappings = map[TableName]interfaces.ReaderConstructor{
	PublicBlocks:          public_blocks.NewReader,
	EthHeaders:            eth_headers.NewReader,
	EthUncles:             eth_uncles.NewReader,
	EthTransactions:       eth_transactions.NewReader,
	EthAccessListElements: eth_access_lists.NewReader,
	EthReceipts:           eth_receipts.NewReader,
	EthLogs:               eth_logs.NewReader,
	EthState:              eth_state.NewReader,
	EthAccounts:           eth_accounts.NewReader,
	EthStorage:            eth_storage.NewReader,
}

var tableWriterConstructorMappings = map[TableName]interfaces.WriterConstructor{
	PublicBlocks:          public_blocks.NewWriter,
	EthHeaders:            eth_headers.NewWriter,
	EthUncles:             eth_uncles.NewWriter,
	EthTransactions:       eth_transactions.NewWriter,
	EthAccessListElements: eth_access_lists.NewWriter,
	EthReceipts:           eth_receipts.NewWriter,
	EthLogs:               eth_logs.NewWriter,
	EthState:              eth_state.NewWriter,
	EthAccounts:           eth_accounts.NewWriter,
	EthStorage:            eth_storage.NewWriter,
}

var tableTransformerConstructorMappings = map[TableName]interfaces.TransformerConstructor{
	PublicBlocks:          public_blocks.NewTransformer,
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
