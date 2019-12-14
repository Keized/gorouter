package gorouter

import (
	"net/http"
	"path"
	"regexp"
)


/***************
**** PARAMS ****
****************/
type Param struct {
	key 	string
	value 	string
}

type Params struct {
	items []Param
}

func (params *Params) Get(key string) (s string){
	for _, p := range params.items {
		if p.key == key {
			return p.value
		}
	}
	return ""
}

func (params *Params) Add(param Param) (err error) {
	params.items = append(params.items, param)
	return
}

/***************
**** ROUTER ****
****************/
type Handler func(http.ResponseWriter, *http.Request, Params)

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
			var params Params
			paramsValues := reg2.FindStringSubmatch(url)
			paramsIndexes := reg1.FindAllString(rt.path, -1)
			paramsValues = paramsValues[2:]

			for i, param := range paramsIndexes {
				err := params.Add(Param{param[1:], paramsValues[i]})
				if err != nil {
					panic(err)
				}
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
	_, err := w.Write([]byte("404: Page not found"))
	if err != nil {
		panic(err)
	}
}