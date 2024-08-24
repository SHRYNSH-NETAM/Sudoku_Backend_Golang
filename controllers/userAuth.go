package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/initializers"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/middleware"
	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signin(w http.ResponseWriter, r *http.Request) {

	var req models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	cleanUsernameoremail := strings.ToLower(req.Usernameoremail)

	existingUserEmail := make(chan *models.User)
	existingUserName := make(chan *models.User)

	go func() {
		existingUserEmail <- initializers.FindData(models.Fuser{Email: cleanUsernameoremail})
	} ()

	go func() {
		existingUserName <- initializers.FindData(models.Fuser{Username: cleanUsernameoremail})
	} ()

	existingUser := <- existingUserEmail
	if(existingUser==nil) {
		existingUser = <- existingUserName
	}

	// existingUser := initializers.FindData(models.Fuser{Email: cleanUsernameoremail}); 
	// if existingUser==nil {
	// 	existingUser = initializers.FindData(models.Fuser{Username: cleanUsernameoremail});
	// }

	if existingUser==nil {
		http.Error(w,"User Not Found", http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "User Not Found"})
		return
	}

	isPasswordCorrect := bcrypt.CompareHashAndPassword([]byte(existingUser.Password),[]byte(req.Password))
	if isPasswordCorrect != nil {
		http.Error(w,"Password is Incorrect",http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "Password is Incorrect"})
		return 
	}

	tokenString, err := middleware.CreateToken(existingUser.Email,existingUser.Username)
	if err!=nil {
		http.Error(w,"Error while creating Token",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ResStruct{Token: tokenString})
}

func Signup(w http.ResponseWriter, r *http.Request) {

	var req models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err!=nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	cleanEmail := strings.ToLower(req.Email)
	cleanUsername := strings.ToLower(req.Username)

	existingEmail := initializers.FindData(models.Fuser{Email: cleanEmail})
	if existingEmail!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "This email is already associated with an existing account"})
		return
	}

	existingUsername := initializers.FindData(models.Fuser{Username: cleanUsername})
	if existingUsername!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "Username already exists"})
		return
	}

	if req.Password!=req.Repeatpassword {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "Passwords don't match"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password),12)
	if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "Error"})
		return
	}

	newUser := models.User{Username: cleanUsername, Email: cleanEmail, Password: string(hashedPwd)}

	done := initializers.AddData(newUser)

	if !done {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResStruct{Message: "Something went wrong!"})
		return
	}

	tokenString, err := middleware.CreateToken(newUser.Email,newUser.Username)

	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ResStruct{Token: tokenString})
}

func DeleteAccount(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	const jwtPayloadKey models.Key = "jwtPayload"

	claims, ok := r.Context().Value(jwtPayloadKey).(jwt.MapClaims) 
	if !ok {
		http.Error(w, "Could not retrieve JWT payload", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
        http.Error(w, "Email not found in JWT payload", http.StatusUnauthorized)
        return
    }

	if Success := initializers.DeleteData(models.Fuser{Email: email}); !Success {
		http.Error(w, `{"message": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.ResStruct{Message: "User deleted successfully"})
}