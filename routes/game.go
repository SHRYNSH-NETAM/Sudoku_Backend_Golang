package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/controllers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func GameRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/",controller.GetSudokuGrid)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuth)
		r.Patch("/validate",controller.ValidateSudoku)
	})

	return r
}