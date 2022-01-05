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

package eth_accounts

const (
	PgReadEthAccountsStr = `SELECT eth.header_cids.block_hash, eth.state_cids.state_path, eth.state_accounts.*
							FROM eth.state_accounts
							INNER JOIN eth.state_cids ON (state_accounts.state_id = state_cids.id)
							INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
							WHERE block_number = $1`

	PgWriteEthAccountsStr = `INSERT INTO eth.state_accounts (header_id, state_path, balance, nonce, code_hash, storage_root)
							VALUES (unnest($1::VARCHAR(66)[]), unnest($2::BYTEA[]), unnest($3::NUMERIC[]), unnest($4::BIGINT[]),
							unnest($5::BYTEA[]), unnest($6::VARCHAR(66)[]))`
)
