package events

import (
	"audit/util"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
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

	entry, _, actionerMem, _ := getAuditActioner(
		c.GuildID,
		discord.Snowflake(c.User.ID),
		api.AuditLogData{},
		// only consider bans and kicks
		mapset.NewSet(discord.MemberBanAdd, discord.MemberKick),
	)

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
	if entry != nil && actionerMem != nil {
		switch entry.ActionType {
		case discord.MemberBanAdd:
			embed.Color = 0x992D22
			embed.Description = fmt.Sprintf("**:rotating_light: %s was banned**", c.User.Mention())
			embed.Fields = append(embed.Fields,
				discord.EmbedField{
					Name:  "Banned by",
					Value: util.FullTag(actionerMem.User),
				},
			)
		case discord.MemberKick:
			embed.Color = 0xA84300
			embed.Description = fmt.Sprintf("**:boot: %s was kicked**", c.User.Mention())
			embed.Fields = append(embed.Fields,
				discord.EmbedField{
					Name:  "Kicked by",
					Value: util.FullTag(actionerMem.User),
				},
			)
		}
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

func getAuditActioner(g discord.GuildID, f discord.Snowflake, data api.AuditLogData, eventTypes mapset.Set[discord.AuditLogEvent]) (
	entry *discord.AuditLogEntry,
	actionerID *discord.UserID,
	actionerMember *discord.Member,
	err error,
) {
	audit, err := s.AuditLog(g, data)

	if err != nil {
		return
	}

	for _, e := range audit.Entries {
		// check if an event type filter was specified
		if eventTypes.Cardinality() != 0 && !eventTypes.Contains(e.ActionType) {
			continue
		}
		// only consider entries within 30 seconds of now
		if e.CreatedAt().Sub(time.Now()).Abs().Seconds() > 30 {
			continue
		}
		if e.TargetID == f {
			entry = &e
			actionerID = &e.UserID
			actionerMember, _ = s.Member(g, e.UserID)
			break
		}
	}

	return
}
