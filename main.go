package main

import (
	"fmt"
	"net/http"
	"os"
	// "embed"
	// "io/fs"

	"github.com/MadAppGang/httplog"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/routes"
	"github.com/go-chi/chi/v5"
)

// //go:embed all:build
// var embedBuildFiles embed.FS

func init() {
	initializers.Initenv()
	initializers.Connect2DB()
}

func main() {
	r := chi.NewRouter()

	// staticFiles, _ := fs.Sub(embedBuildFiles, "build")
	// fileServer := http.FileServer(http.FS(staticFiles))

	r.Use(middleware.Cors)
	r.Use(httplog.Logger)

	r.Mount("/api/v1/game", routes.GameRouter())
	r.Mount("/api/v1/auth", routes.UserAuthRouter())
	r.Mount("/api/v1/statistics", routes.StatisticsRouter())

	// r.Handle("/*", http.StripPrefix("/", fileServer))

	PORT := os.Getenv("PORT")
	if PORT=="" {
		PORT = ":8000"
	}
	fmt.Printf("Server Running on port %v \n",PORT)
	http.ListenAndServe(PORT, r)
}