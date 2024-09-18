package handler

import (
	"net/http"
	"embed"
	"io/fs"

	"github.com/MadAppGang/httplog"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/routes"
	"github.com/go-chi/chi/v5"
)

//go:embed all:build
var embedBuildFiles embed.FS

func init() {
	initializers.Connect2DB()
	initializers.Connect2Redis()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rtr := chi.NewRouter()

	staticFiles, _ := fs.Sub(embedBuildFiles, "build")
	fileServer := http.FileServer(http.FS(staticFiles))

	rtr.Use(middleware.Cors)
	rtr.Use(httplog.Logger)
	
	rtr.Handle("/*", http.StripPrefix("/", fileServer))	
	rtr.Mount("/api/v1/game", routes.GameRouter())
	rtr.Mount("/api/v1/auth", routes.UserAuthRouter())
	rtr.Mount("/api/v1/statistics", routes.StatisticsRouter())

	rtr.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "build/index.html")
	})

	rtr.ServeHTTP(w, r)
}