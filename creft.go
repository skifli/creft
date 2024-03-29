package main

import (
	"fmt"
	"os"
	"time"

	"github.com/skifli/creft/commands"
	"github.com/skifli/creft/database"

	"github.com/alexflint/go-arg"
	"github.com/goccy/go-json"
	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

var args struct {
	Config string        `help:"Config file path." arg:"positional" default:"./config.json"`
	Pause  time.Duration `help:"Time between database writes to disk (in seconds)." default:"5s"`
}

var logger = golog.NewLogger([]*golog.Log{
	golog.NewLog(
		[]*os.File{
			os.Stderr,
		},
		golog.FormatterHuman,
	),
})

func main() {
	parser := arg.MustParse(&args)

	configBytes, err := os.ReadFile(args.Config)

	if err != nil {
		parser.Fail(fmt.Sprintf("failed to read config file: %s", err))
	}

	var config map[string]string

	if err = json.Unmarshal(configBytes, &config); err != nil {
		parser.Fail(fmt.Sprintf("failed to parse config file: %s", err))
	}

	bot := &disgo.Client{
		ApplicationID:  config["appID"],
		Authentication: disgo.BotToken(config["botToken"]),
		Config:         disgo.DefaultConfig(),
		Handlers:       new(disgo.Handlers),
		Sessions:       disgo.NewSessionManager(),
	}

	bot.Config.Gateway.IntentSet[disgo.FlagIntentMESSAGE_CONTENT] = true
	bot.Config.Gateway.Intents |= disgo.FlagIntentMESSAGE_CONTENT

	commands.Init(bot, logger)
	database.Init(logger, args.Pause)

	for {
		logger.Info("Connecting to the Discord Gateway...")

		session := disgo.NewSession()

		if err = session.Connect(bot); err != nil {
			logger.Fatalf("Failed to connect to the Discord Gateway: %s", err)
		}

		logger.Info("Connected to the Discord Gateway. Waiting for interactions...")

		time.Sleep(10 * time.Minute)

		session.Disconnect()

		logger.Info("Disconnected from the Discord Gateway.")
	}
}
