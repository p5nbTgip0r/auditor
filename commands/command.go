package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

//type Command struct {
//	*api.CreateCommandData
//
//	Handle func(*Interaction)
//}
//
//func BuildCommand(name, description string, options ...discord.CommandOption) *Command {
//	cmdData := &api.CreateCommandData{
//		Name:        name,
//		Description: description,
//		Options:     options,
//	}
//
//	return &Command{CreateCommandData: cmdData}
//}

type Interaction struct {
	*gateway.InteractionCreateEvent

	CommandInteraction *discord.CommandInteraction
	State              *state.State
}

func (i Interaction) DeferResponse(ephemeral bool) error {
	var flags discord.MessageFlags
	if ephemeral {
		flags = discord.EphemeralMessage
	}

	return i.State.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
		Type: api.DeferredMessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Flags: flags,
		},
	})
}

func (i Interaction) EditTextResponse(message string) (*discord.Message, error) {
	return i.State.EditInteractionResponse(i.AppID, i.Token, api.EditInteractionResponseData{
		Content: option.NewNullableString(message),
	})
}
