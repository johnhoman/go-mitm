package transformer

import (
	"github.com/johnhoman/go-mitm/internal/context"
	"net/http"
	"strconv"
)

type Header interface {
	Transform(c *context.Context, header http.Header)
}

type HeaderFunc func(c *context.Context, header http.Header)

func (f HeaderFunc) Transform(c *context.Context, header http.Header) {
	f(c, header)
}

type HeaderChain []Header

func (ch HeaderChain) Transform(c *context.Context, header http.Header) {
	for _, f := range ch {
		f.Transform(c, header)
		if c.IsAborted() {
			return
		}
	}
}

func RemoveHeader(key string) Header {
	return HeaderFunc(func(c *context.Context, header http.Header) {
		header.Del(key)
	})
}

func SetContentLength(n int) Header {
	return HeaderFunc(func(c *context.Context, header http.Header) {
		header.Set("Content-Length", strconv.Itoa(n))
	})
}

func ResetContentLength(n int) Header {
	return HeaderChain{RemoveHeader("Content-Length"), SetContentLength(n)}
}
