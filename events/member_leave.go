package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"time"
)

func init() {
	handle := func(c *gateway.GuildMemberRemoveEvent, m *discord.Member) {
		// sleep for a few seconds to give the audit log a chance to catch up
		time.Sleep(time.Second * 3)
		// TODO don't call `check` several times
		// instead, make a new `check` function which takes in a slice/vararg of audittypes and compare the disabled slice to the passed-in slice
		cont := false
		for _, t := range []audit.AuditType{audit.AuditMemberLeave, audit.AuditMemberKick, audit.AuditMemberBan} {
			if check(t, &c.GuildID, nil) {
				cont = true
				break
			}
		}
		if !cont {
			return
		}

		entry, _, actionerMem, _ := getAuditActioner(
			c.GuildID,
			discord.Snowflake(c.User.ID),
			api.AuditLogData{},
			// only consider bans and kicks
			mapset.NewSet(discord.MemberBanAdd, discord.MemberKick),
		)

		var joined string
		if m != nil {
			joined = util.Timestamp(m.Joined.Time(), util.Relative)
		} else {
			joined = "Could not retrieve"
		}

		embed := userBaseEmbed(c.User, "", true)
		embed.Color = color.Red

		embed.Fields = append(embed.Fields,
			discord.EmbedField{
				Name:  "Joined server",
				Value: joined,
			},
		)

		if entry != nil && actionerMem != nil {
			switch entry.ActionType {
			case discord.MemberBanAdd:
				if !check(audit.AuditMemberBan, &c.GuildID, nil) {
					return
				}

				embed.Color = color.DarkRed
				embed.Description = fmt.Sprintf("**:rotating_light: %s was banned**", c.User.Mention())
				embed.Fields = append(embed.Fields,
					discord.EmbedField{
						Name:  "Banned by",
						Value: util.UserTag(actionerMem.User),
					},
				)
			case discord.MemberKick:
				if !check(audit.AuditMemberKick, &c.GuildID, nil) {
					return
				}

				embed.Color = color.DarkOrange
				embed.Description = fmt.Sprintf("**:boot: %s was kicked**", c.User.Mention())
				embed.Fields = append(embed.Fields,
					discord.EmbedField{
						Name:  "Kicked by",
						Value: util.UserTag(actionerMem.User),
					},
				)
			}
		} else {
			if !check(audit.AuditMemberLeave, &c.GuildID, nil) {
				return
			}

			embed.Description = fmt.Sprintf("**:outbox_tray: %s left the server**", c.User.Mention())
			embed.Fields = append(embed.Fields,
				discord.EmbedField{
					Name:  "Account creation",
					Value: util.Timestamp(c.User.CreatedAt(), util.Relative),
				},
			)
		}

		handleAuditError(s.SendMessage(auditChannel, "", *embed))
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildMemberRemoveEvent) {
			m, err := s.MemberStore.Member(c.GuildID, c.User.ID)
			if err != nil {
				go handleError(
					audit.AuditMemberLeave,
					err,
					fmt.Sprintf("Could not retrieve member from cache: %s", util.UserTag(c.User)),
					&c.User,
				)
				return
			}

			go handle(c, m)
		})
	})
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
