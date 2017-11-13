package tuu

import (
	"fmt"
	"net/http"
	"strings"
)

func NewRouter() *DefaultRouter {
	return &DefaultRouter{}
}

type DefaultRouter struct {
	Routes       []*Route
	StaticRoutes []*StaticRoute

	prefix string
}

func (r *Prefix) Prefix(path string) {
	r.prefix = path
}

func (r *DefaultRouter) GET(path string, h Handler) {
	r.addRoute(http.MethodGet, path, h)
}

func (r *DefaultRouter) POST(path string, h Handler) {
	r.addRoute(http.MethodPost, path, h)
}

func (r *DefaultRouter) Static(path string, root http.FileSystem) {
	r.StaticRoutes = append(r.StaticRoutes, &StaticRoute{
		Path:    path,
		Handler: http.StripPrefix(path, http.FileServer(root)),
	})
}

func (r *DefaultRouter) NotFound(path string, h Handler) {

}

func (r *DefaultRouter) GetRoutes() []*Route {
	return r.Routes
}

func (r *DefaultRouter) GetStaticRoutes() []*StaticRoute {
	return r.StaticRoutes
}

func (r *DefaultRouter) addRoute(m, p string, h Handler) {
	path := fmt.Sprintf("/%s/%s", strings.TrimPrefix(r.prefix, "/"), strings.TrimSuffix(p, "/"))

	r.Routes = append(r.Routes, &Route{
		Method:  m,
		Path:    path,
		Handler: h,
	})
}
