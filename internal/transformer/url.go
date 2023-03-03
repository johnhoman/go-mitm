package transformer

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

type URL interface {
	Transform(c *gin.Context, u *url.URL)
}

type URLFunc func(c *gin.Context, u *url.URL)

func (f URLFunc) Transform(c *gin.Context, u *url.URL) {
	f(c, u)
}

type URLChain []URL

func (ch URLChain) Transform(c *gin.Context, u *url.URL) {
	for _, f := range ch {
		f.Transform(c, u)
		if c.IsAborted() {
			return
		}
	}
}
