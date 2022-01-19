package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	migrations := &migrate.FileMigrationSource{Dir: "cmd/db_migrate"}

	c := mysql.Config{
		User:                 "root",
		Passwd:               "mypasswd",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "bsc",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	_, err = migrate.Exec(db, "mysql", migrations, migrate.Down)
	if err != nil {
		log.Fatal(err)
	}
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
