package main

import (
	"fmt"
	"github.com/jessehorne/superchat-core/database"
	"github.com/joho/godotenv"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := database.InitDB()
	if err != nil {
		fmt.Println("problem with database:", err)
		os.Exit(1)
	}

	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///database/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := m.Down(); err != nil {
		fmt.Println("couldn't run migrations:", err)
	}
}
