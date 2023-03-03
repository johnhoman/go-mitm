package mitm

import (
	"encoding/hex"
	"github.com/johnhoman/go-mitm/internal/context"
	"github.com/johnhoman/go-mitm/internal/transformer"
	"regexp"
	"strings"
)

func RegexReplacer(pattern, with string) transformer.String {
	re := regexp.MustCompile(pattern)
	return transformer.StringFunc(func(c *context.Context, s string) string {
		return re.ReplaceAllString(s, with)
	})
}

func HexEncoder() transformer.String {
	return transformer.StringFunc(func(c *context.Context, s string) string {
		return hex.EncodeToString([]byte(s))
	})
}

func TrimEmailDomain() transformer.String {
	return transformer.StringFunc(func(c *context.Context, s string) string {
		s, _, _ = strings.Cut(s, "@")
		return s
	})
}
