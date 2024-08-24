package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/controllers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func UserAuthRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/signin", controller.Signin)
		r.Post("/signup", controller.Signup)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuth)
		r.Delete("/deleteAccount", controller.DeleteAccount)
	})

	return r
}