package transformer

import (
	"net/url"

	"github.com/johnhoman/go-mitm/internal/context"
)

type Query interface {
	Transform(c *context.Context, query url.Values)
}

type QueryFunc func(c *context.Context, query url.Values)

func (f QueryFunc) Transform(c *context.Context, query url.Values) {
	f(c, query)
}

type QueryChain []Query

func (ch QueryChain) Transform(c *context.Context, query url.Values) {
	for _, f := range ch {
		f.Transform(c, query)
		if c.IsAborted() {
			return
		}
	}
}
