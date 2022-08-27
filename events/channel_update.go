package events

import (
	"audit/util/color"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.ChannelUpdateEvent) {
			if !AuditChannelUpdate.check(&c.GuildID, &c.ID) {
				return
			}

			_, err := s.ChannelStore.Channel(c.ID)
			if err != nil {
				log.Warn().Interface("channel", c).Msgf("Couldn't get cached channel for channel update")
				return
			}
			log.Debug().Interface("channel", c).Msg("Channel update")

			e := discord.Embed{
				Description: channelChangeHeader(updated, c.Channel),
				Timestamp:   discord.NowTimestamp(),
				Color:       color.Gold,
			}

			e.Description += "\n\n"

			//handleAuditError(s.SendEmbeds(auditChannel, e))

			// todo
			//perms := util.PermissionString(c.Role.Permissions)
			//perms = strings.ReplaceAll(perms, "Administrator", "**Administrator**")
			//e.Fields = append(e.Fields, discord.EmbedField{Name: "Permissions", Value: perms})

			//handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
