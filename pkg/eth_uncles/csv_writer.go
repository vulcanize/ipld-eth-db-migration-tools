package eth_uncles

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB eth.uncle_cids models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for eth.uncle_cids
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for eth.uncle_cids
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	uModels, ok := models.([]UncleModelV3)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]UncleModelV3), models)
	}
	for _, uModel := range uModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), uModel.BlockHash, uModel.HeaderID, uModel.ParentHash,
			uModel.CID, uModel.Reward, uModel.MhKey); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
