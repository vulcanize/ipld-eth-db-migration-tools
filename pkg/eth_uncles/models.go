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

// UncleModelV2WithMeta is the db model for eth.uncle_cids for v2 DB
// with the additional metadata required to convert to the v3 model
type UncleModelV2WithMeta struct {
	HeaderHash string `db:"header_cids.block_hash"`
	UncleModelV2
}

// UncleModelV2 is the db model for eth.uncle_cids for v2 DB
type UncleModelV2 struct {
	ID         int64  `db:"id"`
	HeaderID   int64  `db:"header_id"`
	BlockHash  string `db:"uncle_cids.block_hash"`
	ParentHash string `db:"parent_hash"`
	CID        string `db:"cid"`
	MhKey      string `db:"mh_key"`
	Reward     string `db:"reward"`
}

// UncleModelV3 is the db model for eth.uncle_cids for v3 DB
type UncleModelV3 struct {
	HeaderID   string `db:"header_id"`
	BlockHash  string `db:"block_hash"`
	ParentHash string `db:"parent_hash"`
	CID        string `db:"cid"`
	MhKey      string `db:"mh_key"`
	Reward     string `db:"reward"`
}
