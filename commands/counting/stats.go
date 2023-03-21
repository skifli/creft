package counting

import (
	"creft/database"
	"fmt"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleStats(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, channel string) {
	var response *disgo.CreateInteractionResponse
	response = nil

	if channelDatabase, ok := database.DatabaseJSON["counting"][channel].(map[string]any); ok {
		response = &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer(fmt.Sprintf("Info for <#%s>", channel)),
							Description: disgo.Pointer(fmt.Sprintf("The current count is **`%d`**.\nThe highest count that has been reached is **`%d`**.\nThe count has been reset **`%d`** time(s).\nThe last user to count was <@%s>. They counted at [this message](https://discord.com/channels/%s/%s/%s)", uint64(channelDatabase["count"].(float64)), uint64(channelDatabase["countMax"].(float64)), uint64(channelDatabase["resetsCount"].(float64)), channelDatabase["lastCountUserID"].(string), *interaction.GuildID, channel, channelDatabase["lastCountMessageID"].(string))),
							Color:       disgo.Pointer(5082199),
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
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer("The specified channel has **not been added**.\nRun **`/counting add {channel}`** to add it."),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				},
			},
		}
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", nil, err)
	} else {
		logger.Infof("Responded to an interaction from %s#%s.", nil, interaction.Member.User.Username, interaction.Member.User.Discriminator)
	}
}
