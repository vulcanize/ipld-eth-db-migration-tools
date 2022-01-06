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

package eth_transactions

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Transformer struct for transforming v2 DB eth.transaction_cids models to v3 DB models
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for eth.transaction_cids
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for eth.transaction_cids
func (t *Transformer) Transform(models interface{}, expectedRange [2]uint64) (interface{}, [][2]uint64, error) {
	v2Models, ok := models.([]TransactionModelV2WithMeta)
	if !ok {
		return nil, [][2]uint64{expectedRange}, fmt.Errorf("expected models of type %T, got %T", make([]TransactionModelV2WithMeta, 0), v2Models)
	}
	v3Models := make([]TransactionModelV3, len(v2Models))
	for i, model := range v2Models {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(model.IPLD); err != nil {
			return nil, [][2]uint64{expectedRange}, err
		}
		var val string
		if tx.Value() != nil {
			val = tx.Value().String()
		}
		v3Models[i] = TransactionModelV3{
			HeaderID: model.BlockHash,
			Index:    model.Index,
			TxHash:   model.TxHash,
			CID:      model.CID,
			MhKey:    model.MhKey,
			Dst:      model.Dst,
			Src:      model.Src,
			Data:     model.Data,
			Type:     *model.Type,
			Value:    val,
		}
	}
	return v3Models, nil, nil
}
