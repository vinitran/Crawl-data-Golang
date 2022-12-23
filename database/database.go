package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
)

func ConnectDatabase() *bun.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	password := os.Getenv("password")
	dsn := "postgres://postgres:@localhost:5432/postgres?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithPassword(password),
	))
	db := bun.NewDB(sqldb, pgdialect.New())
	fmt.Println("connected to database")
	return db
}
