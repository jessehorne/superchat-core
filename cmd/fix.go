package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jessehorne/superchat-core/database"
	"github.com/joho/godotenv"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) <= 1 {
		panic("Not enough params. Try go run cmd/fix.go <NUMBER>")
	}

	// validate number
	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := m.Force(num); err != nil {
		panic(err)
	}
}
