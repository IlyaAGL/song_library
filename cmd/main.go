package main

import (
	"database/sql"
	"os"

	"github.com/agl/music_library/api/handlers"
	. "github.com/agl/music_library/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", DATABASE_URL)

	if err != nil {
		Log.Info("Could not set database driver", "err", err)

		return
	}
	defer db.Close()

	Log.Info("Successfully set database driver")

	err = db.Ping()
	if err != nil {
		Log.Info("Could not ping the database", "err", err)

		return
	}

	Log.Info("Connection was set")

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		Log.Info("Could not create migration driver")

		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
		"postgres", driver)

	if err != nil {
		Log.Info("Could not create migrate instance", "err", err)

	return
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		Log.Error("Could not apply migrations: %v", err)

		return
	}

	Log.Info("Migrations applied successfully!")
	Log.Info("Successfully connected to db")

	handlers.SetAPI(db)
}
