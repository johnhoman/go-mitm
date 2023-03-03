package transformer

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

type Query interface {
	Transform(c *gin.Context, query url.Values)
}

type QueryFunc func(c *gin.Context, query url.Values)

func (f QueryFunc) Transform(c *gin.Context, query url.Values) {
	f(c, query)
}

type QueryChain []Query

func (ch QueryChain) Transform(c *gin.Context, query url.Values) {
	for _, f := range ch {
		f.Transform(c, query)
		if c.IsAborted() {
			return
		}
	}
}
