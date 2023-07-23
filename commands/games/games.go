package games

import (
	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	subCommands := interaction.ApplicationCommand().Options

	switch subCommands[0].Name {
	case "stats":
		HandleStats(bot, logger, interaction, subCommands)
	}
}

func Init(bot *disgo.Client, logger *golog.Logger) {
	request := disgo.CreateGlobalApplicationCommand{
		Name:        "games",
		Description: disgo.Pointer("Commands relating to games."),
		Options: []*disgo.ApplicationCommandOption{
			{
				Name:        "stats",
				Description: "Views game stats.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP,
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "server",
						Description: "Views game stats for the server.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
					},
					{
						Name:        "user",
						Description: "Views game stats for a user.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
						Options: []*disgo.ApplicationCommandOption{
							{
								Name:        "user",
								Description: "The user to view game stats for.",
								Type:        disgo.FlagApplicationCommandOptionTypeUSER,
								Required:    disgo.Pointer(true),
							},
						},
					},
				},
			},
		},
	}

	_, err := request.Send(bot)

	if err != nil {
		logger.Fatalf("Failed to add slash command to bot: %s", err)
	}
}
