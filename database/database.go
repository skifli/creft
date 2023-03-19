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
var DatabaseJSON map[string]map[string]any

func updateDatabase(logger *golog.Logger) {
	databaseBytes, err := json.Marshal(DatabaseJSON)

	if Changed {
		logger.Info("RAM database has changed, updating on-disk database.", nil)
		Changed = false
	}

	if err != nil {
		logger.Errorf("Failed to convert database to bytes for writing to disk: %s", nil, err)
		return
	}

	databaseFile.Truncate(0)
	databaseFile.Seek(0, 0)
	_, err = databaseFile.Write(databaseBytes)

	if err != nil {
		logger.Errorf("Failed to write database to disk: %s", nil, err)
		return
	}

	logger.Info("Updated on-disk database.", nil)
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
			DatabaseJSON = map[string]map[string]any{"counting": {}}
			updateDatabase(logger)
		}

		if err != nil {
			logger.Fatalf("Failed to read database file: %s", nil, err)
		}
	} else {
		databaseTmp, err := io.ReadAll(databaseFile)

		if err != nil {
			logger.Fatalf("Failed to parse database file: %s", nil, err)
		}

		err = json.Unmarshal(databaseTmp, &DatabaseJSON)

		if err != nil {
			logger.Fatalf("Failed to parse database file: %s", nil, err)
		}
	}

	go updateDatabaseLoop(logger, pause)
}
