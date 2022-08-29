package events

import (
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

//go:generate stringer -type=changeType -linecomment
type changeType uint

const (
	created changeType = iota
	deleted
	updated
	permissionsUpdated // permissions updated
	archived
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddHandler(func(c *gateway.ChannelCreateEvent) {
			if !AuditChannelCreate.check(&c.GuildID, &c.ID) {
				return
			}

			e := &discord.Embed{
				Color:     color.Green,
				Timestamp: discord.NewTimestamp(c.CreatedAt()),
			}

			e.Description = channelChangeHeader(created, c.Channel)

			if c.Type == discord.GuildCategory {
				e.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Category ID: %s", c.ID)}
			} else {
				util.AddField(e, "Category", c.ParentID.Mention(), false)
				e.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Channel ID: %s | Category ID: %s", c.ID, c.ParentID)}
			}

			// realistically, none of the fields should be filled at this point
			if c.Topic != "" {
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Topic", Value: c.Topic})
			}

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}

func channelChangeHeader(t changeType, c discord.Channel) string {
	var emoji string
	switch t {
	case created:
		emoji = ":pencil:"
	case updated:
		emoji = ":pencil2:"
	case deleted:
		emoji = ":wastebasket:"
	case permissionsUpdated:
		emoji = ":crossed_swords:"
	case archived:
		emoji = ":package:"
	}

	var chanType string
	mention := c.Mention()
	switch c.Type {
	case discord.GuildText:
		chanType = "Text channel"
	case discord.GuildPublicThread:
		chanType = "Thread"
		mention += " (`#" + c.Name + "`)"
	case discord.GuildPrivateThread:
		chanType = "Private thread"
		mention += " (`#" + c.Name + "`)"
	case discord.GuildNews:
		chanType = "News channel"
	case discord.GuildNewsThread:
		chanType = "News thread"
		mention += " (`#" + c.Name + "`)"
	case discord.GuildVoice:
		chanType = "Voice channel"
	case discord.GuildCategory:
		chanType = "Category"
	default:
		chanType = "Channel"
	}

	return fmt.Sprintf("**%s %s %s: %s**", emoji, chanType, t.String(), mention)
}
