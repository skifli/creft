package counting

import (
	"creft/database"
	"fmt"

	"github.com/maja42/goval"
	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func onMessage(bot *disgo.Client, logger *golog.Logger, message *disgo.MessageCreate) {
	if channelDatabase, ok := database.DatabaseJSON["counting"][message.ChannelID].(map[string]any); ok {
		expression := goval.NewEvaluator()

		if result, err := expression.Evaluate(message.Content, nil, nil); err == nil {
			var response *disgo.CreateMessage
			response = nil

			if message.Author.ID == channelDatabase["lastUser"] {
				response = &disgo.CreateMessage{
					ChannelID: message.ChannelID,
					MessageReference: &disgo.MessageReference{
						MessageID: disgo.Pointer(message.ID),
						ChannelID: disgo.Pointer(message.ChannelID),
						GuildID:   message.GuildID,
					},
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Please Wait"),
							Description: disgo.Pointer("You **counted last**. Please wait for **someone else** to count!"),
							Color:       disgo.Pointer(6591981),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				}
			} else {
				failed := false
				var value float64 = 0

				switch v := result.(type) {
				case int:
					value = float64(v)
				default:
					failed = true
				}

				count := channelDatabase["count"].(float64)

				if failed || value != count+1.0 {
					response = &disgo.CreateMessage{
						ChannelID: message.ChannelID,
						MessageReference: &disgo.MessageReference{
							MessageID: disgo.Pointer(message.ID),
							ChannelID: disgo.Pointer(message.ChannelID),
							GuildID:   message.GuildID,
						},
						Embeds: []*disgo.Embed{
							{
								Title:       disgo.Pointer("Incorrect"),
								Description: disgo.Pointer(fmt.Sprintf("The correct number was **`%0.f`**. The count has reset to **`0`**.", count+1)),
								Color:       disgo.Pointer(13789294),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							},
						},
					}

					channelDatabase["count"] = 0.0
					channelDatabase["resetsCount"] = channelDatabase["resetsCount"].(float64) + 1.0
					channelDatabase["lastUser"] = ""
					database.Changed = true
				} else {
					channelDatabase["count"] = count + 1.0
					channelDatabase["lastUser"] = message.Author.ID

					if count > channelDatabase["countMax"].(float64) {
						channelDatabase["countMax"] = count
					}

					database.Changed = true

					reaction := &disgo.CreateReaction{
						ChannelID: message.ChannelID,
						MessageID: message.ID,
						Emoji:     "\U00002705",
					}

					if err := reaction.Send(bot); err != nil {
						logger.Errorf("Failed to react to a message: %s", nil, err)
					} else {
						logger.Infof("Reacted to a message from %s#%s.", nil, message.Author.Username, message.Author.Discriminator)
					}
				}
			}

			if response != nil {
				if _, err := response.Send(bot); err != nil {
					logger.Errorf("Failed to send message response: %s", nil, err)
				} else {
					logger.Infof("Responded to a message from %s#%s.", nil, message.Author.Username, message.Author.Discriminator)
				}
			}
		}
	}
}

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	subCommands := interaction.ApplicationCommand().Options

	switch subCommands[0].Name {
	case "add":
		HandleAdd(bot, logger, interaction, subCommands[0].Options[0].Value.String())
	case "remove":
		HandleRemove(bot, logger, interaction, subCommands[0].Options[0].Value.String())
	case "stats":
		HandleStats(bot, logger, interaction, subCommands[0].Options[0].Value.String())
	}
}

func Init(bot *disgo.Client, logger *golog.Logger) {
	request := disgo.CreateGlobalApplicationCommand{
		Name:        "counting",
		Description: disgo.Pointer("Commands relating to counting channels."),
		Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		Options: []*disgo.ApplicationCommandOption{
			{
				Name:        "add",
				Description: "Add a counting channel.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "The channel to add.",
						Type:        disgo.FlagApplicationCommandOptionTypeCHANNEL,
					},
				},
			},
			{
				Name:        "remove",
				Description: "Remove a counting channel.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "The channel to remove.",
						Type:        disgo.FlagApplicationCommandOptionTypeCHANNEL,
					},
				},
			},
			{
				Name:        "stats",
				Description: "View stats about a counting channel.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "The channel to view the stats of.",
						Type:        disgo.FlagApplicationCommandOptionTypeCHANNEL,
					},
				},
			},
		},
	}

	_, err := request.Send(bot)

	if err != nil {
		logger.Fatalf("Failed to add slash command to bot: %s", nil, err)
	}

	if err := bot.Handle(disgo.FlagGatewayEventNameMessageCreate, func(message *disgo.MessageCreate) { onMessage(bot, logger, message) }); err != nil {
		logger.Fatalf("Failed to add event handler to bot: %s", nil, err)
	}
}
