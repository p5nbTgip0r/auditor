package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
)

func init() {
	handle := func(old discord.Channel) {
		e := &discord.Embed{
			Description: channelChangeHeader(deleted, old),
			Timestamp:   discord.NowTimestamp(),
			Color:       color.Greyple,
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Channel ID: %s", old.ID)},
		}

		util.AddField(e, "Name", old.Name, true)
		util.AddField(e, "NSFW", util.YesNoBool(old.NSFW), false)
		oldParent, err := s.Channel(old.ParentID)
		if err == nil {
			if !oldParent.ParentID.IsValid() {
				// parent is a category (we're a channel)
				util.AddField(e, "Category", oldParent.Name, true)
			} else {
				// parent is a channel (we're a thread)
				util.AddField(e, "Thread Parent", oldParent.Name, true)
			}
		}

		switch old.Type {
		case discord.GuildText, discord.GuildNews:
			util.AddField(e, "Topic", old.Topic, false)
			if old.UserRateLimit != discord.NullSecond {
				fmtSeconds := func(s discord.Seconds) string {
					if s == 0 {
						return "off"
					} else {
						return humanize.Duration("", s.Duration())
					}
				}(old.UserRateLimit)

				util.AddField(e, "Slowmode", fmtSeconds, false)
			}

			fmtArchiveDuration := func() string {
				switch old.DefaultAutoArchiveDuration {
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
				return old.DefaultAutoArchiveDuration.String()
			}()
			util.AddField(e, "Thread Auto-Archive", fmtArchiveDuration, false)
			util.AddField(e, "News", util.YesNoBool(old.Type == discord.GuildNews), false)

		case discord.GuildVoice, discord.GuildStageVoice:
			util.AddField(e, "Voice Bitrate", fmt.Sprintf("%d kbps", old.VoiceBitrate/1000), true)

			fmtVidQuality := func(s discord.VideoQualityMode) string {
				if s == discord.AutoVideoQuality {
					return "Auto"
				} else if s == discord.FullVideoQuality {
					return "720p"
				} else {
					return "Unknown"
				}
			}(old.VideoQualityMode)
			util.AddField(e, "Video Quality", fmtVidQuality, false)

			fmtUsers := func(s uint) string {
				if s == 0 {
					return "unlimited"
				} else {
					return english.Plural(int(s), "user", "")
				}
			}(old.VoiceUserLimit)
			util.AddField(e, "User Limit", fmtUsers, false)

			fmtRegion := func(s string) string {
				if s == "automatic" {
					return "Automatic"
				}

				region, ok := util.Regions[s]
				if ok {
					return region.Emoji + " " + region.Name
				}
				return "Unknown"
			}(old.RTCRegionID)
			util.AddField(e, "Region Override", fmtRegion, false)
		}

		handleAuditError(s.SendEmbeds(auditChannel, *e))
	}

	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.ChannelDeleteEvent) {
			if !check(audit.AuditChannelDelete, &c.GuildID, &c.ID) {
				return
			}

			go handle(c.Channel)
		})
	})
}
