package models

type User struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Statistics [4]int `json:"statistics,omitempty" bson:"statistics"`
	Grid [][]int `json:"grid" bson:"grid"`
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
}

type Key string