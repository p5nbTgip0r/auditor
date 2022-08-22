package events

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
	"strings"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildEmojisUpdateEvent) {
			if !AuditServerEmoji.check(&c.GuildID, nil) {
				return
			}

			genField := func(emojis []discord.Emoji, name string, link bool) discord.EmbedField {
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

			o, err := s.EmojiStore.Emojis(c.GuildID)
			if err != nil {
				// todo: proper warning for this
				log.Warn().Err(err).Interface("event", c).Msg("Could not retrieve guild from cache for emoji update")
			}

			var fields []discord.EmbedField

			added, removed := emojiDiff(o, c.Emojis)

			if len(added) != 0 {
				fields = append(fields, genField(added, "Added emojis", false))
			}

			if len(removed) != 0 {
				fields = append(fields, genField(removed, "Removed emojis", true))
			}

			if len(fields) != 0 {
				handleAuditError(s.SendEmbeds(auditChannel, discord.Embed{
					Description: "**:pencil: Server's emojis updated!**",
					Timestamp:   discord.NowTimestamp(),
					Color:       0xF1C40F,
					Fields:      fields,
				}))
			}
		})
	})
}

func emojiDiff(old, new []discord.Emoji) (added, removed []discord.Emoji) {
	// todo maybe optimize this somehow? or make it more obvious what it's doing
	oldMap := make(map[discord.EmojiID]discord.Emoji, len(old))
	newMap := make(map[discord.EmojiID]discord.Emoji, len(new))
	added = make([]discord.Emoji, 0)
	removed = make([]discord.Emoji, 0)
	for _, e := range old {
		oldMap[e.ID] = e
	}
	for _, e := range new {
		newMap[e.ID] = e
	}

	for id, e := range newMap {
		_, ok := oldMap[id]
		if !ok {
			added = append(added, e)
		}
	}
	for id, e := range oldMap {
		_, ok := newMap[id]
		if !ok {
			removed = append(removed, e)
		}
	}

	return
}
