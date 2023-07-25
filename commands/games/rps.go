package games

import (
	"creft/database"
	"fmt"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleRPSPlay(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	if !startGame(bot, interaction, logger) {
		return
	}

	response := &disgo.CreateInteractionResponse{
		InteractionID:    interaction.ID,
		InteractionToken: interaction.Token,
		InteractionResponse: &disgo.InteractionResponse{
			Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
			Data: &disgo.Messages{
				Embeds: []*disgo.Embed{
					{
						Title:       disgo.Pointer("Rock, Paper, Scissors"),
						Description: disgo.Pointer(fmt.Sprintf("Choose your weapon <@%s>!", interaction.Member.User.ID)),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
					},
				},
				Components: []disgo.Component{
					disgo.ActionsRow{
						Components: []disgo.Component{
							&disgo.Button{
								Label:    disgo.Pointer("Rock"),
								CustomID: disgo.Pointer("rock"),
								Emoji: &disgo.Emoji{
									Name: disgo.Pointer("U0001F9FA"),
								},
							},
							&disgo.Button{
								Label:    disgo.Pointer("Paper"),
								CustomID: disgo.Pointer("paper"),
								Emoji: &disgo.Emoji{
									Name: disgo.Pointer("U0001F4DC"),
								},
							},
							&disgo.Button{
								Label:    disgo.Pointer("Scissors"),
								CustomID: disgo.Pointer("scissors"),
								Emoji: &disgo.Emoji{
									Name: disgo.Pointer("U00002702"),
								},
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
	case "play":
		HandleRPSPlay(bot, logger, interaction, subCommands)
	case "stats":
		HandleRPSStats(bot, logger, interaction, subCommands)
	}
}
