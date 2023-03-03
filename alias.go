package mitm

import (
	"github.com/gin-gonic/gin"
	"github.com/johnhoman/go-mitm/internal"
	"github.com/johnhoman/go-mitm/internal/transformer"
)

type (
	BodyTransformer   = transformer.Body
	QueryTransformer  = transformer.Query
	HeaderTransformer = transformer.Header
	StringTransformer = transformer.String

	BodyTransformerFunc   = transformer.BodyFunc
	QueryTransformerFunc  = transformer.QueryFunc
	HeaderTransformerFunc = transformer.HeaderFunc
	StringTransformerFunc = transformer.StringFunc

	BodyTransformerChain   = transformer.BodyChain
	QueryTransformerChain  = transformer.QueryChain
	HeaderTransformerChain = transformer.HeaderChain
	StringTransformerChain = transformer.StringChain
)

var (
	ResetContentLength = transformer.ResetContentLength
)

func User(c *gin.Context) string {
	return c.GetString(internal.ContextKeyUsername)
}
func Username(c *gin.Context) string { return User(c) }

func UserT(c *gin.Context) string {
	return c.GetString(internal.ContextKeyUsernameTransformed)
}
func UsernameT(c *gin.Context) string { return UserT(c) }
