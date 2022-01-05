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

package eth_state

const (
	PgReadEthStateStr = `SELECT eth.header_cids.block_hash, eth.state_cids.*
						FROM eth.state_cids
						INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
						WHERE block_number = $1`

	PgWriteEthStateStr = `INSERT INTO eth.state_cids (header_id, state_path, state_leaf_key, node_type, cid, mh_key, diff)
						VALUES (unnest($1::VARCHAR(66)[]), unnest($2::BYTEA[]), unnest($3::VARCHAR(66)[]),
						unnest($4::INTEGER[]), unnest($5::TEXT[]), unnest($6::TEXT[]), unnest($7::BOOLEAN[]))`
)
