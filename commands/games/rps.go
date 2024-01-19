package games

import (
	"fmt"
	"strings"

	"github.com/skifli/creft/database"

	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func HandleRPSInteraction(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, customID string) {
	split := strings.Split(customID, "_")
	gameID := split[2]
	choice := split[3]

	game := games[gameID]

	if game == nil {
		response := &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer("Game *does not* **exist** anymore.\n It was *probably* created in a **previous session** of the bot."),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
					Flags: disgo.Pointer(disgo.FlagMessageEPHEMERAL),
				},
			},
		}

		if err := response.Send(bot); err != nil {
			logger.Errorf("Failed to respond to an interaction: %s", err)
		} else {
			logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
		}

		return
	}

	choices := game["choices"].([]string)
	players := game["players"].([]string)
	token := game["token"].(string)

	if (len(choices) == 0 && interaction.Member.User.ID != players[0]) || (len(choices) == 1 && interaction.Member.User.ID != players[1]) {
		response := &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
				Data: &disgo.Messages{
					Embeds: []*disgo.Embed{
						{
							Title:       disgo.Pointer("Error"),
							Description: disgo.Pointer(fmt.Sprintf("It's not your turn <@%s>!", interaction.Member.User.ID)),
							Color:       disgo.Pointer(13789294),
							Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
						},
					},
					Flags: disgo.Pointer(disgo.FlagMessageEPHEMERAL),
				},
			},
		}

		if err := response.Send(bot); err != nil {
			logger.Errorf("Failed to respond to an interaction: %s", err)
		} else {
			logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
		}

		return
	}

	game["choices"] = append(choices, choice)
	choices = game["choices"].([]string)

	if len(game["choices"].([]string)) == 1 {
		response := &disgo.EditOriginalInteractionResponse{
			ApplicationID:    bot.ApplicationID,
			InteractionToken: token,
			Embeds: disgo.Pointer([]*disgo.Embed{{
				Title:       disgo.Pointer("Rock, Paper, Scissors"),
				Description: disgo.Pointer(fmt.Sprintf("Choose your weapon <@%s>!", players[1])),
				Color:       disgo.Pointer(6591981),
				Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
			}}),
			Components: disgo.Pointer([]disgo.Component{&disgo.ActionRow{
				Type: disgo.FlagComponentTypeActionRow,
				Components: []disgo.Component{
					&disgo.Button{
						Type:     disgo.FlagComponentTypeButton,
						Style:    disgo.FlagButtonStylePRIMARY,
						Label:    disgo.Pointer("Rock"),
						CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_rock", gameID)),
						Emoji: &disgo.Emoji{
							ID:   nil,
							Name: disgo.Pointer("\U0001F5FB"),
						},
					},
					&disgo.Button{
						Type:     disgo.FlagComponentTypeButton,
						Style:    disgo.FlagButtonStylePRIMARY,
						Label:    disgo.Pointer("Paper"),
						CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_paper", gameID)),
						Emoji: &disgo.Emoji{
							ID:   nil,
							Name: disgo.Pointer("\U0001F4DC"),
						},
					},
					&disgo.Button{
						Type:     disgo.FlagComponentTypeButton,
						Style:    disgo.FlagButtonStylePRIMARY,
						Label:    disgo.Pointer("Scissors"),
						CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_scissors", gameID)),
						Emoji: &disgo.Emoji{
							ID:   nil,
							Name: disgo.Pointer("\U00002702"),
						},
					},
				},
			}}),
		}

		if _, err := response.Send(bot); err != nil {
			logger.Errorf("Failed to respond to an interaction: %s", err)
		} else {
			logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
		}
	} else if len(game["choices"].([]string)) == 2 {
		winner := -1

		if choices[0] == choices[1] {
			winner = -1
		} else if choices[0] == "rock" && choices[1] == "paper" {
			winner = 1
		} else if choices[0] == "rock" && choices[1] == "scissors" {
			winner = 0
		} else if choices[0] == "paper" && choices[1] == "rock" {
			winner = 0
		} else if choices[0] == "paper" && choices[1] == "scissors" {
			winner = 1
		} else if choices[0] == "scissors" && choices[1] == "rock" {
			winner = 1
		} else if choices[0] == "scissors" && choices[1] == "paper" {
			winner = 0
		}

		description := fmt.Sprintf("<@%s> chose **%s**.\n<@%s> chose **%s**.\n\n", players[0], choices[0], players[1], choices[1])

		if winner == -1 {
			description = fmt.Sprintf("%s<@%s> and <@%s> **tied**!", description, players[0], players[1])
		} else {
			description = fmt.Sprintf("%s<@%s> **won** against <@%s>!", description, players[winner], players[1-winner])
		}

		response := &disgo.EditOriginalInteractionResponse{
			ApplicationID:    bot.ApplicationID,
			InteractionToken: token,
			Embeds: disgo.Pointer([]*disgo.Embed{{
				Title:       disgo.Pointer("Game Results!"),
				Description: disgo.Pointer(description),
				Color:       disgo.Pointer(6591981),
				Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
			}}),
			Components: disgo.Pointer([]disgo.Component{}),
		}

		if _, err := response.Send(bot); err != nil {
			logger.Errorf("Failed to respond to an interaction: %s", err)
		} else {
			logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
		}

		if _, ok := database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]]; !ok {
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]] = map[string]any{"rps": map[string]any{"gamesPlayed": 0.0, "gamesWon": 0.0, "gamesLost": 0.0}}
		}

		if _, ok := database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]]; !ok {
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]] = map[string]any{"rps": map[string]any{"gamesPlayed": 0.0, "gamesWon": 0.0, "gamesLost": 0.0}}
		}

		database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]].(map[string]any)["rps"].(map[string]any)["gamesPlayed"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]].(map[string]any)["rps"].(map[string]any)["gamesPlayed"].(float64) + 1
		database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]].(map[string]any)["rps"].(map[string]any)["gamesPlayed"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]].(map[string]any)["rps"].(map[string]any)["gamesPlayed"].(float64) + 1

		if winner != -1 {
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[winner]].(map[string]any)["rps"].(map[string]any)["gamesWon"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[winner]].(map[string]any)["rps"].(map[string]any)["gamesWon"].(float64) + 1
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1-winner]].(map[string]any)["rps"].(map[string]any)["gamesLost"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1-winner]].(map[string]any)["rps"].(map[string]any)["gamesLost"].(float64) + 1
		} else if winner == -1 {
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]].(map[string]any)["rps"].(map[string]any)["gamesDrawn"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[0]].(map[string]any)["rps"].(map[string]any)["gamesDrawn"].(float64) + 1
			database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]].(map[string]any)["rps"].(map[string]any)["gamesDrawn"] = database.DatabaseJSON["games"].(map[string]any)["userStats"].(map[string]any)[players[1]].(map[string]any)["rps"].(map[string]any)["gamesDrawn"].(float64) + 1
		}

		database.DatabaseJSON["games"].(map[string]any)["rps"].(map[string]any)["gamesCount"] = database.DatabaseJSON["games"].(map[string]any)["rps"].(map[string]any)["gamesCount"].(float64) + 1

		database.Changed = true

		endGame(game["interaction"].(*disgo.InteractionCreate))
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
						Description: disgo.Pointer(fmt.Sprintf("You chose **%s** <@%s>!", choice, interaction.Member.User.ID)),
						Color:       disgo.Pointer(6591981),
						Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
					},
				},
				Flags: disgo.Pointer(disgo.FlagMessageEPHEMERAL),
			},
		},
	}

	if err := response.Send(bot); err != nil {
		logger.Errorf("Failed to respond to an interaction: %s", err)
	} else {
		logger.Infof("Responded to an interaction from %s.", interaction.Member.User.Username)
	}
}

