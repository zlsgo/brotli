package brotli

import (
	"io/ioutil"

	"github.com/andybalholm/brotli"
)

type (
	poolCap struct {
		c chan *brotli.Writer
		l int
	}
)

func (bp *poolCap) Get() (g *brotli.Writer) {
	select {
	case g = <-bp.c:
	default:
		g = brotli.NewWriterLevel(ioutil.Discard, bp.l)
	}

	return
}

func (bp *poolCap) Put(g *brotli.Writer) {
	select {
	case bp.c <- g:
	default:
	}
}
