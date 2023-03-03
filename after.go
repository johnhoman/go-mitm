package mitm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/johnhoman/go-mitm/internal"
	"github.com/johnhoman/go-mitm/internal/transformer"
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
		c.Set(internal.AfterResponseFuncKey, internal.AfterResponseFunc(func(res *http.Response) error {
			if res.Body != nil && len(bodyChain) > 0 {
				m := into()
				if err := json.NewDecoder(res.Body).Decode(m); err != nil {
					_ = res.Body.Close()
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
				res.Body = io.NopCloser(buf)
				res.ContentLength = int64(buf.Len())
				headerChain = append(headerChain, ResetContentLength(buf.Len()))
			}
			headerChain.Transform(c, res.Header)
			return nil
		}))
	}
}
