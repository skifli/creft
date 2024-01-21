package counting

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/skifli/creft/database"

	"github.com/maja42/goval"
	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

const COOLDOWN_TIME_SECONDS = 5

var messageTracker = make(map[string]time.Time)

func checkIfOnCooldown(userID string) bool {
	if _, ok := messageTracker[userID]; ok {
		if time.Since(messageTracker[userID]) < time.Second*COOLDOWN_TIME_SECONDS {
			return true
		} else {
			delete(messageTracker, userID)
		}
	}

	return false
}

func addCooldown(userID string) {
	messageTracker[userID] = time.Now()
}

func onMessageCreate(bot *disgo.Client, logger *golog.Logger, message *disgo.MessageCreate) {
	if checkIfOnCooldown(message.Author.ID) {
		logger.Infof("Ignored a counting message from %s due to cooldown.", message.Author.Username)

		return
	}

	defer func() {
		if r := recover(); r != nil {
			var response *disgo.CreateMessage
			response = nil

			if strings.Contains(r.(runtime.Error).Error(), "integer divide by zero") {
				response = &disgo.CreateMessage{
					ChannelID: message.ChannelID,
					MessageReference: &disgo.MessageReference{
						MessageID: disgo.Pointer(message.ID),
						ChannelID: disgo.Pointer(message.ChannelID),
						GuildID:   message.GuildID,
					},
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Cheeky"),
							Description: disgo.Pointer("Imagine trying to divide by **`0`**. Couldn't be you.\nThe count has stayed the same."),
							Color:       disgo.Pointer(6591981),
							Footer:      &disgo.EmbedFooter{Text: "Nothing eventful happened. Run /about for more information about the bot."},
						},
					},
				}
			} else {
				response = &disgo.CreateMessage{
					ChannelID: message.ChannelID,
					MessageReference: &disgo.MessageReference{
						MessageID: disgo.Pointer(message.ID),
						ChannelID: disgo.Pointer(message.ChannelID),
						GuildID:   message.GuildID,
					},
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer(fmt.Sprintf("An unknown error occurred: **`%s`**.\nThe count has stayed the same.", r.(string))),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
				}
			}

			if _, err := response.Send(bot); err != nil {
				logger.Errorf("Failed to respond to a message: %s", err)
			} else {
				logger.Infof("Responded to a message from %s.", message.Author.Username)
			}
		}
	}()

	if channelDatabase, ok := database.DatabaseJSON["counting"].(map[string]any)[message.ChannelID].(map[string]any); ok {
		expression := goval.NewEvaluator()

		if result, err := expression.Evaluate(message.Content, nil, nil); err == nil {
			var response *disgo.CreateMessage
			response = nil

			if message.Author.ID == channelDatabase["lastCountUserID"] {
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
							Description: disgo.Pointer("You **counted last**. Please wait for **someone else** to count!\nYou now have a **5 second** cooldown."),
							Color:       disgo.Pointer(6591981),
							Footer:      &disgo.EmbedFooter{Text: "Idk if that was even correct, but I will let it slide. Run /about for more information about the bot."},
						},
					},
				}

				addCooldown(message.Author.ID)
				logger.Infof("%s is now on a counting cooldown for %d seconds.", message.Author.Username, COOLDOWN_TIME_SECONDS)
			} else {
				failed := false
				var value float64 = 0

				switch v := result.(type) {
				case int:
					value = float64(v)
				case float32:
				case float64:
					value = math.Round(v)
				default:
					failed = true
				}

				count := channelDatabase["count"].(float64) + 1.0

				if failed || value != count {
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
								Description: disgo.Pointer(fmt.Sprintf("The correct number was **`%0.f`**, but you said **`%0.f`**. The count has reset to **`0`**.", count, value)),
								Color:       disgo.Pointer(13789294),
								Footer:      &disgo.EmbedFooter{Text: "RIP streak. Run /about for more information about the bot."},
							},
						},
					}

					database.DatabaseJSON["counting"].(map[string]any)[message.ChannelID] = map[string]any{"count": 0.0, "countMax": channelDatabase["countMax"].(float64), "lastCountUserID": message.Author.ID, "lastCountMessageID": message.ID, "resetsCount": channelDatabase["resetsCount"].(float64) + 1}
					database.Changed = true
				} else {
					channelDatabase["count"] = count
					channelDatabase["lastCountUserID"] = message.Author.ID
					channelDatabase["lastCountMessageID"] = message.ID

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
						logger.Errorf("Failed to react to a message: %s", err)
					} else {
						logger.Infof("Reacted to a message from %s.", message.Author.Username)
					}
				}
			}

			if response != nil {
				if _, err := response.Send(bot); err != nil {
					logger.Errorf("Failed to respond to a message: %s", err)
				} else {
					logger.Infof("Responded to a message from %s.", message.Author.Username)
				}
			}
		}
	}
}

func onMessageDelete(bot *disgo.Client, logger *golog.Logger, message *disgo.MessageDelete) {
	if channelDatabase, ok := database.DatabaseJSON["counting"].(map[string]any)[message.ChannelID].(map[string]any); ok {
		if channelDatabase["lastCountMessageID"] == message.MessageID {
			response := &disgo.CreateMessage{
				ChannelID: message.ChannelID,
				Embeds: []*disgo.Embed{
					{
						Title:       disgo.Pointer("I saw that"),
						Description: disgo.Pointer(fmt.Sprintf("<@%s> **deleted** their message. The count is at **`%.0f`**.", channelDatabase["lastCountUserID"], channelDatabase["count"])),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Cheeky! Run /about for more information about the bot."},
					},
				},
			}

			if _, err := response.Send(bot); err != nil {
				logger.Errorf("Failed to respond to a deleted message: %s", err)
			} else {
				logger.Infof("Responded to a deleted message from %s.", channelDatabase["lastCountMessageAuthor"])
			}
		}
	}
}

func Handle(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate) {
	subCommands := interaction.ApplicationCommand().Options

	var channelID string

	if len(subCommands[0].Options) == 0 {
		channelID = *interaction.ChannelID
	} else {
		channelID = subCommands[0].Options[0].Value.String()
	}

	switch subCommands[0].Name {
	case "add":
		HandleAdd(bot, logger, interaction, channelID)
	case "remove":
		HandleRemove(bot, logger, interaction, channelID)
	case "stats":
		HandleStats(bot, logger, interaction, channelID)
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
		logger.Fatalf("Failed to add slash command to bot: %s", err)
	}

	if err := bot.Handle(disgo.FlagGatewayEventNameMessageCreate, func(message *disgo.MessageCreate) { onMessageCreate(bot, logger, message) }); err != nil {
		logger.Fatalf("Failed to add event handler to bot: %s", err)
	} else if err := bot.Handle(disgo.FlagGatewayEventNameMessageDelete, func(message *disgo.MessageDelete) { onMessageDelete(bot, logger, message) }); err != nil {
		logger.Fatalf("Failed to add event handler to bot: %s", err)
	}
}
