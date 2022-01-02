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

// StateModelV2 is the db model for eth.state_cids for v2 DB
type StateModelV2 struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	Path     []byte `db:"state_path"`
	StateKey string `db:"state_leaf_key"`
	NodeType int    `db:"node_type"`
	CID      string `db:"cid"`
	MhKey    string `db:"mh_key"`
	Diff     bool   `db:"diff"`
}

// StateModelV3 is the db model for eth.state_cids for v3 DB
type StateModelV3 struct {
	HeaderID string `db:"header_id"`
	Path     []byte `db:"state_path"`
	StateKey string `db:"state_leaf_key"`
	NodeType int    `db:"node_type"`
	CID      string `db:"cid"`
	MhKey    string `db:"mh_key"`
	Diff     bool   `db:"diff"`
}
