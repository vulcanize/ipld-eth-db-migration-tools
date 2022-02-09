package eth_access_lists

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.access_list_elements models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.access_list_elements
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies csv.Writer for v3 database for eth.access_list_elements
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	alModels, ok := models.([]AccessListElementModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]AccessListElementModelV3), models)
	}
	for _, alModel := range alModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), alModel.TxID, alModel.Index, alModel.Address,
			alModel.StorageKeys); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
