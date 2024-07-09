package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GameRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/",GetSudokuGrid)

	return r
}

func GetSudokuGrid(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Get Game"))
}