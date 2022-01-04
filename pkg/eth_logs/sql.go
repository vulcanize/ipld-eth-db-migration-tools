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

package eth_logs

const (
	PgReadEthLogsStr = `SELECT eth.transaction_cis.tx_hash, eth.log_cids.*
						FROM eth.log_cids
						INNER JOIN eth.receipt_cids ON (log_cids.receipt_id = receipt_cids.id)
						INNER JOIN eth.transaction_cids ON (receipt_cids.tx_id = transaction_cids.id)
						INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
						WHERE block_number = $1`

	PgWriteEthLogsStr = `INSERT INTO eth.logs_cids (rct_id, leaf_cid, leaf_mh_key, address, index, log_data, topic0,
							topic1, topic2, topic3)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::TEXT[]), unnest($3::TEXT[]), unnest($4::VARCHAR(66)[]),
							unnest($5::INTEGER[]), unnest($6::BYTEA[]), unnest($7::VARCHAR(66)[]), unnest($8::VARCHAR(66)[]),
							unnest($9::VARCHAR(66)[]), unnest($10::VARCHAR(66)[]))`
)
