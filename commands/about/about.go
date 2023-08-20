package about

import (
	"fmt"
	"time"

	"github.com/skifli/creft/utils"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	response := &disgo.CreateInteractionResponse{
		InteractionID:    interaction.ID,
		InteractionToken: interaction.Token,
		InteractionResponse: &disgo.InteractionResponse{
			Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
			Data: &disgo.Messages{
				Embeds: []*disgo.Embed{
					{
						Title:       disgo.Pointer("About"),
						Description: disgo.Pointer(fmt.Sprintf("This instance has been online for **`%.2f`** seconds.\nBot created by **skifli#8774**.\nGitHub Repository: https://github.com/skifli/creft.", time.Since(utils.StartTime).Seconds())),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Hello there."},
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
}

func Init(bot *disgo.Client, logger *golog.Logger) {
	request := disgo.CreateGlobalApplicationCommand{
		Name:        "about",
		Description: disgo.Pointer("About the bot."),
		Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
	}

	_, err := request.Send(bot)

	if err != nil {
		logger.Fatalf("Failed to add slash command to bot: %s", err)
	}
}
