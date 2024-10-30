package middleware

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func Decompress() gin.HandlerFunc {
	decompressor := NewDecompressor()

	return decompressor.Handler
}

type DecodeFunc func(io.Reader) io.ReadCloser

type Decompressor struct {
	pooledDecoders map[string]*sync.Pool
}

func NewDecompressor() *Decompressor {
	d := &Decompressor{
		pooledDecoders: make(map[string]*sync.Pool),
	}

	d.SetDecoder("deflate", DecoderDeflate)
	d.SetDecoder("gzip", DecoderGzip)

	return d
}

func (d *Decompressor) SetDecoder(encoding string, fn DecodeFunc) {
	encoding = strings.ToLower(encoding)

	delete(d.pooledDecoders, encoding)

	if fn(nil) != nil {
		pool := &sync.Pool{
			New: func() interface{} {
				return fn
			},
		}
		d.pooledDecoders[encoding] = pool
	}
}

func (d *Decompressor) selectDecoder(h http.Header, r io.ReadCloser) (io.ReadCloser, string) {
	encoded := h.Get("Content-Encoding")

	// content is not encoded
	if encoded == "" {
		return r, ""
	}

	// try to get from pooledDecoders
	for name := range d.pooledDecoders {
		if name == encoded {
			if pool, ok := d.pooledDecoders[name]; ok {
				if decoder, ok := pool.Get().(DecodeFunc); ok {
					return decoder(r), encoded
				}
			}
		}
	}

	return nil, encoded
}

func (d *Decompressor) Handler(c *gin.Context) {
	decoder, _ := d.selectDecoder(c.Request.Header, c.Request.Body)

	if decoder == nil {
		c.String(http.StatusInternalServerError, "Content could not be decoded")

		return
	}

	c.Request.Body = decoder
	defer decoder.Close()
	c.Next()
}
