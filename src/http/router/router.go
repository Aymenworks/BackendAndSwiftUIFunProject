package router

import (
	"net/http"
)

type Router interface {
	http.Handler
	Get(uri string, f http.HandlerFunc)
	Post(uri string, f http.HandlerFunc)
	Put(uri string, f http.HandlerFunc)
	Delete(uri string, f http.HandlerFunc)
	UseMiddleware(f func(http.Handler) http.Handler)
}
