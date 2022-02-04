package public_nodes

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB public.nodes models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for public.nodes
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for nModel
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	nModels, ok := models.([]NodeModel)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]NodeModel), models)
	}
	for _, nModel := range nModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), nModel.GenesisBlock, nModel.NetworkID, nModel.NodeID,
			nModel.ClientName, nModel.ChainID); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
