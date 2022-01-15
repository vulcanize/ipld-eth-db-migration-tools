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

package eth_accounts

import (
	"fmt"
	"strconv"

	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Transformer struct for transforming v2 DB eth.state_accounts models into v3 DB models
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for eth.state_accounts
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for eth.state_accounts
func (t *Transformer) Transform(models interface{}, expectedRange [2]uint64) (interface{}, [][2]uint64, error) {
	v2Models, ok := models.(*[]AccountModelV2WithMeta)
	if !ok {
		return nil, [][2]uint64{expectedRange}, fmt.Errorf("expected models of type %T, got %T", new([]AccountModelV2WithMeta), models)
	}
	v3Models := make([]AccountModelV3, len(*v2Models))
	expectedHeight := expectedRange[0]
	missingHeights := make([][2]uint64, 0)
	for i, model := range *v2Models {
		height, err := strconv.ParseUint(model.BlockNumber, 10, 64)
		if err != nil {
			return nil, [][2]uint64{expectedRange}, fmt.Errorf("EthAccount transformer unable to parse blocknumber %s", model.BlockHash)
		}
		// if the expected height doesn't match the actual current block height, we have a gap between the two
		if expectedHeight < height {
			missingHeights = append(missingHeights, [2]uint64{expectedHeight, height - 1})
			expectedHeight = height
		} else if height < expectedHeight {
			return nil, [][2]uint64{expectedRange}, fmt.Errorf("it should not be possible for the current"+
				"expected height (%d) to be greater than the actual current height (%d)", expectedHeight, height)
		}
		v3Models[i] = AccountModelV3{
			HeaderID:    model.BlockHash,
			StatePath:   model.StatePath,
			Balance:     model.Balance,
			Nonce:       model.Nonce,
			CodeHash:    model.CodeHash,
			StorageRoot: model.StorageRoot,
		}
		expectedHeight++
	}
	// if the last processed height isn't the last block in the range, we have a gap at the end of the range
	if expectedHeight-1 != expectedRange[1] {
		missingHeights = append(missingHeights, [2]uint64{expectedHeight, expectedRange[1]})
	}
	return v3Models, missingHeights, nil
}
