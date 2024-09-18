package models

import (
	"time"
)

type User struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Statistics [4]int `json:"statistics,omitempty" bson:"statistics"`
	CurrentSudoku Sudokugrid `json:"currentSudoku" bson:"currentSudoku"`
}

type Sudokugrid struct {
	SolvedGrid [][]int `json:"solvedgrid" bson:"solvedgrid"`
	UnSolvedGrid [][]int `json:"unsolvedgrid" bson:"unsolvedgrid"`
	Time time.Time
} 

type Fuser struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type LoginUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Repeatpassword string `json:"repeatpassword"`
	Usernameoremail string `json:"usernameoremail"`
	Token string `json:"token"`
	Mode string `json:"mode"`
}

type ResStruct struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Token string `json:"token"`
	Message string `json:"message"`
	MyStatistics [4]int `json:"myStatistics,omitempty"`
	Result []float64 `json:"result"`
}

type Key string

type SolutionRedis struct {
	ID string `json:"id"`
	Result []float64 `json:"result"`
}