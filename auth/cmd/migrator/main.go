package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var storagePath, migrationPath string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage file")
	flag.StringVar(&migrationPath, "migration-path", "", "path to migration files")

	flag.Parse()

	if storagePath == "" {
		panic("storage-path is empty")
	}
	if migrationPath == "" {
		panic("migration-path is empty")
	}

	m, err := migrate.New("file://"+migrationPath, "sqlite3://"+storagePath)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		} else {
			panic(err)
		}
	}
	fmt.Println("applied migrations")
}
