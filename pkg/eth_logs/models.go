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

// LogModelV2WithMeta is the db model for eth.logs for v2 DB with the additional metadata required to convert to v3 model
type LogModelV2WithMeta struct {
	TxHash string `db:"tx_hash"`
	LogModelV2
}

// LogModelV2 is the db model for eth.logs for v2 DB
type LogModelV2 struct {
	ID        int64  `db:"id"`
	LeafCID   string `db:"leaf_cid"`
	LeafMhKey string `db:"leaf_mh_key"`
	ReceiptID int64  `db:"receipt_id"`
	Address   string `db:"address"`
	Index     int64  `db:"index"`
	Data      []byte `db:"log_data"`
	Topic0    string `db:"topic0"`
	Topic1    string `db:"topic1"`
	Topic2    string `db:"topic2"`
	Topic3    string `db:"topic3"`
}

// LogModelV3 is the db model for eth.logs for v3 DB
type LogModelV3 struct {
	LeafCID   string `db:"leaf_cid"`
	LeafMhKey string `db:"leaf_mh_key"`
	ReceiptID string `db:"rct_id"`
	Address   string `db:"address"`
	Index     int64  `db:"index"`
	Topic0    string `db:"topic0"`
	Topic1    string `db:"topic1"`
	Topic2    string `db:"topic2"`
	Topic3    string `db:"topic3"`
	Data      []byte `db:"log_data"`
}
