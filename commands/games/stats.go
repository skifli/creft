package games

import (
	"creft/database"
	"fmt"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleStatsServer(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	rpsStats := fmt.Sprintf("Games Played: %d", uint64(database.DatabaseJSON["games"].(map[string]any)["rps"].(map[string]any)["gamesCount"].(float64)))

	response := &disgo.CreateInteractionResponse{
		InteractionID:    interaction.ID,
		InteractionToken: interaction.Token,
		InteractionResponse: &disgo.InteractionResponse{
			Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
			Data: &disgo.Messages{
				Embeds: []*disgo.Embed{
					{
						Title:       disgo.Pointer("Server Game Stats"),
						Description: disgo.Pointer(fmt.Sprintf("Stats for all games on this server. There are **%d** game(s) in total.", len(database.DatabaseJSON["games"].(map[string]any))-1)),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						Fields: []*disgo.EmbedField{
							{
								Name:  "Rock Paper Scissors",
								Value: rpsStats,
							},
						},
					},
				},
			},
		},
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", err)
	} else {
		logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
	}
}

func HandleStatsUser(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	var response *disgo.CreateInteractionResponse
	var userID string

	if len(subCommands[0].Options[0].Options) == 0 {
		userID = interaction.Member.User.ID
	} else {
		userID = subCommands[0].Options[0].Options[0].Value.String()
	}

	if stats, ok := database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[userID].(map[string]any); ok {
		rpsStats := fmt.Sprintf("Games Played: %d\nWins: %d\nLosses: %d\nDraws: %d", uint64(stats["rps"].(map[string]any)["gamesPlayed"].(float64)), uint64(stats["rps"].(map[string]any)["gamesWon"].(float64)), uint64(stats["rps"].(map[string]any)["gamesLost"].(float64)), uint64(stats["rps"].(map[string]any)["gamesDrawn"].(float64)))

		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("User Game Stats"),
							Description: disgo.Pointer(fmt.Sprintf("Stats for all games played by <@%s>.", userID)),

							Color:  disgo.Pointer(6591981),
							Footer: &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							Fields: []*disgo.EmbedField{
								{
									Name:  "Rock Paper Scissors",
									Value: rpsStats,
								},
							},
						},
					},
				},
			},
		}
	} else {
		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer(fmt.Sprintf("<@%s> has **not played any games**.\nAnd if you haven't, what are you waiting for?", userID)),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				},
			},
		}
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", err)
	} else {
		logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
	}
}

func HandleStats(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	switch subCommands[0].Options[0].Name {
	case "server":
		HandleStatsServer(bot, logger, interaction)
	case "user":
		HandleStatsUser(bot, logger, interaction, subCommands)
	}
}
