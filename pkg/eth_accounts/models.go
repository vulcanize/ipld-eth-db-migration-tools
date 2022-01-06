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

// AccountModelV2WithMeta is the db model for eth.state_accounts for v2 DB, with the additional metadata
// required for converting to the v3 DB model
type AccountModelV2WithMeta struct {
	BlockHash   string `db:"block_hash"`
	BlockNumber string `db:"block_number"`
	StatePath   []byte `db:"state_path"`
	AccountModelV2
}

// AccountModelV2 is the db model for an eth state account (decoded value of state leaf node) for v2 DB
type AccountModelV2 struct {
	ID          int64  `db:"id"`
	StateID     int64  `db:"state_id"`
	Balance     string `db:"balance"`
	Nonce       uint64 `db:"nonce"`
	CodeHash    []byte `db:"code_hash"`
	StorageRoot string `db:"storage_root"`
}

// AccountModelV3 is a db model for an eth state account (decoded value of state leaf node) for v3 DB
type AccountModelV3 struct {
	HeaderID    string `db:"header_id"`
	StatePath   []byte `db:"state_path"`
	Balance     string `db:"balance"`
	Nonce       uint64 `db:"nonce"`
	CodeHash    []byte `db:"code_hash"`
	StorageRoot string `db:"storage_root"`
}
