package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func StatisticsRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JwtAuth)

	r.Get("/myStatistics", GetMyStatistics)
	r.Patch("/myStatistics", UpdateMyStatistics)
	
	return r
}

func GetMyStatistics(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(""))
}

func UpdateMyStatistics(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(""))
}