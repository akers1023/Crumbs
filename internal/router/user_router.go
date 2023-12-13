package router

import (
	handler "crumbs/internal/handler/user"
	"crumbs/internal/middleware"

	"github.com/gorilla/mux"
)

func UserRouter(r *mux.Router) *mux.Router {
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login/email", middleware.Authenticate(handler.Login)).Methods("POST")
	r.HandleFunc("/login/phone", middleware.Authenticate(handler.Login)).Methods("POST")
	r.HandleFunc("/login/user_name", middleware.Authenticate(handler.Login)).Methods("POST")
	r.HandleFunc("/users", middleware.Authenticate(handler.GetUsers)).Methods("GET")
	r.HandleFunc("/users/{user_id}", middleware.Authenticate(handler.GetUser)).Methods("GET")
	// r.HandleFunc("/users", ).Methods("GET")
	return r
}
