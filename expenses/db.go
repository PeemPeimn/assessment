package expenses

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(url string) *sql.DB {

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Cannot connect to the database.", err)
	}

	createStatement := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createStatement)

	if err != nil {
		log.Fatal("Cannot create the table.", err)
	}

	return db
}
