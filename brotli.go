package brotli

import (
	"bytes"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/sohaha/zlsgo/znet"
)

type Config struct {
	// CompressionLevel brotli compression level to use
	CompressionLevel int
	// PoolMaxSize maximum number of resource pools
	PoolMaxSize int
	// MinContentLength minimum content length to trigger brotli, the unit is in byte.
	MinContentLength int
}

func Default() znet.HandlerFunc {
	return New(Config{
		CompressionLevel: 7,
		PoolMaxSize:      1024,
		MinContentLength: 1024,
	})
}

func New(conf Config) znet.HandlerFunc {
	pool := &poolCap{
		c: make(chan *brotli.Writer, conf.PoolMaxSize),
		l: conf.CompressionLevel,
	}
	return func(c *znet.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "br") {
			c.Next()
			return
		}
		c.Next()
		p := c.PrevContent()

		if len(p.Content) < conf.MinContentLength {
			return
		}

		if encoding := c.GetSetHeader("Content-Encoding"); len(encoding) > 0 {
			return
		}

		g := pool.Get()
		defer pool.Put(g)

		be := &bytes.Buffer{}
		g.Reset(be)

		_, err := g.Write(p.Content)
		if err != nil {
			return
		}
		_ = g.Flush()

		c.Byte(p.Code, be.Bytes())
		c.SetHeader("Content-Encoding", "br")
	}
}
