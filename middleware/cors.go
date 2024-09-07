package middleware

import "net/http"

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// w.Header().Set("Access-Control-Allow-Origin", "https://sudoku-frontend-xi.vercel.app")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		// w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		// w.Header().Set("Access-Control-Allow-Credentials","true")
		// w.Header().Set("Access-Control-Max-Age", "300")

		// if r.Method == http.MethodOptions {
		// 	w.WriteHeader(http.StatusOK)
		// 	return
		// }
		
		// next.ServeHTTP(w, r)
		origin := r.Header.Get("Origin")
        if origin == "https://sudoku-frontend-xi.vercel.app" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        } else {
            
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
	})
}