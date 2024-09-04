package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Laptop@2K"
	dbname   = "sudokudb"
  )

var db *sql.DB
var err error

func Connect2psql() error {
	fmt.Println("Connecting to DB....")
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	for i := 0; i < 3; i++ {
		if db == nil { // Only attempt to connect if db is not already set
			db, err = sql.Open("postgres", psqlInfo)
			if err != nil {
				fmt.Println("Failed to connect to Postgres:", err)
				if i < 2 {
					time.Sleep(2 * time.Second) // Wait 2 seconds before retrying
				} else {
					return err
				}
				continue
			}
		}
	
		err = db.Ping()
		if err != nil {
			fmt.Println("Failed to ping Postgres:", err)
			if i < 2 {
				time.Sleep(2 * time.Second) // Wait 2 seconds before retrying
			} else {
				return err
			}
		} else {
			fmt.Println("Successfully connected and pinged DB")
			break
		}
	}
	
	return nil
}