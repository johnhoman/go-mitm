package mitm

import (
	"github.com/gin-gonic/gin"
	"github.com/johnhoman/go-mitm/internal/context"
	"github.com/johnhoman/go-mitm/internal/handler"
	"github.com/johnhoman/go-mitm/internal/transformer"
	"net/http"
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

	Context      = context.Context
	HandlerFunc  = handler.Func
	HandlerChain []handler.Func
)

type Engine struct {
	*gin.Engine
}

func (e *Engine) Handle(method string, relativePath string, handlers ...HandlerFunc) {
	chain := gin.HandlersChain{}
	for _, f := range handlers {
		chain = append(chain, handler.WrapF(f))
	}
	e.Engine.Handle(method, relativePath, chain...)
}

func (e *Engine) POST(relativePath string, handlers ...HandlerFunc) {
	e.Handle(http.MethodPost, relativePath, handlers...)
}

func (e *Engine) GET(relativePath string, handlers ...HandlerFunc) {
	e.Handle(http.MethodGet, relativePath, handlers...)
}

func (e *Engine) DELETE(relativePath string, handlers ...HandlerFunc) {
	e.Handle(http.MethodDelete, relativePath, handlers...)
}

func (e *Engine) PATCH(relativePath string, handlers ...HandlerFunc) {
	e.Handle(http.MethodPatch, relativePath, handlers...)
}

var (
	ResetContentLength = transformer.ResetContentLength
)

func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func New() *Engine {
	return &Engine{Engine: gin.New()}
}

func Default() *Engine {
	return &Engine{Engine: gin.Default()}
}
