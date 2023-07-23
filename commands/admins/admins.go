package admins

import (
	"github.com/skifli/golog"
	"github.com/switchupcb/disgo"
)

func Init(bot *disgo.Client, logger *golog.Logger) {
	request := disgo.CreateGlobalApplicationCommand{
		Name:        "admins",
		Description: disgo.Pointer("Commands relating to admins."),
		Type:        disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
		Options: []*disgo.ApplicationCommandOption{
			{
				Name:        "list",
				Description: "Lists all admins.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,
			},
		},
	}

	_, err := request.Send(bot)

	if err != nil {
		logger.Fatalf("Failed to add slash command to bot: %s", err)
	}
}
