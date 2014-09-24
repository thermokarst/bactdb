package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()
	m.Path("/users/{Id:.+}").Methods("GET").Name(User)
	return m
}
