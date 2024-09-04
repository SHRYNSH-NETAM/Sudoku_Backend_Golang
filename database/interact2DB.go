package database

import (
	"database/sql"
	"log"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/lib/pq"
)

func FindUser(loginData models.Fuser) *models.User{
	var result models.User

	sqlStatement := `SELECT password FROM users WHERE username=$1 OR email=$2;`
	row := db.QueryRow(sqlStatement, loginData.Username, loginData.Email)

	switch err := row.Scan(&result.Password); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return &result
	default:
		log.Panic("Error occured:",err)
	}
	return nil
}

func AddUser(loginData models.User) bool {
    sqlStatement := `INSERT INTO users (username, email, password, statistics) VALUES ($1, $2, $3, $4)`
    _, err := db.Exec(sqlStatement, loginData.Username, loginData.Email, loginData.Password, pq.Array(loginData.Statistics))
    if err != nil {
        log.Panic("Error occurred:", err)
        return false
    }
    return true
}