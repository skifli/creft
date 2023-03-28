package counting

import (
	"creft/database"
	"creft/utils"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleAdd(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, channel string) {
	var response *disgo.CreateInteractionResponse
	response = nil

	if !utils.HasAdminPerms(interaction.Member.Permissions) {
		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Missing Permissions"),
							Description: disgo.Pointer("You do not have **sufficient permissions** to run this command.\nYou need a role with the **Administrator Permission**."),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				},
			},
		}
	} else {
		if _, ok := database.DatabaseJSON["counting"][channel]; ok {
			response = &disgo.CreateInteractionResponse{
				InteractionID:    interaction.ID,
				InteractionToken: interaction.Token,
				InteractionResponse: &disgo.InteractionResponse{
					Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
					Data: &disgo.Messages{
						Embeds: []*disgo.Embed{
							{
								Title:       disgo.Pointer("Error"),
								Description: disgo.Pointer("The specified channel has **already been added**.\nRun **`/counting remove {channel}`** to remove it."),
								Color:       disgo.Pointer(13789294),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
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
								Title:       disgo.Pointer("Success"),
								Description: disgo.Pointer("The specified channel was **successfully** added."),
								Color:       disgo.Pointer(5082199),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							},
						},
					},
				},
			}

			database.DatabaseJSON["counting"][channel] = map[string]any{"count": 0.0, "countMax": 0.0, "lastCountUserID": "", "lastCountMessageID": "", "resetsCount": 0.0}
			database.Changed = true
		}
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", err)
	} else {
		logger.Infof("Responded to an interaction from %s#%s.", interaction.Member.User.Username, interaction.Member.User.Discriminator)
	}
}
