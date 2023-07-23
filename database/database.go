package database

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/skifli/golog"
)

var Changed = false
var databaseFile *os.File
var DatabaseJSON map[string]any

func updateDatabase(logger *golog.Logger) {
	databaseBytes, err := json.Marshal(DatabaseJSON)

	if Changed {
		logger.Info("RAM database has changed, updating on-disk database.")
		Changed = false
	}

	if err != nil {
		logger.Errorf("Failed to convert database to bytes for writing to disk: %s", err)
		return
	}

	databaseFile.Truncate(0)
	databaseFile.Seek(0, 0)
	_, err = databaseFile.Write(databaseBytes)

	if err != nil {
		logger.Errorf("Failed to write database to disk: %s", err)
		return
	}

	logger.Info("Updated on-disk database.")
}

func updateDatabaseLoop(logger *golog.Logger, pause time.Duration) {
	for {
		time.Sleep(pause)

		if Changed { // If the database has changed...
			updateDatabase(logger)
		}
	}
}

func Init(logger *golog.Logger, pause time.Duration) {
	var err error

	databaseFile, err = os.OpenFile("./database.json", os.O_RDWR, 0777)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			databaseFile, err = os.OpenFile("./database.json", os.O_RDWR|os.O_CREATE, 0777)
			DatabaseJSON = map[string]any{"counting": map[string]any{}, "admins": map[string]any{}}
			updateDatabase(logger)
		}

		if err != nil {
			logger.Fatalf("Failed to read database file: %s", err)
		}
	} else {
		databaseTmp, err := io.ReadAll(databaseFile)

		if err != nil {
			logger.Fatalf("Failed to parse database file: %s", err)
		}

		err = json.Unmarshal(databaseTmp, &DatabaseJSON)

		if err != nil {
			logger.Fatalf("Failed to parse database file: %s", err)
		}
	}

	go updateDatabaseLoop(logger, pause)
}
