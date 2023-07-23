package admins

import (
	"creft/utils"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleAdd(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
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
	}
}
