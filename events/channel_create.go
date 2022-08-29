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
				parent, err := s.Channel(c.ParentID)
				if err == nil {
					util.AddField(e, "Category", parent.Name, false)
				}
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
	emoji := ":pencil2:"
	if t == created {
		emoji = ":pencil:"
	} else if t == deleted {
		emoji = ":wastebasket:"
	} else if t == permissionsUpdated {
		emoji = ":crossed_swords:"
	}

	var chanType string
	var mention string
	switch c.Type {
	case discord.GuildText:
		chanType = "Text channel"
	case discord.GuildNews:
		chanType = "News channel"
	case discord.GuildVoice:
		chanType = "Voice channel"
	case discord.GuildCategory:
		chanType = "Category"
		mention = fmt.Sprintf("`%s`", c.Name)
	default:
		chanType = "Channel"
		mention = fmt.Sprintf("`%s`", c.Name)
	}
	if mention == "" {
		mention = fmt.Sprintf("`#%s`", c.Name)
	}

	return fmt.Sprintf("**%s %s %s: %s**", emoji, chanType, t.String(), mention)
}
