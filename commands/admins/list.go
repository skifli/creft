package admins

import (
	"creft/database"
	"fmt"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleList(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	admins := ""

	for admin := range database.DatabaseJSON["admins"].(map[string]any) {
		admins += fmt.Sprintf("<@%s>\n", admin)
	}

	response := &disgo.CreateInteractionResponse{
		InteractionID:    interaction.ID,
		InteractionToken: interaction.Token,
		InteractionResponse: &disgo.InteractionResponse{
			Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
			Data: &disgo.Messages{
				Embeds: []*disgo.Embed{
					{
						Title:       disgo.Pointer("Admins"),
						Description: disgo.Pointer(fmt.Sprintf("The following users are admins:\n%s", admins)),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
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
