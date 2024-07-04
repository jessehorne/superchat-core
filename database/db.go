package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"database/sql"
	"fmt"
	"os"
)

func GetDSN() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	name := os.Getenv("MYSQL_DB")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		user, pass, host, port, name)
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", GetDSN())

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitGDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
