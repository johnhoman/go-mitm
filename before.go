package mitm

import (
	"bytes"
	"encoding/json"
	"github.com/johnhoman/go-mitm/internal/handler"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin/binding"

	"github.com/johnhoman/go-mitm/internal/transformer"
)

func BeforeRequest(into func() any, opts ...any) handler.Func {
	var (
		bodyChain   transformer.BodyChain
		headerChain transformer.HeaderChain
		queryChain  transformer.QueryChain
	)
	for _, f := range opts {
		switch fn := f.(type) {
		case transformer.Query:
			queryChain = append(queryChain, fn)
		case transformer.Body:
			bodyChain = append(bodyChain, fn)
		case transformer.HeaderChain:
			headerChain = append(headerChain, fn)
		default:
			panic("invalid transformer passed to BeforeRequest")
		}
	}
	v := into()
	if v == nil || reflect.ValueOf(v).Kind() != reflect.Pointer {
		panic("argument into() must return a non nil pointer")
	}
	return func(c *Context) {
		if c.Request.Body != nil && len(bodyChain) > 0 {
			m := into()
			if err := c.MustBindWith(m, binding.JSON); err != nil {
				return
			}
			bodyChain.Transform(c, m)
			if c.IsAborted() {
				return
			}
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(m); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Request.Body = io.NopCloser(buf)
			c.Request.ContentLength = int64(buf.Len())
			headerChain = append(headerChain, ResetContentLength(buf.Len()))
		}
		q := c.Request.URL.Query()
		queryChain.Transform(c, q)
		c.Request.URL.RawQuery = q.Encode()
		headerChain.Transform(c, c.Request.Header)
	}
}
