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

package eth_storage

// StorageModelV2WithMeta is the db model for eth.storage_cids for v2 DB
// with the additional metadata required to convert to the v3 model
type StorageModelV2WithMeta struct {
	BlockHash string `db:"block_hash"`
	StatePath []byte `db:"state_path"`
	StorageModelV2
}

// StorageModelV2 is the db model for eth.storage_cids for v2 DB
type StorageModelV2 struct {
	ID         int64  `db:"id"`
	StateID    int64  `db:"state_id"`
	Path       []byte `db:"storage_path"`
	StorageKey string `db:"storage_leaf_key"`
	NodeType   int    `db:"node_type"`
	CID        string `db:"cid"`
	MhKey      string `db:"mh_key"`
	Diff       bool   `db:"diff"`
}

// StorageModelV3 is the db model for eth.storage_cids for v3 DB
type StorageModelV3 struct {
	HeaderID   string `db:"header_id"`
	StatePath  []byte `db:"state_path"`
	StorageKey string `db:"storage_leaf_key"`
	CID        string `db:"cid"`
	Path       []byte `db:"storage_path"`
	NodeType   int    `db:"node_type"`
	Diff       bool   `db:"diff"`
	MhKey      string `db:"mh_key"`
}
