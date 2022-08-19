package events

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddHandler(func(c *gateway.GuildBanRemoveEvent) {
			e := userBaseEmbed(c.User, "", true)
			e.Color = 0x1F8B4C

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
