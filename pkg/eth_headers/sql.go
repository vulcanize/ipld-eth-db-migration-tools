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

package eth_headers

const (
	PgReadEthHeadersStr = `SELECT public.blocks.data, eth.header_cids.*
							FROM eth.header_cids
							INNER JOIN public.blocks ON (header_cids.mh_key = blocks.key)
							WHERE block_number = $1`

	PgWriteEthHeadersStr = `INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, mh_key, td, node_id,
							reward, state_root, uncle_root, tx_root, receipt_root, bloom, timestamp, times_validated, coinbase)
							VALUES (unnest($1::BIGINT[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR(66)[]), unnest($4::TEXT[]),
							unnest($5::TEXT[]), unnest($6::NUMERIC[]), unnest($7::VARCHAR(128)[]), unnest($8::NUMERIC[]),
							unnest($9::VARCHAR(66)[]), unnest($10::VARCHAR(66)[]), unnest($11::VARCHAR(66)[]), unnest($12::VARCHAR(66)[]),
							unnest($13::BYTEA[]), unnest($14::BIGINT[]), unnest($15::INTEGER[]), unnest($16::VARCHAR(66)[]))`
)
