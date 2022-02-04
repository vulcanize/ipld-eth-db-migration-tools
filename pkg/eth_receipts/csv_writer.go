package eth_receipts

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.receipt_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.receipt_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.receipt_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	rModels, ok := models.([]ReceiptModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]ReceiptModelV3), models)
	}
	for _, rModel := range rModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), rModel.TxID, rModel.LeafCID, rModel.Contract,
			rModel.ContractHash, rModel.LeafMhKey, rModel.PostState, rModel.PostStatus, rModel.LogRoot); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
