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

package eth_uncles

const (
	PgReadEthUnclesStr = `SELECT eth.header_cids.block_hash, eth.uncle_cids.*
							FROM eth.uncle_cids
							INNER JOIN eth.header_cids ON (uncle_cids.header_id = header_cids.id)
							WHERE block_number = $1`

	PgWriteEthUnclesStr = `INSERT INTO eth.uncle_cids (header_id, block_hash, parent_hash, cid, mh_key, reward)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR(66)[]),
							unnest($4::TEXT[]), unnest($5::TEXT[]), unnest($6::NUMERIC[]))`
)
