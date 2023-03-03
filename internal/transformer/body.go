package transformer

import (
	"github.com/johnhoman/go-mitm/internal/context"
)

type Body interface {
	Transform(c *context.Context, body any)
}

type BodyFunc func(c *context.Context, body any)

func (f BodyFunc) Transform(c *context.Context, body any) {
	f(c, body)
}

type BodyChain []Body

func (ch BodyChain) Transform(c *context.Context, body any) {
	for _, f := range ch {
		f.Transform(c, body)
		if c.IsAborted() {
			return
		}
	}
}
