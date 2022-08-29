package events

import (
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(old, new discord.Channel) {
		log.Debug().Interface("thread", new).Msg("Thread update")

		e := &discord.Embed{
			Description: channelChangeHeader(updated, new),
			Timestamp:   discord.NowTimestamp(),
			Color:       color.Gold,
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Thread ID: %s | Parent ID: %s", new.ID, new.ParentID)},
		}

		if old.Name != new.Name {
			util.AddUpdatedField(e, "Name", old.Name, new.Name, true)
		}

		if old.ThreadMetadata.Archived != new.ThreadMetadata.Archived {
			if new.ThreadMetadata.Archived {
				e.Description = channelChangeHeader(archived, new)
				util.AddField(e, "Archived on", util.Timestamp(new.ThreadMetadata.ArchiveTimestamp.Time(), util.Relative), false)
			} else {
				util.AddUpdatedField(e, "Archived", util.YesNoBool(true), util.YesNoBool(false), false)
			}
		}

		if old.ThreadMetadata.Locked != new.ThreadMetadata.Locked {
			util.AddUpdatedField(e, "Locked", util.YesNoBool(old.ThreadMetadata.Locked), util.YesNoBool(new.ThreadMetadata.Locked), false)
		}

		if old.ThreadMetadata.Invitable != new.ThreadMetadata.Invitable {
			util.AddUpdatedField(e, "Invitable", util.YesNoBool(old.ThreadMetadata.Invitable), util.YesNoBool(new.ThreadMetadata.Invitable), false)
		}

		if old.ThreadMetadata.AutoArchiveDuration != new.ThreadMetadata.AutoArchiveDuration {
			fmtArchiveDuration := func(s discord.ArchiveDuration) string {
				switch s {
				case discord.OneHourArchive:
					return "1 Hour"
				case 0, discord.OneDayArchive:
					// default value
					return "24 Hours"
				case discord.ThreeDaysArchive:
					return "3 Days"
				case discord.SevenDaysArchive:
					return "1 Week"
				}
				return s.String()
			}

			util.AddUpdatedField(e, "Auto-Archive Duration", fmtArchiveDuration(old.ThreadMetadata.AutoArchiveDuration), fmtArchiveDuration(new.ThreadMetadata.AutoArchiveDuration), false)
		}

		if old.UserRateLimit != new.UserRateLimit {
			fmtSeconds := func(s discord.Seconds) string {
				if s == 0 {
					return "off"
				} else {
					return humanize.Duration("", s.Duration())
				}
			}

			util.AddUpdatedField(e, "Slowmode", fmtSeconds(old.UserRateLimit), fmtSeconds(new.UserRateLimit), false)
		}

		if len(e.Fields) == 0 {
			log.Debug().Interface("threadId", new.ID).Msg("Ignoring thread update as no known fields were changed")
			return
		}

		handleAuditError(s.SendEmbeds(auditChannel, *e))
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.ThreadUpdateEvent) {
			if !AuditChannelUpdate.check(&c.GuildID, &c.ID) {
				return
			}

			old, err := s.ChannelStore.Channel(c.ID)
			if err != nil {
				go handleError(
					AuditChannelUpdate,
					err,
					fmt.Sprintf("Could not retrieve thread from cache: `%s` / `%s`", c.Channel.Name, c.Channel.ID),
					nil,
				)
				return
			}

			go handle(*old, c.Channel)
		})
	})
}
