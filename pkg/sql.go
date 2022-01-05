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

type ReadPgStr string

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
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::BYTEA[]), unnest($3::BYTEA[]),
							unnest($4::VARCHAR(66)[]), unnest($5::INTEGER[]), unnest($6::TEXT[]), unnest($7::TEXT[]),
							unnest($8::BOOLEAN[]))`

	PgWriteEthStateStr WritePgStr = `INSERT INTO eth.state_cids (header_id, state_path, state_leaf_key, node_type, cid, mh_key, diff)
						VALUES (unnest($1::VARCHAR(66)[]), unnest($2::BYTEA[]), unnest($3::VARCHAR(66)[]),
						unnest($4::INTEGER[]), unnest($5::TEXT[]), unnest($6::TEXT[]), unnest($7::BOOLEAN[]))`

	PgWriteNodesStr WritePgStr = `INSERT INTO public.nodes (client_name, genesis_block, network_id, node_id, chain_id)
						VALUES (unnest($1::VARCHAR[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR[]),
						unnest($4::VARCHAR(128)[]), unnest($5::INTEGER[]))`

	PgWriteEthReceiptsStr WritePgStr = `INSERT INTO eth.receipt_cids (tx_id, leaf_cid, leaf_mh_key, post_status, post_state,
								contract, contract_hash, log_root)
								VALUES (unnest($1::VARCHAR(66)[]), unnest($2::TEXT[]), unnest($3::TEXT[]),
								unnest($4::INTEGER[]), unnest($5::VARCHAR(66)[]), unnest($6::VARCHAR(66)[]),
								unnest($7::VARCHAR(66)[]), unnest($8::VARCHAR(66)[]))`

	PgWriteEthLogsStr WritePgStr = `INSERT INTO eth.logs_cids (rct_id, leaf_cid, leaf_mh_key, address, index, log_data, topic0,
							topic1, topic2, topic3)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::TEXT[]), unnest($3::TEXT[]), unnest($4::VARCHAR(66)[]),
							unnest($5::INTEGER[]), unnest($6::BYTEA[]), unnest($7::VARCHAR(66)[]), unnest($8::VARCHAR(66)[]),
							unnest($9::VARCHAR(66)[]), unnest($10::VARCHAR(66)[]))`

	PgWriteEthHeadersStr WritePgStr = `INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, mh_key, td, node_id,
							reward, state_root, uncle_root, tx_root, receipt_root, bloom, timestamp, times_validated, coinbase)
							VALUES (unnest($1::BIGINT[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR(66)[]), unnest($4::TEXT[]),
							unnest($5::TEXT[]), unnest($6::NUMERIC[]), unnest($7::VARCHAR(128)[]), unnest($8::NUMERIC[]),
							unnest($9::VARCHAR(66)[]), unnest($10::VARCHAR(66)[]), unnest($11::VARCHAR(66)[]), unnest($12::VARCHAR(66)[]),
							unnest($13::BYTEA[]), unnest($14::BIGINT[]), unnest($15::INTEGER[]), unnest($16::VARCHAR(66)[]))`

	PgWriteEthAccountsStr WritePgStr = `INSERT INTO eth.state_accounts (header_id, state_path, balance, nonce, code_hash, storage_root)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::BYTEA[]), unnest($3::NUMERIC[]), unnest($4::BIGINT[]),
							unnest($5::BYTEA[]), unnest($6::VARCHAR(66)[]))`

	PgWriteAccessListElementsStr WritePgStr = `INSERT INTO eth.access_list_elements (tx_id, index, address, storage_keys)
									VALUES (unnest($1::VARCHAR(66)[]), unnest($2::INTEGER[]), unnest($3::VARCHAR(66)[]),
									unnest($4::VARCHAR(66)[][]))`

	PgGapFinderStr = `SELECT s.i AS missing_cmd
						FROM generate_series($1,$2) s(i)
						WHERE NOT EXISTS (SELECT 1 FROM eth.header_cids WHERE block_number = s.i)`
)
