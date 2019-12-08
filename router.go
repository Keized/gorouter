package routergo

import (
	"net/http"
	"path"
	"regexp"
)

type Handler func(http.ResponseWriter, *http.Request, map[string]string)

type Route struct {
	method  string
	path 	string
	handler Handler
}

type Router struct{
	Routes []Route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := path.Clean(req.URL.Path)

	for _, rt := range r.Routes {
		reg1 := regexp.MustCompile(`(:[\w]+)`)
		p := reg1.ReplaceAllString(rt.path, `([^/]+)`)
		reg2 := regexp.MustCompile(`(^` + p + `)`)
		match := reg2.FindString(url)

		if match == url {
			params := map[string]string {}
			paramsValues := reg2.FindStringSubmatch(url)
			paramsIndexes := reg1.FindAllString(rt.path, -1)
			paramsValues = paramsValues[2:]

			for i, param := range paramsIndexes {
				params[param[1:]] = paramsValues[i]
			}

			rt.handler(w, req, params)
			return
		}
	}

	notFound(w, req)
	return
}

func (r *Router) Register(method, path string, handler Handler) {
	route := Route{method: method, path: path, handler: handler}
	r.Routes = append(r.Routes, route)
}

func (r *Router) GET(path string, handler Handler) {
	r.Register(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.Register(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.Register(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler Handler) {
	r.Register(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.Register(http.MethodDelete, path, handler)
}

func (r *Router) OPTIONS(path string, handler Handler) {
	r.Register(http.MethodOptions, path, handler)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404: Page not found"))
}