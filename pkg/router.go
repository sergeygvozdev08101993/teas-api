package pkg

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type router struct {
	*httprouter.Router
}

func NewRouter() *router {
	return &router{httprouter.New()}
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

func (r *router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

func wrapHandler(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		handler.ServeHTTP(w, r)
	}
}
