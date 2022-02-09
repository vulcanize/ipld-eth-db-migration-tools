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

// HeaderModelV2WithMeta is the db model for eth.header_cids for v2 DB with the additional data required to convert to the v3 model
type HeaderModelV2WithMeta struct {
	IPLD []byte `db:"data"`
	HeaderModelV2
}

// HeaderModelV2 is the db model for eth.header_cids for v2 DB
type HeaderModelV2 struct {
	ID              int64  `db:"id"`
	BlockNumber     string `db:"block_number"`
	BlockHash       string `db:"block_hash"`
	ParentHash      string `db:"parent_hash"`
	CID             string `db:"cid"`
	MhKey           string `db:"mh_key"`
	TotalDifficulty string `db:"td"`
	NodeID          string `db:"node_id"`
	Reward          string `db:"reward"`
	StateRoot       string `db:"state_root"`
	UncleRoot       string `db:"uncle_root"`
	TxRoot          string `db:"tx_root"`
	RctRoot         string `db:"receipt_root"`
	Bloom           []byte `db:"bloom"`
	Timestamp       uint64 `db:"timestamp"`
	TimesValidated  int64  `db:"times_validated"`
	BaseFee         *int64 `db:"base_fee"`
}

// HeaderModelV3 is the db model for eth.header_cids for v3 DB
type HeaderModelV3 struct {
	BlockNumber     string `db:"block_number"`
	BlockHash       string `db:"block_hash"`
	ParentHash      string `db:"parent_hash"`
	CID             string `db:"cid"`
	TotalDifficulty string `db:"td"`
	NodeID          string `db:"node_id"`
	Reward          string `db:"reward"`
	StateRoot       string `db:"state_root"`
	TxRoot          string `db:"tx_root"`
	RctRoot         string `db:"receipt_root"`
	UncleRoot       string `db:"uncle_root"`
	Bloom           []byte `db:"bloom"`
	Timestamp       uint64 `db:"timestamp"`
	MhKey           string `db:"mh_key"`
	TimesValidated  int64  `db:"times_validated"`
	Coinbase        string `db:"coinbase"`
}
