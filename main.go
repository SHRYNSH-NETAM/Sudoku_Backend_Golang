package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Mount("/api/v1/game", routes.GameRouter())
	r.Mount("/api/v1/auth", routes.UserAuthRouter())
	r.Mount("/api/v1/statistics", routes.StatisticsRouter())

	PORT := os.Getenv("PORT")
	fmt.Printf("Server Running on port %v \n",PORT)
	http.ListenAndServe(PORT, r)
}