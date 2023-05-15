package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewDBSqlite() *sql.DB {
	db, err := sql.Open("sqlite3", "./forumdb.sqlite")
	if err != nil {
		log.Fatalln("NewDBSqlite:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("NewDBSqlitePing:", err)
	}

	table, err := os.ReadFile("repository/sqlite/migrations/init_tables.sql")
	if err != nil {
		log.Fatal("NewDBSqliteTable:", err)
	}
	_, err = db.Exec(string(table))
	if err != nil {
		log.Fatal("NewDBSqliteExec:", err)
	}
	return db
}
