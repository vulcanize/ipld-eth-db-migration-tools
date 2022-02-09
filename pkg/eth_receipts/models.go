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

// ReceiptModelV2WithMeta is the db model for eth.receipt_cids for v2 DB
// with the additional metadata required to convert to the v3 model
type ReceiptModelV2WithMeta struct {
	TxHash string `db:"tx_hash"`
	ReceiptModelV2
}

// ReceiptModelV2 is the db model for eth.receipt_cids for v2 DB
type ReceiptModelV2 struct {
	ID           int64  `db:"id"`
	TxID         int64  `db:"tx_id"`
	LeafCID      string `db:"leaf_cid"`
	LeafMhKey    string `db:"leaf_mh_key"`
	PostStatus   uint64 `db:"post_status"`
	PostState    string `db:"post_state"`
	Contract     string `db:"contract"`
	ContractHash string `db:"contract_hash"`
	LogRoot      string `db:"log_root"`
}

// ReceiptModelV3 is the db model for eth.receipt_cids for v3 DB
type ReceiptModelV3 struct {
	TxID         string `db:"tx_id"`
	LeafCID      string `db:"leaf_cid"`
	Contract     string `db:"contract"`
	ContractHash string `db:"contract_hash"`
	LeafMhKey    string `db:"leaf_mh_key"`
	PostState    string `db:"post_state"`
	PostStatus   uint64 `db:"post_status"`
	LogRoot      string `db:"log_root"`
}
