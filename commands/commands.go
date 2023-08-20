package commands

import (
	"strings"

	"github.com/skifli/creft/commands/about"
	"github.com/skifli/creft/commands/admins"
	"github.com/skifli/creft/commands/counting"
	"github.com/skifli/creft/commands/games"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	defer func() {
		if err := recover(); err != nil {
			interactionCustomID := interaction.MessageComponent().CustomID

			if strings.HasPrefix(interactionCustomID, "games_") {
				games.HandleInteraction(bot, logger, interaction, interactionCustomID)
			} else {
				// TODO: Add error
			}
		}
	}()

	commmandName := interaction.ApplicationCommand().Name

	switch commmandName {
	case "about":
		about.Handle(bot, logger, interaction)
	case "admins":
		admins.Handle(bot, logger, interaction)
	case "counting":
		counting.Handle(bot, logger, interaction)
	case "games":
		games.Handle(bot, logger, interaction)
	}
}

func Init(bot *disgo.Client, logger *golog.Logger) {
	about.Init(bot, logger)
	admins.Init(bot, logger)
	counting.Init(bot, logger)
	games.Init(bot, logger)

	if err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(interaction *disgo.InteractionCreate) { Handle(bot, logger, interaction) }); err != nil {
		logger.Fatalf("Failed to add event handler to bot: %s", err)
	}
}
