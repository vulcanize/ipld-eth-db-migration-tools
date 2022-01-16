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

package repair

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff/indexer/ipld"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dshelp "github.com/ipfs/go-ipfs-ds-help"
	"github.com/multiformats/go-multihash"

	"github.com/vulcanize/migration-tools/pkg/eth_logs"
	"github.com/vulcanize/migration-tools/pkg/interfaces"
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
)

// Transformer struct for repairing log_cids models in the v3 DB
type Transformer struct {
}

// NewTransformer satisfies interfaces.TransformerConstructor for repairing log_cids
func NewTransformer() interfaces.Transformer {
	return &Transformer{}
}

// Transform satisfies interfaces.Transformer for eth.log_cids
func (t *Transformer) Transform(models interface{}, expectedRange [2]uint64) (interface{}, [][2]uint64, error) {
	v3Models, ok := models.(*[]eth_logs.LogModelV3)
	if !ok {
		return nil, [][2]uint64{expectedRange}, fmt.Errorf("expected models of type %T, got %T", new([]eth_logs.LogModelV3), models)
	}
	missingIPLDs := make([]public_blocks.IPLDModel, len(*v3Models))
	for i, model := range *v3Models {
		log := new(types.Log)
		log.Address = common.HexToAddress(model.Address)
		log.Data = model.Data
		log.Topics = make([]common.Hash, 0, 4)
		if model.Topic0 != "" {
			log.Topics = append(log.Topics, common.HexToHash(model.Topic0))
		}
		if model.Topic1 != "" {
			log.Topics = append(log.Topics, common.HexToHash(model.Topic1))
		}
		if model.Topic2 != "" {
			log.Topics = append(log.Topics, common.HexToHash(model.Topic2))
		}
		if model.Topic3 != "" {
			log.Topics = append(log.Topics, common.HexToHash(model.Topic3))
		}
		data, key, err := rlpAndBlockStoreKey(log)
		if err != nil {
			return nil, nil, err
		}
		if model.LeafMhKey != key {
			return nil, [][2]uint64{expectedRange}, fmt.Errorf("log_cids record mh key (%s) does not match dervied mh key (%s)", model.LeafMhKey, key)
		}
		missingIPLDs[i] = public_blocks.IPLDModel{
			Key:  key,
			Data: data,
		}
	}
	return missingIPLDs, nil, nil
}

func rlpAndBlockStoreKey(log *types.Log) ([]byte, string, error) {
	logRaw, err := rlp.EncodeToBytes(log)
	if err != nil {
		return nil, "", err
	}
	c, err := ipld.RawdataToCid(ipld.MEthLog, logRaw, multihash.KECCAK_256)
	if err != nil {
		return nil, "", err
	}
	return logRaw, blockstore.BlockPrefix.String() + dshelp.MultihashToDsKey(c.Hash()).String(), nil
}
