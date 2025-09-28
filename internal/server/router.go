package server

import "regexp"

type Router struct {
	routes   []route
	notFound Handler
}

type route struct {
	pattern string
	handler Handler
}

// temporarily use regex, we will write better routing in a while
func (router *Router) route(url string) Handler {

	for _, route := range router.routes {
		if match, err := regexp.MatchString(route.pattern, url); match && err == nil {
			return route.handler
		}
	}
	return router.notFound
}

func (router *Router) Register(url string, handler Handler) {
	router.routes = append(router.routes, route{url, handler})
}

func (router *Router) RegisterNotFound(handler Handler) {
	router.notFound = handler
}
