package brotli_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/zlsgo/brotli"
)

const host = "127.0.0.1"

var (
	r    *znet.Engine
	size = 10000
)

func init() {
	r = znet.New()
	r.SetMode(znet.ProdMode)

	r.GET("/brotli", func(c *znet.Context) {
		c.String(200, zstring.Rand(size, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	}, brotli.Default())

}

func TestGzip(t *testing.T) {
	tt := zlsgo.NewTest(t)

	go func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brotli", nil)
		r.ServeHTTP(w, req)
		tt.Equal(200, w.Code)
		tt.Equal(size, w.Body.Len())
	}()

	var g sync.WaitGroup
	for i := 0; i < 1000; i++ {
		g.Add(1)
		go func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/brotli", nil)
			req.Header.Add("Accept-Encoding", "br")
			req.Host = host
			r.ServeHTTP(w, req)
			tt.Equal(200, w.Code)
			t.Log(w.Body.Len(), size)
			tt.EqualTrue(w.Body.Len() > 100)
			tt.EqualTrue(w.Body.Len() < size)
			g.Done()
		}()
	}
	g.Wait()
}

func BenchmarkGzipDoNotUse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gzip", nil)
		req.Header.Add("Accept-Encoding1", "not-gzip")
		r.ServeHTTP(w, req)
		if 200 != w.Code || size != w.Body.Len() {
			b.Fail()
		}
	}
}

func BenchmarkGzipUse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gzip", nil)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Host = host
		r.ServeHTTP(w, req)
		if 200 != w.Code || size <= w.Body.Len() || 100 >= w.Body.Len() {
			b.Fail()
		}
	}
}
