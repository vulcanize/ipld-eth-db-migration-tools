package eth_transactions

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.transaction_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.transaction_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.transaction_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	tModels, ok := models.([]TransactionModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]TransactionModelV3), models)
	}
	for _, tModel := range tModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), tModel.HeaderID, tModel.TxHash, tModel.CID, tModel.Dst,
			tModel.Src, tModel.Index, tModel.MhKey, tModel.Data, tModel.Type, tModel.Value); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
