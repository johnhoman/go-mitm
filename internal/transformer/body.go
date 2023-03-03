package transformer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	// TODO: remove this -- it has too many dependencies
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
)

type Body interface {
	Transform(c *gin.Context, body any)
}

type BodyFunc func(c *gin.Context, body any)

func (f BodyFunc) Transform(c *gin.Context, body any) {
	f(c, body)
}

type BodyChain []Body

func (ch BodyChain) Transform(c *gin.Context, body any) {
	for _, f := range ch {
		f.Transform(c, body)
		if c.IsAborted() {
			return
		}
	}
}

// BeforeRequest(
//   Mapper("notebooks",
//     TransformStringField("name", TrimPrefix(), TrimLeft("-")))
// )
//

func TransformStringField(path string, f ...String) Body {
	segments, err := fieldpath.Parse(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse path"))
	}
	return BodyFunc(func(c *gin.Context, body any) {
		current := reflect.ValueOf(body)
		for _, segment := range segments {
			switch segment.Type {
			case fieldpath.SegmentField:
				current = current.FieldByName(segment.Field)
			case fieldpath.SegmentIndex:
				current = current.Index(int(segment.Index))
			}
		}
		chain := StringChain(f)
		v := chain.Transform(c, current.Interface().(string))
		current.Set(reflect.ValueOf(v))
	})
}

func Mapper(path string, f Body) Body {
	segments, err := fieldpath.Parse(path)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse path"))
	}
	return BodyFunc(func(c *gin.Context, body any) {
		v := reflect.ValueOf(body)
		current := v
		for _, sg := range segments {
			switch sg.Type {
			case fieldpath.SegmentField:
				current = current.FieldByName(sg.Field)
			case fieldpath.SegmentIndex:
				if current.Kind() != reflect.Slice {
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				current = current.Index(int(sg.Index))
			}
		}
		if current.Kind() != reflect.Slice {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for k := 0; k < current.Len(); k++ {
			f.Transform(c, current.Index(k).Interface())
		}
	})
}
