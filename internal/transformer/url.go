package transformer

import (
	"github.com/johnhoman/go-mitm/internal/context"
	"net/url"
)

type URL interface {
	Transform(c *context.Context, u *url.URL)
}

type URLFunc func(c *context.Context, u *url.URL)

func (f URLFunc) Transform(c *context.Context, u *url.URL) {
	f(c, u)
}

type URLChain []URL

func (ch URLChain) Transform(c *context.Context, u *url.URL) {
	for _, f := range ch {
		f.Transform(c, u)
		if c.IsAborted() {
			return
		}
	}
}
