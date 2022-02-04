package eth_state

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.state_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.state_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.state_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	sModels, ok := models.([]StateModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]StateModelV3), models)
	}
	for _, sModel := range sModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), sModel.HeaderID, sModel.StateKey, sModel.CID, sModel.Path,
			sModel.NodeType, sModel.Diff, sModel.MhKey); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
