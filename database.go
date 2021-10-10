package main

import (
	"database/sql"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func connectDB() error {
	var err error
	log.Info("connecting to database...")
	db, err = sql.Open("sqlite3", "file:counter.sqlite?cache=shared")
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return err
	}
	db.SetMaxOpenConns(1)
	return nil
}

func createDB() error {
	log.Info("creating database...")
	os.Create("counter.sqlite")

	db, err := sql.Open("sqlite3", "file:counter.sqlite?cache=shared")
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return err
	}
	db.SetMaxOpenConns(1)

	log.Info("creating counters table...")

	createStatement := `CREATE TABLE counters (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"path" TEXT,
		"count" integer
	);`

	statement, err := db.Prepare(createStatement)
	if err != nil {
		err = fmt.Errorf("failed to create counters table: %w", err)
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		err = fmt.Errorf("failed to create counters table: %w", err)
		return err
	}

	return nil
}

func addPath(path string) error {
	log.Debug("adding new path", path)
	insertSQL := `INSERT INTO counters (path, count) VALUES (?, 0)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		err = fmt.Errorf("failed to prepare insert statement: %w", err)
		return err
	}
	_, err = statement.Exec(path)
	if err != nil {
		err = fmt.Errorf("failed to execute insert statement: %w", err)
		return err
	}

	return nil
}

func getCount(path string) (int, error) {
	var count int
	log.Debug("fetching count for path", path)
	query := fmt.Sprintf(`SELECT count FROM counters WHERE path = '%s' LIMIT 1`, path)
	if err := db.QueryRow(query).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err := addPath(path)
			if err != nil {
				err = fmt.Errorf("failed to add new count: %w", err)
				return 0, err
			}
		}
		err = fmt.Errorf("failed to fetch count: %w", err)
		return 0, err
	}

	return count, nil
}

func incrementCount(path string) (int, error) {
	log.Debug("incrementing count...")
	currentCount, err := getCount(path)
	if err != nil {
		err = fmt.Errorf("failed to get count to increment: %w", err)
		return 0, err
	}
	currentCount++

	updateSQL := fmt.Sprintf(`UPDATE counters SET count = %d WHERE path = '%s'`, currentCount, path)
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		err = fmt.Errorf("failed to prepare increment count statement: %w", err)
		return 0, err
	}
	_, err = statement.Exec()
	if err != nil {
		err = fmt.Errorf("failed to execute increment count statement: %w", err)
		return 0, err
	}

	return currentCount, nil
}
