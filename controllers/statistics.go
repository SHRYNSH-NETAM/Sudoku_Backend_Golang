package controller

import (
	"encoding/json"
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/golang-jwt/jwt/v5"
)

func UpdateMyStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userData models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&userData); err!=nil {
		http.Error(w,"Error occured while decoding userData",http.StatusInternalServerError)
		return
	}

	mode := 0
	switch userData.Mode {
		case "medium":
			mode = 1
		case "hard":
			mode = 2
		case "extreme":
			mode = 3
	}

	var jwtPayload models.Key = "jwtPayload"
	claims, err := r.Context().Value(jwtPayload).(jwt.MapClaims)
	if !err {
		http.Error(w,"Could not retrieve JWT Payload",http.StatusUnauthorized)
		return
	}

	userEmail, err := claims["email"].(string)
	if !err {
		http.Error(w, "Email not found in JWT payload", http.StatusUnauthorized)
        return
	}

	result := initializers.FindData(models.Fuser{Email: userEmail})
	if result==nil {
		http.Error(w,"User data not Found",http.StatusNotFound)
		return
	}

	newStats := result.Statistics
	newStats[mode]++

	if Success := initializers.UpdateData(models.Fuser{Email: result.Email},models.User{Statistics: newStats}); !Success {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.ResStruct{MyStatistics: result.Statistics})
}

func GetMyStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jwtPayload models.Key = "jwtPayload"

	claims, err := r.Context().Value(jwtPayload).(jwt.MapClaims)
	if !err {
		http.Error(w,"Could not retrieve JWT Payload",http.StatusUnauthorized)
		return
	}

	userEmail, err := claims["email"].(string)
	if !err {
		http.Error(w, "Email not found in JWT payload", http.StatusUnauthorized)
        return
	}

	result := initializers.FindData(models.Fuser{Email: userEmail})
	if result==nil {
		http.Error(w,"User data not Found",http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(models.ResStruct{MyStatistics: result.Statistics})
}