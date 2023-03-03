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
