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

package public_nodes

// NodeModel is the v2 and v3 model for public.nodes
type NodeModel struct {
	GenesisBlock string `db:"genesis_block"`
	NetworkID    string `db:"network_id"`
	NodeID       string `db:"node_id"`
	ClientName   string `db:"client_name"`
	ChainID      int    `db:"chain_id"`
}
