package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/golang-jwt/jwt/v5"
)

func UpdateMyStatistics(userEmail string, rmode string) error{

	mode := 0
	switch rmode {
		case "medium":
			mode = 1
		case "hard":
			mode = 2
		case "extreme":
			mode = 3
	}

	result := initializers.FindData(models.Fuser{Email: userEmail})
	if result==nil {
		return errors.New("user data not Found")
	}

	newStats := result.Statistics
	newStats[mode]++

	if Success := initializers.UpdateData(models.Fuser{Email: result.Email},models.User{Statistics: newStats}); !Success {
		return errors.New("something went wrong")
	}
	return nil
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


