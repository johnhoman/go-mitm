package transformer

import (
	"github.com/johnhoman/go-mitm/internal/context"
)

type String interface {
	Transform(c *context.Context, s string) string
}

type StringFunc func(c *context.Context, s string) string

func (f StringFunc) Transform(c *context.Context, s string) string {
	return f(c, s)
}

type StringChain []String

func (ch StringChain) Transform(c *context.Context, s string) string {
	for _, f := range ch {
		s = f.Transform(c, s)
	}
	return s
}
