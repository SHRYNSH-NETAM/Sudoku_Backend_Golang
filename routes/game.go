package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func GameRouter() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/",GetSudokuGrid)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuth)
		r.Patch("/validate",ValidateSudoku)
	})

	return r
}

func GetSudokuGrid(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Get Game"))
}

func ValidateSudoku(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Validate Sudoku"))
}