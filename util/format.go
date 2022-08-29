package util

import (
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"strings"
)

// YesNoBool creates a string of "Yes" or "No" depending on whether the boolean parameter is true or false
func YesNoBool(b bool) string {
	if b {
		return "Yes"
	} else {
		return "No"
	}
}

func UserTag(u discord.User) string {
	return fmt.Sprintf("%s (`%s` | `%d`)", u.Mention(), u.Tag(), u.ID)
}

func ChannelTag(c discord.ChannelID) string {
	return fmt.Sprintf("%s (`%s`)", c.Mention(), c.String())
}

func DiscordJumpLink(
	guildID discord.GuildID,
	channelID discord.ChannelID,
	messageID discord.MessageID,
) (string, error) {
	var link strings.Builder

	if guildID == discord.NullGuildID || channelID == discord.NullChannelID {
		return "", errors.New("cannot create jump link without guild and channel IDs")
	} else {
		link.WriteString("https://discord.com/channels/")
		link.WriteString(guildID.String())
		link.WriteString("/")
		link.WriteString(channelID.String())
		link.WriteString("/")
	}

	if messageID != discord.NullMessageID {
		link.WriteString(messageID.String())
	}

	return link.String(), nil
}
