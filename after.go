package mitm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/johnhoman/mitm/internal"
	"github.com/johnhoman/mitm/internal/transformer"
)

const (
	ErrAbortAfterResponse = "the request was aborted after the response from the upstream server"
)

func AfterResponse(into func() any, opts ...any) gin.HandlerFunc {
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
			panic("invalid transformer passed to AfterRequest")
		}
	}
	v := into()
	if v == nil || reflect.ValueOf(v).Kind() != reflect.Pointer {
		panic("argument into() must return a non nil pointer")
	}
	return func(c *gin.Context) {
		c.Set(internal.AfterResponseFuncKey, internal.AfterResponseFunc(func(req *http.Response) error {
			if req.Body != nil && len(bodyChain) > 0 {
				m := into()
				if err := json.NewDecoder(req.Body).Decode(m); err != nil {
					_ = req.Body.Close()
					return err
				}
				bodyChain.Transform(c, m)
				if c.IsAborted() {
					return errors.New(ErrAbortAfterResponse)
				}
				buf := new(bytes.Buffer)
				if err := json.NewEncoder(buf).Encode(m); err != nil {
					return err
				}
				c.Request.Body = io.NopCloser(buf)
				c.Request.ContentLength = int64(buf.Len())
				headerChain = append(headerChain, ResetContentLength(buf.Len()))
			}
			q := c.Request.URL.Query()
			queryChain.Transform(c, q)
			c.Request.URL.RawQuery = q.Encode()
			headerChain.Transform(c, c.Request.Header)
			return nil
		}))
	}
}
