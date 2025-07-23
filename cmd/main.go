package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"miniapp/handlers"
	"net/http"
	"time"
)

func main() {
	db, err := initDb()
	if err != nil {
		fmt.Printf("datadase init err %v: ", err.Error())
	}
	defer db.Close()
	migrationsPath := "migrations"
	err = RunMigrations(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	r := handlers.Manager(db)
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
func RunMigrations(db *sql.DB, migrationsDir string) error {
	_, err := db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	if err != nil {
		log.Fatalf("❌ Failed to reset schema: %v", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func initDb() (*sql.DB, error) {
	var err error
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "miniapp"
	dbPassword := "miniapp"
	dbName := "postgres"
	connectTimeOut := 5

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		dbHost, dbPort, dbUser, dbPassword, dbName, connectTimeOut)
	//connStr := "postgres://postgres_y6xc_user:1s6VRgCIWEgkMPOo7ztvCmqYxWiUd5AB@dpg-d1vil0re5dus739qkqs0-a:5432/postgres_y6xc"
	for i := 0; i < 3; i++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(10)

		if err = db.Ping(); err == nil {
			return db, nil
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", 3, err)
}
