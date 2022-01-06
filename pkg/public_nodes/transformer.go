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

import "github.com/vulcanize/migration-tools/pkg/interfaces"

// Transformer for transforming v2 DB public.nodes models to v3 DB models
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for public.nodes
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for public.nodes
func (t *Transformer) Transform(models interface{}, _ [2]uint64) (interface{}, [][2]uint64, error) {
	return models, nil, nil
}
