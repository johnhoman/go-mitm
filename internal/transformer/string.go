package transformer

import (
	"github.com/gin-gonic/gin"
)

type String interface {
	Transform(c *gin.Context, s string) string
}

type StringFunc func(c *gin.Context, s string) string

func (f StringFunc) Transform(c *gin.Context, s string) string {
	return f(c, s)
}

type StringChain []String

func (ch StringChain) Transform(c *gin.Context, s string) string {
	for _, f := range ch {
		s = f.Transform(c, s)
	}
	return s
}
