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

package public_blocks

// TransactionModelV2 is the db model for eth.transaction_cids for v2 DB
type TransactionModelV2 struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	Index    int64  `db:"index"`
	TxHash   string `db:"tx_hash"`
	CID      string `db:"cid"`
	MhKey    string `db:"mh_key"`
	Dst      string `db:"dst"`
	Src      string `db:"src"`
	Data     []byte `db:"tx_data"`
	Type     *uint8 `db:"tx_type"`
}

// TransactionModelV3 is the db model for eth.transaction_cids for v3 DB
type TransactionModelV3 struct {
	HeaderID string `db:"header_id"`
	Index    int64  `db:"index"`
	TxHash   string `db:"tx_hash"`
	CID      string `db:"cid"`
	MhKey    string `db:"mh_key"`
	Dst      string `db:"dst"`
	Src      string `db:"src"`
	Data     []byte `db:"tx_data"`
	Type     uint8  `db:"tx_type"`
	Value    string `db:"value"`
}
