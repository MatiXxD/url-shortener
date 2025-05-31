package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter used instead of http.ResponseWriter in compress middleware
type compressWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	gzipWriter, _ := gzip.NewWriterLevel(w, gzip.BestCompression) // ignore error because using fixed lvl
	return &compressWriter{
		ResponseWriter: w,
		gw:             gzipWriter,
	}
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	return cw.gw.Write(b)
}

func (cw *compressWriter) Close() error {
	return cw.gw.Close()
}

// compressReader used instead of default request's body in compress middleware
type compressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		gr: gz,
	}, nil
}

func (cr *compressReader) Read(b []byte) (int, error) {
	return cr.gr.Read(b)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}
	return cr.gr.Close()
}

// CompressMiddleware encode/decode with gzip
func CompressMiddleware(h http.Handler) http.Handler {
	cf := func(w http.ResponseWriter, r *http.Request) {
		cw := w

		// compress response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") &&
			shouldEncode(r.Header.Get("Content-Type")) {
			tmp := newCompressWriter(w)
			defer tmp.Close()

			cw = tmp
		}

		// decompress request body
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "failed to decompress body", http.StatusInternalServerError) // can't decompress body -> 500
				return
			}
			defer cr.Close()

			r.Body = cr
		}

		h.ServeHTTP(cw, r)
	}

	return http.HandlerFunc(cf)
}

// encodeList lists of type to encode
var encodeList = []string{
	"application/json",
	"text/html",
}

func shouldEncode(contentType string) bool {
	for _, ct := range encodeList {
		if strings.Contains(ct, contentType) {
			return true
		}
	}
	return false
}
