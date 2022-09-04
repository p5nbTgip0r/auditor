package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.ThreadCreateEvent) {
			if !audit.AuditChannelCreate.Check(&c.GuildID, &c.ID) {
				return
			}

			e := discord.NewEmbed()
			user, _ := s.User(c.OwnerID)
			if user != nil {
				e = userBaseEmbed(*user, "", false)
			}
			e.Description = channelChangeHeader(created, c.Channel)
			e.Color = color.Green
			e.Timestamp = discord.NewTimestamp(c.CreatedAt())
			e.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Thread ID: %s | Parent ID: %s", c.ID, c.ParentID)}

			util.AddField(e, "Parent", c.ParentID.Mention(), false)

			if user != nil {
				util.AddField(e, "Author", util.UserTag(*user), false)
			} else {
				util.AddField(e, "Author", c.OwnerID.Mention(), false)
			}

			if c.Type == discord.GuildPrivateThread {
				util.AddField(e, "Private", "Yes", false)
			}

			// realistically, none of the fields should be filled at this point
			if c.Topic != "" {
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Topic", Value: c.Topic})
			}

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
