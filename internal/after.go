package internal

import "net/http"

const (
	AfterResponseFuncKey = "ContextKey-AfterResponseFunc"
)

type AfterResponseFunc func(req *http.Response) error
