package utils

import (
	"compress/flate"
	"io"
)

// NewMaxCompressionWriter returns a WriteCloser that writes to w with maximum compression level
func NewMaxCompressionWriter(w io.Writer) io.WriteCloser {
	fw, err := flate.NewWriter(w, flate.BestCompression)
	if err != nil {
		// If there's an error, fallback to a default writer
		// This should never happen with BestCompression, but we handle it just in case
		fw, _ = flate.NewWriter(w, flate.DefaultCompression)
	}
	return fw
} 