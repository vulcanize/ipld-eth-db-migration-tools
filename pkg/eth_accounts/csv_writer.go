package eth_accounts

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.state_accounts models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.state_accounts
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies csv.Writer for v3 database for eth.state_accounts
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	aModels, ok := models.([]AccountModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]AccountModelV3), models)
	}
	for _, aModel := range aModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), aModel.HeaderID, aModel.StatePath, aModel.Balance, aModel.Nonce,
			aModel.CodeHash, aModel.StorageRoot); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
