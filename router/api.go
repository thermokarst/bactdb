package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()

	// Users
	m.Path("/users").Methods("GET").Name(Users)
	m.Path("/users").Methods("POST").Name(CreateUser)
	m.Path("/users/{Id:.+}").Methods("GET").Name(User)

	// Genera
	m.Path("/genera").Methods("GET").Name(Genera)
	m.Path("/genera").Methods("POST").Name(CreateGenus)
	m.Path("/genera/{Id:.+}").Methods("GET").Name(Genus)
	return m
}
