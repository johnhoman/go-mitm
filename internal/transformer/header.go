package transformer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Header interface {
	Transform(c *gin.Context, header http.Header)
}

type HeaderFunc func(c *gin.Context, header http.Header)

func (f HeaderFunc) Transform(c *gin.Context, header http.Header) {
	f(c, header)
}

type HeaderChain []Header

func (ch HeaderChain) Transform(c *gin.Context, header http.Header) {
	for _, f := range ch {
		f.Transform(c, header)
		if c.IsAborted() {
			return
		}
	}
}

func RemoveHeader(key string) Header {
	return HeaderFunc(func(c *gin.Context, header http.Header) {
		header.Del(key)
	})
}

func SetContentLength(n int) Header {
	return HeaderFunc(func(c *gin.Context, header http.Header) {
		header.Set("Content-Length", strconv.Itoa(n))
	})
}

func ResetContentLength(n int) Header {
	return HeaderChain{RemoveHeader("Content-Length"), SetContentLength(n)}
}
