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

package csv

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/vulcanize/migration-tools/pkg/eth_access_lists"
	"github.com/vulcanize/migration-tools/pkg/eth_accounts"
	"github.com/vulcanize/migration-tools/pkg/eth_headers"
	"github.com/vulcanize/migration-tools/pkg/eth_logs"
	"github.com/vulcanize/migration-tools/pkg/eth_receipts"
	"github.com/vulcanize/migration-tools/pkg/eth_state"
	"github.com/vulcanize/migration-tools/pkg/eth_storage"
	"github.com/vulcanize/migration-tools/pkg/eth_transactions"
	"github.com/vulcanize/migration-tools/pkg/eth_uncles"
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
	"github.com/vulcanize/migration-tools/pkg/public_nodes"
)

type Writer interface {
	Write(pgStr WriteCSVStr, models interface{}) error
	io.Closer
}

// Writer struct for writing v3 DB public.nodes models
type writer struct {
	dsts map[WriteCSVStr]io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for public.nodes
func NewWriter(dsts map[WriteCSVStr]io.WriteCloser) Writer {
	return &writer{dsts: dsts}
}

// Write satisfies interfaces.Writer for v3 database
func (w *writer) Write(pgStr WriteCSVStr, models interface{}) error {
	dst := w.dsts[pgStr]
	vs := reflect.ValueOf(models)
	if vs.Kind() == reflect.Slice {
		for i := 0; i < vs.Len(); i++ {
			v := vs.Index(i)
			fields := make([]interface{}, v.NumField())
			for j := 0; j < v.NumField(); j++ {
				fields[j] = v.Field(j).Interface()
			}
			if _, err := fmt.Fprintf(dst, string(pgStr), fields...); err != nil {
				return err
			}
		}
	} else {
		return errors.New("expected models to be a Slice")
	}
	return nil
}

// WriteAlt satisfies interfaces.Writer for v3 database
// this uses a type switch rather than reflect to iterate over all the fields
func (w *writer) WriteAlt(pgStr WriteCSVStr, models interface{}) error {
	dst := w.dsts[pgStr]
	switch vs := models.(type) {
	case []public_blocks.IPLDModel:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.Key, v.Data); err != nil {
				return err
			}
		}
	case []public_nodes.NodeModel:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.GenesisBlock, v.NetworkID, v.NodeID,
				v.ClientName, v.ChainID); err != nil {
				return err
			}
		}
	case []eth_headers.HeaderModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.BlockNumber, v.BlockHash, v.ParentHash, v.CID,
				v.TotalDifficulty, v.NodeID, v.Reward, v.StateRoot, v.TxRoot, v.RctRoot, v.UncleRoot, v.Bloom,
				v.Timestamp, v.MhKey, v.TimesValidated, v.Coinbase); err != nil {
				return err
			}
		}
	case []eth_uncles.UncleModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.BlockHash, v.HeaderID, v.ParentHash, v.CID,
				v.Reward, v.MhKey); err != nil {
				return err
			}
		}
	case []eth_transactions.TransactionModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.HeaderID, v.TxHash, v.CID, v.Dst, v.Src, v.Index,
				v.MhKey, v.Data, v.Type, v.Value); err != nil {
				return err
			}
		}
	case []eth_access_lists.AccessListElementModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.TxID, v.Index, v.Address, v.StorageKeys); err != nil {
				return err
			}
		}
	case []eth_receipts.ReceiptModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.TxID, v.LeafCID, v.Contract, v.ContractHash, v.LeafMhKey,
				v.PostState, v.PostStatus, v.LogRoot); err != nil {
				return err
			}
		}
	case []eth_logs.LogModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.LeafCID, v.LeafMhKey, v.ReceiptID, v.Address, v.Index,
				v.Topic0, v.Topic1, v.Topic2, v.Topic3, v.Data); err != nil {
				return err
			}
		}
	case []eth_state.StateModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.HeaderID, v.StateKey, v.CID, v.Path, v.NodeType, v.Diff,
				v.MhKey); err != nil {
				return err
			}
		}
	case []eth_accounts.AccountModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.HeaderID, v.StatePath, v.Balance, v.Nonce, v.CodeHash,
				v.StorageRoot); err != nil {
				return err
			}
		}
	case []eth_storage.StorageModelV3:
		for _, v := range vs {
			if _, err := fmt.Fprintf(dst, string(pgStr), v.HeaderID, v.StatePath, v.StorageKey, v.CID, v.Path,
				v.NodeType, v.Diff, v.MhKey); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unrecognized model type: %T", models)
	}
	return nil
}

// Close satisfies io.Closer
func (w *writer) Close() error {
	for _, dst := range w.dsts {
		if err := dst.Close(); err != nil {
			return err
		}
	}
	return nil
}
