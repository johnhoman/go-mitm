package mitm

import (
	"github.com/johnhoman/go-mitm/internal/transformer"
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
)

var (
	ResetContentLength = transformer.ResetContentLength
)
