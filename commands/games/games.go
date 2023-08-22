package games

import (
	"strings"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

var userGames map[string]struct{} = make(map[string]struct{})
var games map[string]map[string]any = make(map[string]map[string]any)

func startGame(bot *disgo.Client, interaction *disgo.InteractionCreate, logger *golog.Logger) bool {
	if _, ok := userGames[interaction.Member.User.ID]; ok {
		response := &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer("You are **already in a game**. You can't play multiple games at once!\nPlease **finish your current game** before starting another one."),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				},
			},
		}

		if err := response.Send(bot); err != nil {
			logger.Errorf("Failed to send slash command response: %s", err)
		} else {
			logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
		}

		return false
	}

	userGames[interaction.Member.User.ID] = struct{}{}
	return true
}

func endGame(interaction *disgo.InteractionCreate) {
	delete(userGames, interaction.Member.User.ID)
}

func HandleInteraction(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, interactionCustomID string) {
	if strings.HasPrefix(interactionCustomID, "games_rps") {
		HandleRPSInteraction(bot, logger, interaction, interactionCustomID)
	} else {
		// TODO: Add error
	}
}

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	subCommands := interaction.ApplicationCommand().Options

	switch subCommands[0].Name {
	case "rps":
		HandleRPS(bot, logger, interaction, subCommands)
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
				Description: "View game stats.",
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
							},
						},
					},
				},
			},
			{
				Name:        "rps",
				Description: "The classic rock paper scissors game.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP,
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "play",
						Description: "Play rock paper scissors.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
						Options: []*disgo.ApplicationCommandOption{
							{
								Name:        "user",
								Description: "The user to play against.",
								Type:        disgo.FlagApplicationCommandOptionTypeUSER,
								Required:    disgo.Pointer(true),
							},
						},
					},
					{
						Name:        "stats",
						Description: "View a user's rock paper scissors stats.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
						Options: []*disgo.ApplicationCommandOption{
							{
								Name:        "user",
								Description: "The user to view rock paper scissors stats for.",
								Type:        disgo.FlagApplicationCommandOptionTypeUSER,
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
