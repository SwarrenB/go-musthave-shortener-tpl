package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

func EncoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}
	return gw
}

func EncoderDeflate(w io.Writer, level int) io.Writer {
	dw, err := flate.NewWriter(w, level)
	if err != nil {
		return nil
	}

	return dw
}

func DecoderGzip(r io.Reader) io.ReadCloser {
	if r == nil {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	dr, err := gzip.NewReader(r)
	if err != nil {
		return nil
	}
	return dr
}

func DecoderDeflate(r io.Reader) io.ReadCloser {
	return flate.NewReader(r)
}
