package games

import (
	"creft/database"
	"fmt"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleRPSStats(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	var response *disgo.CreateInteractionResponse
	var userID string

	if len(subCommands[0].Options[0].Options) == 0 {
		userID = interaction.Member.User.ID
	} else {
		userID = subCommands[0].Options[0].Options[0].Value.String()
	}

	var found bool = false

	if _, ok := database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[userID]; ok {
		if stats, ok := database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[userID].(map[string]any)["rps"]; ok {
			fmt.Println(stats) // TODO: Finish this

			found = true
		}
	}

	if !found {
		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer(fmt.Sprintf("<@%s> has **not played any rock, paper, scissor games**.\nAll I can say is that they aeren't at the cutting edge of these technologies.", userID)),
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

func HandleRPS(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	switch subCommands[0].Options[0].Name {
	case "stats":
		HandleRPSStats(bot, logger, interaction, subCommands)
	}
}
