package public_blocks

import (
	"fmt"
	"io"

	"github.com/vulcanize/migration-tools/pkg/csv"
)

// CSVWriter struct for writing v3 DB public.blocks models to a csv file
type CSVWriter struct {
	dst io.WriteCloser
}

// NewWriter satisfies interfaces.WriterConstructor for public.blocks
func NewWriter(dst io.WriteCloser) *CSVWriter {
	return &CSVWriter{dst: dst}
}

// Write satisfies cshModel.Writer for v3 database for public.blocks
func (cw *CSVWriter) Write(pgStr csv.WriteCSVStr, models interface{}) error {
	iModels, ok := models.([]IPLDModel)
	if !ok {
		return fmt.Errorf("expected models of type %T, got %T", new([]IPLDModel), models)
	}
	for _, iModel := range iModels {
		if _, err := fmt.Fprintf(cw.dst, string(pgStr), iModel.Key, iModel.Data); err != nil {
			return err
		}
	}
	return nil
}

// Close satisfies io.Closer
func (cw *CSVWriter) Close() error {
	return cw.dst.Close()
}
