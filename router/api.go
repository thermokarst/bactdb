package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()
	m.Path("/users").Methods("GET").Name(Users)
	m.Path("/users").Methods("POST").Name(CreateUser)
	m.Path("/users/{Id:.+}").Methods("GET").Name(User)
	return m
}
