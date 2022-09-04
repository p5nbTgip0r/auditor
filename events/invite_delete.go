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
		s.AddHandler(func(c *gateway.InviteDeleteEvent) {
			if !audit.AuditInviteDelete.Check(&c.GuildID, nil) {
				return
			}

			e := &discord.Embed{
				Description: "**:wastebasket: An invite has been deleted**",
				Timestamp:   discord.NowTimestamp(),
				Color:       color.Red,
			}
			util.AddField(e, "Code", fmt.Sprintf("[**%s**](https://discord.gg/%s)", c.Code, c.Code), false)
			util.AddField(e, "Channel", util.ChannelTag(c.ChannelID), false)

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
