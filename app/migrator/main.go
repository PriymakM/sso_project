package main

import (
	"errors"
	"flag"
	"fmt"

	// "github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationPath, "migrations-path", "", "path to migration")
	flag.StringVar(&migrationTable, "migrations-table", "", "path to migration-table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path empty!")
	}
	if migrationPath == "" {
		panic("migration path empty!")
	}

	m, err := migrate.New("file://"+migrationPath, fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationTable))
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply!")
			return
		}
		panic(err)
	}

	fmt.Println("Migration successfully!")
}
