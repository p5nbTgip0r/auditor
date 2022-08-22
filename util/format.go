package util

import (
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"strings"
)

func Plural(count int) string {
	return CPlural(count, "s")
}

func CPlural(count int, plural string) string {
	if count != 1 {
		return plural
	}

	return ""
}

func FullTag(u discord.User) string {
	return fmt.Sprintf("%s (`%s` | `%d`)", u.Mention(), u.Tag(), u.ID)
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
