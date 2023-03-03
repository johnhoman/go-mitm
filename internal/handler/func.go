package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/johnhoman/go-mitm/internal/context"
)

type Func func(c *context.Context)

func WrapF(f Func) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(context.New(c))
	}
}

func WrapChain(fns ...Func) gin.HandlersChain {
	chain := gin.HandlersChain{}
	for _, f := range fns {
		chain = append(chain, WrapF(f))
	}
	return chain
}
