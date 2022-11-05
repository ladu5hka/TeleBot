package Database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

var db *sql.DB

func StartDatabase() {
	var err error
	db, err = sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	createTables()
}

func createTables() {
	_, err := db.Exec(
		`create table if not exists "user"
					(
						id      varchar not null
							constraint user_pk
								primary key,
						"group" varchar not null
					);`)

	if err != nil {
		log.Fatal(err)
	}
}
