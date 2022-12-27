package expenses

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDb() {

	var err error

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	createTb := `
    CREATE TABLE IF NOT EXISTS expenses (
        id SERIAL PRIMARY KEY,
        title TEXT,
        amount FLOAT,
        note TEXT,
        tags TEXT[]
    );`

	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

}
