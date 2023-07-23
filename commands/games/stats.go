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
	userID := subCommands[0].Options[0].Options[0].Value.String()

	if _, ok := database.DatabaseJSON["users"].(map[string]any)[userID]; !ok {
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
