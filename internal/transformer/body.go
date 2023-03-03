package transformer

import (
	"github.com/gin-gonic/gin"
)

type Body interface {
	Transform(c *gin.Context, body any)
}

type BodyFunc func(c *gin.Context, body any)

func (f BodyFunc) Transform(c *gin.Context, body any) {
	f(c, body)
}

type BodyChain []Body

func (ch BodyChain) Transform(c *gin.Context, body any) {
	for _, f := range ch {
		f.Transform(c, body)
		if c.IsAborted() {
			return
		}
	}
}
