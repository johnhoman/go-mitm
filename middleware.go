package mitm

import (
	"github.com/johnhoman/go-mitm/internal/context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/johnhoman/go-mitm/internal"
	"github.com/johnhoman/go-mitm/internal/transformer"
)

const (
	ErrMiddlewareRequireUsernameHeader = "Missing required prerequisite middleware RequiredUsernameHeader"
)

func RequireHeader(key string) HandlerFunc {
	return func(c *Context) {
		if c.GetHeader(key) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}

func RequireUsernameHeader(key string) HandlerFunc {
	return func(c *Context) {
		if c.GetHeader(key) == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(context.ContextKeyUsername, c.GetHeader(key))
	}
}

func TransformUsername(f ...transformer.String) HandlerFunc {
	return func(c *Context) {
		s := c.GetString(context.ContextKeyUsername)
		chain := transformer.StringChain(f)
		c.Set(context.ContextKeyUsernameTransformed, chain.Transform(c, s))
	}
}

func ProxyAfter(upstream *url.URL, transport http.RoundTripper) HandlerFunc {
	if upstream == nil {
		panic("argument upstream cannot be nil")
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	return func(c *Context) {
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

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type Director interface {
	ModifyRequest(req *http.Request)
}

type DirectorFunc func(req *http.Request)

func (f DirectorFunc) ModifyRequest(req *http.Request) {
	f(req)
}

type proxyOptions struct {
	transport http.RoundTripper
	director  Director
}

type ProxyOption func(o *proxyOptions)

func WithTransport(t http.RoundTripper) ProxyOption {
	return func(o *proxyOptions) {
		o.transport = t
	}
}

func WithDirector(d Director) ProxyOption {
	return func(o *proxyOptions) {
		o.director = d
	}
}

func ProxyTo(upstream *url.URL, opts ...ProxyOption) HandlerFunc {
	o := &proxyOptions{}
	for _, f := range opts {
		f(o)
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	if o.director != nil {
		proxy.Director = func(req *http.Request) {
			o.director.ModifyRequest(req)
		}
	}
	if o.transport != nil {
		proxy.Transport = o.transport
	}
	return WrapH(proxy)
}