func HandleRPSPlay(bot *disgo.Client, logger *golog.Logger, interaction *disgo.InteractionCreate, subCommands []*disgo.ApplicationCommandInteractionDataOption) {
	if !startGame(bot, interaction, logger) {
		return
	}

	gameID := interaction.Token[:10]

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
					&disgo.ActionRow{
						Type: disgo.FlagComponentTypeActionRow,
						Components: []disgo.Component{
							&disgo.Button{
								Type:     disgo.FlagComponentTypeButton,
								Style:    disgo.FlagButtonStylePRIMARY,
								Label:    disgo.Pointer("Rock"),
								CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_rock", gameID)),
								Emoji: &disgo.Emoji{
									ID:   nil,
									Name: disgo.Pointer("\U0001F5FB"),
								},
							},
							&disgo.Button{
								Type:     disgo.FlagComponentTypeButton,
								Style:    disgo.FlagButtonStylePRIMARY,
								Label:    disgo.Pointer("Paper"),
								CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_paper", gameID)),
								Emoji: &disgo.Emoji{
									ID:   nil,
									Name: disgo.Pointer("\U0001F4DC"),
								},
							},
							&disgo.Button{
								Type:     disgo.FlagComponentTypeButton,
								Style:    disgo.FlagButtonStylePRIMARY,
								Label:    disgo.Pointer("Scissors"),
								CustomID: disgo.Pointer(fmt.Sprintf("games_rps_%s_scissors", gameID)),
								Emoji: &disgo.Emoji{
									ID:   nil,
									Name: disgo.Pointer("\U00002702"),
								},
							},
						},
					},
				},
			},
		},
	}

	games[gameID] = map[string]any{"game": "rps", "interaction": interaction, "choices": []string{}, "players": []string{interaction.Member.User.ID, subCommands[0].Options[0].Options[0].Value.String()}, "token": interaction.Token}

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
			response = &disgo.CreateInteractionResponse{
				InteractionID:    interaction.ID,
				InteractionToken: interaction.Token,
				InteractionResponse: &disgo.InteractionResponse{
					Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
					Data: &disgo.Messages{
						Embeds: []*disgo.Embed{
							{
								Title:       disgo.Pointer("Rock, Paper, Scissors"),
								Description: disgo.Pointer(fmt.Sprintf("<@%s> has played **%d** game(s).\nThey have won **%d** game(s) and lost **%d** game(s).", userID, int(stats.(map[string]any)["gamesPlayed"].(float64)), int(stats.(map[string]any)["gamesWon"].(float64)), int(stats.(map[string]any)["gamesLost"].(float64)))),
								Color:       disgo.Pointer(6591981),
								Footer:      &disgo.EmbedFooter{Text: "Run /about for more information about the bot."},
							},
						},
						Flags: disgo.Pointer(disgo.FlagMessageEPHEMERAL),
					},
				},
			}

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
					Flags: disgo.Pointer(disgo.FlagMessageEPHEMERAL),
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
