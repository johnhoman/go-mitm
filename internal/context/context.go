package context

import "github.com/gin-gonic/gin"

const (
	ContextKeyUsername            = "ContextKeyUsername"
	ContextKeyUsernameTransformed = "ContextKeyUsername-Transformed"
)

type Context struct {
	*gin.Context
}

func (c *Context) Username() string {
	return c.GetString(ContextKeyUsername)
}

func (c *Context) UsernameTransformed() string {
	return c.GetString(ContextKeyUsernameTransformed)
}

func New(c *gin.Context) *Context {
	return &Context{Context: c}
}
