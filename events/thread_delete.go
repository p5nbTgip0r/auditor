package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(old discord.Channel) {
		if !check(audit.ChannelDelete, &old.GuildID, &old.ID) {
			return
		}

		log.Debug().Interface("thread", old).Msg("Thread deleted")

		e := &discord.Embed{
			Description: channelChangeHeader(deleted, old),
			Timestamp:   discord.NowTimestamp(),
			Color:       color.Greyple,
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Thread ID: %s | Parent ID: %s", old.ID, old.ParentID)},
		}

		util.AddField(e, "Name", old.Name, true)
		user, _ := s.User(old.OwnerID)
		if user != nil {
			util.AddField(e, "Author", util.UserTag(*user), false)
		} else {
			util.AddField(e, "Author", old.OwnerID.Mention(), false)
		}

		if old.ThreadMetadata.Archived {
			util.AddField(e, "Archived on", util.Timestamp(old.ThreadMetadata.ArchiveTimestamp.Time(), util.Relative), false)
		} else {
			util.AddField(e, "Archived", util.YesNoBool(false), false)
		}

		util.AddField(e, "Locked", util.YesNoBool(old.ThreadMetadata.Locked), false)
		util.AddField(e, "Invitable", util.YesNoBool(old.ThreadMetadata.Invitable), false)

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
		util.AddField(e, "Auto-Archive Duration", fmtArchiveDuration(old.ThreadMetadata.AutoArchiveDuration), false)

		if old.UserRateLimit != 0 {
			util.AddField(e, "Slowmode", humanize.Duration("", old.UserRateLimit.Duration()), false)
		}

		bot.QueueEmbed(audit.ChannelDelete, old.GuildID, *e)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.ThreadDeleteEvent) {
			old, err := s.ChannelStore.Channel(c.ID)
			if err != nil {
				go handleError(
					audit.ChannelDelete,
					c.GuildID,
					err,
					fmt.Sprintf("Could not retrieve thread from cache: `%s` / `%s`", c.ParentID.Mention(), c.ID),
					discord.User{},
				)
				return
			}

			go handle(*old)
		})
	})
}
