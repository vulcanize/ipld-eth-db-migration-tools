package eth_logs

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.log_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.log_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.log_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	lModels, ok := models.([]LogModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]LogModelV3), models)
	}
	for _, lModel := range lModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), lModel.LeafCID, lModel.LeafMhKey, lModel.ReceiptID,
			lModel.Address, lModel.Index, lModel.Topic0, lModel.Topic1, lModel.Topic2, lModel.Topic3, lModel.Data); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
