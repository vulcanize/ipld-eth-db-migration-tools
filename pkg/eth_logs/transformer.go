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

import (
	"fmt"

	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Transformer struct for transforming v2 DB eth.log_cids models to v3 DB models
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for eth.log_cids
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for eth.log_cids
func (t *Transformer) Transform(models interface{}, expectedRange [2]uint64) (interface{}, [][2]uint64, error) {
	v2Models, ok := models.(*[]LogModelV2WithMeta)
	if !ok {
		return nil, [][2]uint64{expectedRange}, fmt.Errorf("expected models of type %T, got %T", new([]LogModelV2WithMeta), models)
	}
	v3Models := make([]LogModelV3, len(*v2Models))
	for i, model := range *v2Models {
		v3Models[i] = LogModelV3{
			ReceiptID: model.TxHash,
			LeafCID:   model.LeafCID,
			LeafMhKey: model.LeafMhKey,
			Address:   model.Address,
			Index:     model.Index,
			Data:      model.Data,
			Topic0:    model.Topic0,
			Topic1:    model.Topic1,
			Topic2:    model.Topic2,
			Topic3:    model.Topic3,
		}
	}
	return v3Models, nil, nil
}
