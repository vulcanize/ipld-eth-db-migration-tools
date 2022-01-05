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

// ReadPgStr provides explicit typing for read postgres statements
type ReadPgStr string

// WritePgStr provides explicit typing for write postgres statements
type WritePgStr string

const (
	PgReadNodesStr ReadPgStr = `SELECT client_name, genesis_block, network_id, node_id, nodes_chain_id
						FROM public.nodes`

	PgReadEthUnclesStr ReadPgStr = `SELECT eth.header_cids.block_hash, eth.uncle_cids.*
							FROM eth.uncle_cids
							INNER JOIN eth.header_cids ON (uncle_cids.header_id = header_cids.id)
							WHERE block_number BETWEEN $1 AND $2`

	PgReadEthTransactionsStr ReadPgStr = `SELECT public.blocks.data, eth.header_cids.block_hash, eth.transaction_cids.*
								FROM eth.transaction_cids
								INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
								INNER JOIN public.blocks ON (transaction_cids.mh_key = blocks.key)
								WHERE block_number BETWEEN $1 AND $2`

	PgReadEthStorageStr ReadPgStr = `SELECT eth.header_cids.block_hash, eth.state_cids.state_path, eth.storage_cids.*
							FROM eth.storage_cids
							INNER JOIN eth.state_cids ON (storage_cids.state_id = state_cids.id)
							INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
							WHERE block_number BETWEEN $1 AND $2`

	PgReadEthStateStr ReadPgStr = `SELECT eth.header_cids.block_hash, eth.state_cids.*
						FROM eth.state_cids
						INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
						WHERE block_number = $1`

	PgReadEthReceiptsStr ReadPgStr = `SELECT eth.transaction_cids.tx_hash, eth.receipt_cids.*
							FROM eth.receipt_cids
							INNER JOIN eth.transaction_cids ON (receipt_cids.tx_id = transaction_cids.id)
							INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
							WHERE block_number BETWEEN $1 AND $2`

	PgReadEthLogsStr ReadPgStr = `SELECT eth.transaction_cis.tx_hash, eth.log_cids.*
						FROM eth.log_cids
						INNER JOIN eth.receipt_cids ON (log_cids.receipt_id = receipt_cids.id)
						INNER JOIN eth.transaction_cids ON (receipt_cids.tx_id = transaction_cids.id)
						INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
						WHERE block_number BETWEEN $1 AND $2`

	PgReadEthHeadersStr ReadPgStr = `SELECT public.blocks.data, eth.header_cids.*
							FROM eth.header_cids
							INNER JOIN public.blocks ON (header_cids.mh_key = blocks.key)
							WHERE block_number BETWEEN $1 AND $2`

	PgReadEthAccountsStr ReadPgStr = `SELECT eth.header_cids.block_hash, eth.state_cids.state_path, eth.state_accounts.*
							FROM eth.state_accounts
							INNER JOIN eth.state_cids ON (state_accounts.state_id = state_cids.id)
							INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
							WHERE block_number BETWEEN $1 AND $2`

	PgReadAccessListElementsStr ReadPgStr = `SELECT eth.transaction_cids.tx_hash, eth.access_list_elements.*
									FROM eth.access_list_elements
									INNER JOIN eth.transaction_cids ON (access_list_elements.tx_id = transaction_cids.id)
									INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
									WHERE block_number BETWEEN $1 AND $2`

	PgWriteEthUnclesStr WritePgStr = `INSERT INTO eth.uncle_cids (header_id, block_hash, parent_hash, cid, mh_key, reward)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR(66)[]),
							unnest($4::TEXT[]), unnest($5::TEXT[]), unnest($6::NUMERIC[]))`

	PgWriteEthTransactionsStr WritePgStr = `INSERT INTO eth.transaction_cids (header_id, index, tx_hash, cid, mh_key, dst, src,
									tx_data, tx_type, value)
									VALUES (unnest($1::VARCHAR(66)[]), unnest($2::INTEGER[]), unnest($3::VARCHAR(66)[]),
									unnest($4::TEXT[]), unnest($5::TEXT[]), unnest($6::VARCHAR(66)[]),
									unnest($7::VARCHAR(66)[]), unnest($8::BYTEA[]), unnest($9::INTEGER[]),
									unnest($10::NUMERIC[]))`

	PgWriteEthStorageStr WritePgStr = `INSERT INTO eth.storage_cids (header_id, state_path, storage_path, storage_leaf_key, node_type,
							cid, mh_key, diff)
							VALUES (:header_id, :state_path, :storage_path, :storage_leaf_key, :node_type,
							:cid, :mh_key, :diff)`

	PgWriteEthStateStr WritePgStr = `INSERT INTO eth.state_cids (header_id, state_path, state_leaf_key, node_type, cid, mh_key, diff)
						VALUES (:header_id, :state_path, :state_leaf_key, :node_type, :cid, :mh_key, :diff)`

	PgWriteNodesStr WritePgStr = `INSERT INTO public.nodes (client_name, genesis_block, network_id, node_id, chain_id)
						VALUES (:client_name, :genesis_block, :network_id, :node_id, :chain_id)`

	PgWriteEthReceiptsStr WritePgStr = `INSERT INTO eth.receipt_cids (tx_id, leaf_cid, leaf_mh_key, post_status, post_state,
								contract, contract_hash, log_root)
								VALUES (:tx_id, :leaf_cid, :leaf_mh_key, :post_status, :post_state,
								:contract, :contract_hash, :log_root)`

	PgWriteEthLogsStr WritePgStr = `INSERT INTO eth.logs_cids (rct_id, leaf_cid, leaf_mh_key, address, index, log_data, topic0,
							topic1, topic2, topic3)
							VALUES (:rct_id, :leaf_cid, :leaf_mh_key, :address, :index, :log_data, :topic0,
							:topic1, :topic2, :topic3)`

	PgWriteEthHeadersStr WritePgStr = `INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, mh_key, td, node_id,
							reward, state_root, uncle_root, tx_root, receipt_root, bloom, timestamp, times_validated, coinbase)
							VALUES (:block_number, :block_hash, :parent_hash, :cid, :mh_key, :td, :node_id, :reward,
							:state_root, :uncle_root, :tx_root, :receipt_root, :bloom, :timestamp, :times_validated, :coinbase)`

	PgWriteEthAccountsStr WritePgStr = `INSERT INTO eth.state_accounts (header_id, state_path, balance, nonce, code_hash, storage_root)
							VALUES (:header_id, :state_path, :balance, :nonce, :code_hash, :storage_root)`

	PgWriteAccessListElementsStr WritePgStr = `INSERT INTO eth.access_list_elements (tx_id, index, address, storage_keys)
									VALUES (:tx_id, :index, :address, :storage_keys)`
)
