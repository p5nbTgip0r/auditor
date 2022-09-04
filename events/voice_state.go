package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(old *discord.VoiceState, new discord.VoiceState) {
		user, _ := s.User(new.UserID)
		var e *discord.Embed
		if user != nil {
			e = userBaseEmbed(*user, "", false)
		} else {
			e = discord.NewEmbed()
		}
		send := func(embed discord.Embed) {
			handleAuditError(s.SendEmbeds(auditChannel, embed))
		}

		if audit.AuditVoiceConnection.Check(&new.GuildID, &new.ChannelID) {
			// joined voice
			if (old == nil || !old.ChannelID.IsValid()) && new.ChannelID.IsValid() {
				e.Description = "**:inbox_tray: " + new.UserID.Mention() + " joined voice in " + new.ChannelID.Mention() + "**"
				e.Color = color.Green
				send(*e)
				// none of these other events should take place
				return
			}

			// protect against nil pointer access on `old`
			if old == nil {
				log.Error().Interface("voiceState", new).Msg("Old voice state is nil but new channel is not valid. This shouldn't happen.")
				return
			}

			// left voice
			if old.ChannelID.IsValid() && !new.ChannelID.IsValid() {
				e.Description = "**:outbox_tray: " + new.UserID.Mention() + " left voice from " + old.ChannelID.Mention() + "**"
				e.Color = color.Red
				send(*e)
				// none of these other events should take place
				return
			}

			// switched channels
			if old.ChannelID.IsValid() && new.ChannelID.IsValid() && old.ChannelID != new.ChannelID {
				e.Description = fmt.Sprintf("**:twisted_rightwards_arrows: %s switched voice channels: %s -> %s**", new.UserID.Mention(), old.ChannelID.Mention(), new.ChannelID.Mention())
				e.Color = color.Blue
				send(*e)
			}
		} else if old == nil {
			return
		}

		if !audit.AuditVoiceAudioState.Check(&new.GuildID, &new.ChannelID) {
			return
		}

		audioStatus := func(id discord.UserID, actionerID discord.UserID, status bool, text string) {
			var emoji string
			if status {
				emoji = ":loud_sound:"
				e.Color = color.Green
				text = "un" + text
			} else {
				emoji = ":mute:"
				e.Color = color.Red
			}
			if actionerID.IsValid() {
				text += " by " + actionerID.Mention()
			}

			e.Description = fmt.Sprintf("**%s %s was %s**", emoji, id.Mention(), text)
			send(*e)
		}

		filter := func(key discord.AuditLogChangeKey, wantedValue bool) func(entry discord.AuditLogEntry) bool {
			return func(entry discord.AuditLogEntry) bool {
				for _, change := range entry.Changes {
					if change.Key == key {
						var b bool
						_ = change.NewValue.UnmarshalTo(&b)
						return b == wantedValue
					}
				}
				return false
			}
		}

		if old.Mute != new.Mute {
			_, actioner, _, _ := util.GetAuditActioner(s, new.GuildID, discord.Snowflake(new.UserID),
				api.AuditLogData{ActionType: discord.MemberUpdate},
				filter(discord.AuditUserMute, new.Mute),
			)

			audioStatus(new.UserID, actioner, !new.Mute, "muted")
		}

		if old.Deaf != new.Deaf {
			_, actioner, _, _ := util.GetAuditActioner(s, new.GuildID, discord.Snowflake(new.UserID),
				api.AuditLogData{ActionType: discord.MemberUpdate},
				filter(discord.AuditUserDeaf, new.Deaf),
			)

			audioStatus(new.UserID, actioner, !new.Deaf, "deafened")
		}
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.VoiceStateUpdateEvent) {
			if !audit.AuditVoiceConnection.Check(&c.GuildID, &c.ChannelID) || !audit.AuditVoiceAudioState.Check(&c.GuildID, &c.ChannelID) {
				return
			}
			log.Debug().Interface("event", c).Msg("Received updated voice state")
			old, _ := s.VoiceStateStore.VoiceState(c.GuildID, c.UserID)

			go handle(old, c.VoiceState)
		})
	})
}
