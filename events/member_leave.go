package events

import (
	"audit/util"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"time"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildMemberRemoveEvent) {
			m, _ := s.MemberStore.Member(c.GuildID, c.User.ID)

			go handleMemberRemove(c, m)
		})
	})
}

func handleMemberRemove(c *gateway.GuildMemberRemoveEvent, m *discord.Member) {
	// sleep for a few seconds to give the audit log a chance to catch up
	time.Sleep(time.Second * 3)
	_, bannerMember, _ := getAuditActioner(c.GuildID, discord.Snowflake(c.User.ID), api.AuditLogData{ActionType: discord.MemberKick})

	var joined string
	if m != nil {
		joined = fmt.Sprintf("<t:%d:R>", m.Joined.Time().Unix())
	} else {
		joined = "Could not retrieve"
	}

	embed := userBaseEmbed(c.User, "", true)
	embed.Color = 0xE74C3C

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:  "Joined server",
			Value: joined,
		},
	)
	if bannerMember != nil {
		embed.Description = fmt.Sprintf("**:boot: %s was kicked**", c.User.Mention())
		embed.Fields = append(embed.Fields,
			discord.EmbedField{
				Name:  "Kicked by",
				Value: util.FullTag(bannerMember.User),
			},
		)
	} else {
		embed.Description = fmt.Sprintf("**:outbox_tray: %s left the server**", c.User.Mention())
		embed.Fields = append(embed.Fields,
			discord.EmbedField{
				Name:  "Account creation",
				Value: fmt.Sprintf("<t:%d:R>", c.User.CreatedAt().Unix()),
			},
		)
	}

	handleAuditError(s.SendMessage(auditChannel, "", *embed))
}
