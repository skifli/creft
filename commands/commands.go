package commands

import (
	"creft/commands/about"
	"creft/commands/counting"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	commmandName := interaction.ApplicationCommand().Name

	switch commmandName {
	case "about":
		about.Handle(bot, logger, interaction)
	case "counting":
		counting.Handle(bot, logger, interaction)
	}
}

func Init(bot *disgo.Client, logger *golog.Logger) {
	about.Init(bot, logger)
	counting.Init(bot, logger)

	if err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(interaction *disgo.InteractionCreate) { Handle(bot, logger, interaction) }); err != nil {
		logger.Fatalf("Failed to add event handler to bot: %s", err)
	}
}
