package events

import (
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
)

func init() {
	invites := make(map[discord.GuildID]map[string]int)

	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.GuildCreateEvent) {
			if !AuditMemberJoin.check(&c.ID, nil) {
				return
			}

			inv, err := s.GuildInvites(c.ID)
			if err != nil {
				log.Warn().
					Err(err).
					Interface("guild_id", c.ID).
					Msg("Could not retrieve invites for guild")
				return
			}

			tempInvites := make(map[string]int, len(inv))
			for _, invite := range inv {
				tempInvites[invite.Code] = invite.Uses
			}

			invites[c.ID] = tempInvites
		})
		s.AddHandler(func(c *gateway.GuildDeleteEvent) {
			if !AuditMemberJoin.check(&c.ID, nil) {
				return
			}

			log.Debug().
				Interface("guild_id", c.ID).
				Msg("Removing cached invites for guild")

			delete(invites, c.ID)
		})

		s.AddHandler(func(c *gateway.InviteCreateEvent) {
			if !AuditMemberJoin.check(&c.GuildID, nil) {
				return
			}

			i := invites[c.GuildID]
			i[c.Code] = 0
			invites[c.GuildID] = i
		})
		s.AddHandler(func(c *gateway.InviteDeleteEvent) {
			if !AuditMemberJoin.check(&c.GuildID, nil) {
				return
			}

			delete(invites[c.GuildID], c.Code)
		})

		s.AddHandler(func(c *gateway.GuildMemberAddEvent) {
			if !AuditMemberJoin.check(&c.GuildID, nil) {
				return
			}

			oldInvites, ok := invites[c.GuildID]
			var usedInvite *discord.Invite
			if ok {
				newInvites, err := s.GuildInvites(c.GuildID)
				if err != nil {
					log.Warn().
						Err(err).
						Interface("guild_id", c.GuildID).
						Interface("user_id", c.User.ID).
						Msg("Could not retrieve invites for guild on user join")
				} else {
					for _, newInvite := range newInvites {
						if newInvite.Uses > oldInvites[newInvite.Code] {
							usedInvite = &newInvite
							break
						}
					}
				}
			}

			embed := userBaseEmbed(c.User, "", true)
			embed.Color = color.Green

			embed.Description = fmt.Sprintf("**:inbox_tray: %s joined the server**", c.User.Mention())
			embed.Fields = append(embed.Fields,
				discord.EmbedField{
					Name:  "Account Creation",
					Value: util.Timestamp(c.User.CreatedAt(), util.Relative),
				},
			)
			if usedInvite != nil {
				invMsg := fmt.Sprintf(
					"Inviter: %s\nCode: `%s`\nUses: `%d`",
					util.UserTag(*usedInvite.Inviter),
					usedInvite.Code,
					usedInvite.Uses,
				)

				embed.Fields = append(embed.Fields,
					discord.EmbedField{
						Name:  "Used Invite",
						Value: invMsg,
					},
				)
			}
			if c.IsPending {
				embed.Description += "\n\n**:clipboard: User is currently in membership screening**"
			}
			handleAuditError(s.SendMessage(auditChannel, "", *embed))
		})
	})
}
