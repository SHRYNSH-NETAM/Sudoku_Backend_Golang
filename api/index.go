package handler

import (
	"net/http"

	"github.com/MadAppGang/httplog"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/routes"
	"github.com/go-chi/chi/v5"
)

func init() {
	initializers.Initenv()
	initializers.Connect2DB()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rtr := chi.NewRouter()

	rtr.Use(middleware.Cors)
	rtr.Use(httplog.Logger)
	rtr.Get("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})	
	rtr.Mount("/api/v1/game", routes.GameRouter())
	rtr.Mount("/api/v1/auth", routes.UserAuthRouter())
	rtr.Mount("/api/v1/statistics", routes.StatisticsRouter())

	rtr.ServeHTTP(w, r)
}