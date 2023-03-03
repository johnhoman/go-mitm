package transformer

import (
	"github.com/gin-gonic/gin"
	"strings"
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

type StringGetter func(c *gin.Context) string

func TrimPrefix(getter StringGetter) String {
	return StringFunc(func(c *gin.Context, s string) string {
		return strings.TrimPrefix(s, getter(c))
	})
}

func TrimLeft(v string) String {
	return StringFunc(func(c *gin.Context, s string) string {
		for strings.HasPrefix(s, v) {
			s = strings.TrimPrefix(s, v)
		}
		return s
	})
}

func TrimRight(v string) String {
	return StringFunc(func(c *gin.Context, s string) string {
		for strings.HasSuffix(s, v) {
			s = strings.TrimSuffix(s, v)
		}
		return s
	})
}
