package mitm

import (
	"github.com/gin-gonic/gin"
	"github.com/johnhoman/mitm/internal"
	"github.com/johnhoman/mitm/internal/transformer"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	ContextKeyUsername            = "ContextKeyUsername"
	ContextKeyUsernameTransformed = "ContextKeyUsername-Transformed"

	ErrMiddlewareRequireUsernameHeader = "Missing required prerequisite middleware RequiredUsernameHeader"
)

func RequireHeader(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(key) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}

func RequireUsernameHeader(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(key) == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ContextKeyUsername, c.GetHeader(key))
	}
}

func TransformUsername(f ...transformer.String) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.GetString(ContextKeyUsername)
		chain := transformer.StringChain(f)
		c.Set(ContextKeyUsernameTransformed, chain.Transform(c, s))
	}
}

func ProxyAfter(upstream *url.URL, transport http.RoundTripper) gin.HandlerFunc {
	if upstream == nil {
		panic("argument upstream cannot be nil")
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	return func(c *gin.Context) {
		c.Next()
		if c.IsAborted() || c.Writer.Written() {
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(upstream)
		proxy.Transport = transport
		f, ok := c.Get(internal.AfterResponseFuncKey)
		if ok {
			proxy.ModifyResponse = f.(internal.AfterResponseFunc)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
