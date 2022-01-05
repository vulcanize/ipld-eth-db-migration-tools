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

const (
	PgReadNodesStr = `SELECT client_name, genesis_block, network_id, node_id, nodes_chain_id
						FROM public.nodes`

	PgWriteNodesStr = `INSERT INTO public.nodes (client_name, genesis_block, network_id, node_id, chain_id)
						VALUES (unnest($1::VARCHAR[]), unnest($2::VARCHAR(66)[]), unnest($3::VARCHAR[]),
						unnest($4::VARCHAR(128)[]), unnest($5::INTEGER[]))`
)
