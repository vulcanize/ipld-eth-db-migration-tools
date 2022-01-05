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

package eth_transactions

const (
	PgReadEthTransactionsStr = `SELECT public.blocks.data, eth.header_cids.block_hash, eth.transaction_cids.*
								FROM eth.transaction_cids
								INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
								INNER JOIN public.blocks ON (transaction_cids.mh_key = blocks.key)
								WHERE block_number = $1`

	PgWriteEthTransactionsStr = `INSERT INTO eth.transaction_cids (header_id, index, tx_hash, cid, mh_key, dst, src,
									tx_data, tx_type, value)
									VALUES (unnest($1::VARCHAR(66)[]), unnest($2::INTEGER[]), unnest($3::VARCHAR(66)[]),
									unnest($4::TEXT[]), unnest($5::TEXT[]), unnest($6::VARCHAR(66)[]),
									unnest($7::VARCHAR(66)[]), unnest($8::BYTEA[]), unnest($9::INTEGER[]),
									unnest($10::NUMERIC[]))`
)
