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

package eth_receipts

const (
	PgReadEthReceiptsStr = `SELECT eth.transaction_cids.tx_hash, eth.receipt_cids.*
							FROM eth.receipt_cids
							INNER JOIN eth.transaction_cids ON (receipt_cids.tx_id = transaction_cids.id)
							INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
							WHERE block_number = $1`

	PgWriteEthReceiptsStr = `INSERT INTO eth.receipt_cids (tx_id, leaf_cid, leaf_mh_key, post_status, post_state,
								contract, contract_hash, log_root)
								VALUES (unnest($1::VARCHAR(66)[]), unnest($2::TEXT[]), unnest($3::TEXT[]),
								unnest($4::INTEGER[]), unnest($5::VARCHAR(66)[]), unnest($6::VARCHAR(66)[]),
								unnest($7::VARCHAR(66)[]), unnest($8::VARCHAR(66)[]))`
)
