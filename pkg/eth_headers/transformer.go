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

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Transformer struct for transforming v2 DB eth.header_cids models to v3 DB models
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for eth.header_cids
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for eth.header_cids
func (t *Transformer) Transform(models interface{}, expectedRange [2]uint64) (interface{}, [][2]uint64, error) {
	v2Models, ok := models.([]HeaderModelV2WithMeta)
	if !ok {
		return nil, [][2]uint64{expectedRange}, fmt.Errorf("expected models of type %T, got %T", make([]HeaderModelV2WithMeta, 0), v2Models)
	}
	v3Models := make([]HeaderModelV3, len(v2Models))
	for i, model := range v2Models {
		header := new(types.Header)
		if err := rlp.DecodeBytes(model.IPLD, header); err != nil {
			return nil, [][2]uint64{expectedRange}, err
		}
		v3Models[i] = HeaderModelV3{
			BlockNumber:     model.BlockNumber,
			BlockHash:       model.BlockHash,
			ParentHash:      model.ParentHash,
			CID:             model.CID,
			MhKey:           model.MhKey,
			TotalDifficulty: model.TotalDifficulty,
			NodeID:          model.NodeID,
			Reward:          model.Reward,
			StateRoot:       model.StateRoot,
			UncleRoot:       model.UncleRoot,
			TxRoot:          model.TxRoot,
			RctRoot:         model.RctRoot,
			Bloom:           model.Bloom,
			Timestamp:       model.Timestamp,
			TimesValidated:  model.TimesValidated,
			Coinbase:        header.Coinbase.String(),
		}
	}
	return v3Models, nil, nil
}
