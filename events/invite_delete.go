package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.InviteDeleteEvent) {
			if !check(audit.InviteDelete, &c.GuildID, nil) {
				return
			}

			e := &discord.Embed{
				Description: "**:wastebasket: An invite has been deleted**",
				Timestamp:   discord.NowTimestamp(),
				Color:       color.Red,
			}
			util.AddField(e, "Code", fmt.Sprintf("[**%s**](https://discord.gg/%s)", c.Code, c.Code), false)
			util.AddField(e, "Channel", util.ChannelTag(c.ChannelID), false)

			bot.QueueEmbed(audit.InviteDelete, c.GuildID, *e)
		})
	})
}
