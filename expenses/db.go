package expenses

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDb() *handler {

	var err error

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	h := NewApplication(db)

	createTb := `
    CREATE TABLE IF NOT EXISTS expenses (
        id SERIAL PRIMARY KEY,
        title TEXT,
        amount FLOAT,
        note TEXT,
        tags TEXT[]
    );`

	_, err = h.DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
	return h
}
