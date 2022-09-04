package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize/english"
)

func init() {
	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.InviteCreateEvent) {
			if !audit.AuditInviteCreate.Check(&c.GuildID, nil) {
				return
			}

			var e *discord.Embed
			if c.Inviter != nil {
				e = userBaseEmbed(*c.Inviter, "", false)
				util.AddField(e, "Inviter", util.UserTag(*c.Inviter), false)
			} else {
				e = discord.NewEmbed()
			}

			e.Description = "**:envelope_with_arrow: An invite has been created**"
			e.Timestamp = c.CreatedAt
			e.Color = color.Green
			util.AddField(e, "Code", fmt.Sprintf("[**%s**](https://discord.gg/%s)", c.Code, c.Code), false)

			maxAge := "Never"
			if c.InviteMetadata.MaxAge != 0 {
				maxAge = util.Timestamp(c.CreatedAt.Time().Add(c.InviteMetadata.MaxAge.Duration()), util.Relative)
			}
			util.AddField(e, "Expiration", maxAge, false)

			maxUses := "No limit"
			if c.InviteMetadata.MaxUses != 0 {
				maxUses = english.Plural(c.InviteMetadata.MaxUses, "use", "")
			}
			util.AddField(e, "Max uses", maxUses, false)

			util.AddField(e, "Channel", util.ChannelTag(c.ChannelID), false)

			if c.InviteMetadata.Temporary {
				util.AddField(e, "Temporary membership", "Yes", false)
			}

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
