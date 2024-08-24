package routes

import (
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/controllers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/go-chi/chi/v5"
)

func StatisticsRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.JwtAuth)

	r.Get("/myStatistics", controller.GetMyStatistics)
	
	return r
}