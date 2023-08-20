package admins

import (
	"github.com/skifli/creft/database"
	"github.com/skifli/creft/utils"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleRemove(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	subCommands := interaction.ApplicationCommand().Options
	userID := subCommands[0].Options[0].Value.String()

	var response *disgo.CreateInteractionResponse
	response = nil

	if !utils.HasAdminPerms(interaction.Member.User.ID) {
		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Missing Permissions"),
							Description: disgo.Pointer("You do not have **sufficient permissions** to run this command.\nYou need to be on the **admins list**."),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				},
			},
		}
	} else {
		if _, ok := database.DatabaseJSON["admins"].(map[string]any)[userID]; ok {
			response = &disgo.CreateInteractionResponse{
				InteractionID:    interaction.ID,
				InteractionToken: interaction.Token,
				InteractionResponse: &disgo.InteractionResponse{
					Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
					Data: &disgo.Messages{
						Embeds: []*disgo.Embed{
							{
								Title:       disgo.Pointer("Success"),
								Description: disgo.Pointer("The specified user was **successfully** removed."),
								Color:       disgo.Pointer(5082199),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							},
						},
					},
				},
			}

			delete(database.DatabaseJSON["admins"].(map[string]any), userID)
			database.Changed = true
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
								Description: disgo.Pointer("The specified user has **not been added**.\nRun **`/admins add {user}`** to add them."),
								Color:       disgo.Pointer(13789294),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							},
						},
					},
				},
			}
		}
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", err)
	} else {
		logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
	}
}
