package eth_headers

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.header_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.header_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.header_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	hModels, ok := models.([]HeaderModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]HeaderModelV3), models)
	}
	for _, hModel := range hModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), hModel.BlockNumber, hModel.BlockHash, hModel.ParentHash, hModel.CID,
			hModel.TotalDifficulty, hModel.NodeID, hModel.Reward, hModel.StateRoot, hModel.TxRoot, hModel.RctRoot, hModel.UncleRoot,
			hModel.Bloom, hModel.Timestamp, hModel.MhKey, hModel.TimesValidated, hModel.Coinbase); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
