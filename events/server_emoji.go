package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strings"
)

func init() {
	handle := func(c gateway.GuildEmojisUpdateEvent, old []discord.Emoji) {
		if !check(audit.ServerEmoji, &c.GuildID, nil) {
			return
		}

		genField := func(emojis map[discord.EmojiID]discord.Emoji, name string, link bool) discord.EmbedField {
			var text strings.Builder
			for _, emoji := range emojis {
				if text.Len() != 0 {
					text.WriteString("\n")
				}

				if link {
					text.WriteString(fmt.Sprintf("[`:%s:`](%s)", emoji.Name, emoji.EmojiURL()))
				} else {
					text.WriteString(fmt.Sprintf("%s `:%s:`", emoji.String(), emoji.Name))
				}
			}

			return discord.EmbedField{Name: name, Value: text.String()}
		}

		var fields []discord.EmbedField

		added, removed := util.SliceDiffIdentifier(old, c.Emojis, func(emoji discord.Emoji) discord.EmojiID {
			return emoji.ID
		})

		if len(added) != 0 {
			fields = append(fields, genField(added, "Added emojis", false))
		}

		if len(removed) != 0 {
			fields = append(fields, genField(removed, "Removed emojis", true))
		}

		if len(fields) != 0 {
			bot.QueueEmbed(audit.ServerEmoji, c.GuildID, discord.Embed{
				Description: "**:pencil: Server's emojis updated!**",
				Timestamp:   discord.NowTimestamp(),
				Color:       color.Gold,
				Fields:      fields,
			})
		}
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildEmojisUpdateEvent) {
			o, err := s.EmojiStore.Emojis(c.GuildID)
			if err != nil {
				go handleError(
					audit.ServerEmoji,
					c.GuildID,
					err,
					"Could not retrieve guild from cache: `"+c.GuildID.String()+"`",
					discord.User{},
				)
				return
			}

			go handle(*c, o)
		})
	})
}
