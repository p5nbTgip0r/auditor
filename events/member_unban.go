package events

import (
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddHandler(func(c *gateway.GuildBanRemoveEvent) {
			if !AuditMemberUnban.check(&c.GuildID, nil) {
				return
			}

			e := userBaseEmbed(c.User, "", true)
			e.Color = color.DarkGreen

			e.Description = fmt.Sprintf("**:ballot_box_with_check: %s was unbanned**", c.User.Mention())
			e.Fields = append(e.Fields,
				discord.EmbedField{
					Name:  "Account Creation",
					Value: fmt.Sprintf("<t:%d:R>", c.User.CreatedAt().Unix()),
				},
			)

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
