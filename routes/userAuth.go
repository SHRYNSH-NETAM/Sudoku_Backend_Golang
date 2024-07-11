package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func UserAuthRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/signin", Signin)
		r.Post("/signup", Signup)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuth)
		r.Delete("/deleteAccount", DeleteAccount)
	})

	return r
}

func DeleteAccount(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(""))
}

func Signin(w http.ResponseWriter, r *http.Request){

}

func Signup(w http.ResponseWriter, r *http.Request){
	
}