package main

import (
	"fmt"
	"net/http"
	"os"

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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Cors)
	r.Use(httplog.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../build/index.html")
	})

	// Serve static files from the "build" directory
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("../build/static"))))

	r.Mount("/api/v1/game", routes.GameRouter())
	r.Mount("/api/v1/auth", routes.UserAuthRouter())
	r.Mount("/api/v1/statistics", routes.StatisticsRouter())

	PORT := os.Getenv("PORT")
	if PORT=="" {
		PORT = ":8000"
	}
	fmt.Printf("Server Running on port %v \n",PORT)
	http.ListenAndServe(PORT, r)
}