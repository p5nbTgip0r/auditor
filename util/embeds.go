package util

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
)

// AddField generates a discord.EmbedField from the given name, value, and codeTags parameters, then appends that
// field to the embed.
//
// The value of the discord.EmbedField looks like:
//
//	"value"
//
// If codeTags is true, the given value will be surrounded in `code tags`:
//
//	"`value`"
//
// This is mainly for cutting down on the boilerplate code with generating a value string, appending, and assigning
// to discord.Embed.Fields
func AddField(embed *discord.Embed, name, value string, codeTags bool) {
	v := value
	if codeTags {
		v = fmt.Sprintf("`%s`", value)
	}
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:  name,
			Value: v,
		})
}

// AddUpdatedField generates a discord.EmbedField from the given name, old + new values, and codeTags parameters, then
// appends that field to the embed.
//
// The value of the discord.EmbedField looks like:
//
//	"old -> new"
//
// If codeTags is true, the given values will be surrounded in `code tags`:
//
//	"`old` -> `new`"
//
// This is mainly for cutting down on the boilerplate code with generating a value string, appending, and assigning
// to discord.Embed.Fields
func AddUpdatedField(embed *discord.Embed, name, old, new string, codeTags bool) {
	format := "%s -> %s"
	if codeTags {
		format = "`%s` -> `%s`"
	}
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:  name,
			Value: fmt.Sprintf(format, old, new),
		})
}
