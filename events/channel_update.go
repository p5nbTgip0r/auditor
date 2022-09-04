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
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(old, new discord.Channel) {
		log.Debug().Interface("channel", new).Msg("Channel update")

		e := &discord.Embed{
			Description: channelChangeHeader(updated, new),
			Timestamp:   discord.NowTimestamp(),
			Color:       color.Gold,
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Channel ID: %s", new.ID)},
		}

		if old.Name != new.Name {
			util.AddUpdatedField(e, "Name", old.Name, new.Name, true)
		}
		if old.NSFW != new.NSFW {
			util.AddUpdatedField(e, "NSFW", util.YesNoBool(old.NSFW), util.YesNoBool(new.NSFW), false)
		}
		if old.ParentID != new.ParentID {
			oldParent, err := s.Channel(old.ParentID)
			newParent, err2 := s.Channel(new.ParentID)
			if err == nil && err2 == nil {
				util.AddUpdatedField(e, "Category", oldParent.Name, newParent.Name, true)
			}
		}
		if old.Position != new.Position {
			// TODO in its current state, this spams anytime a channel is moved.
			// it seems like discord updates every channel's Position field when that happens.
			// TODO might be nice to format this:
			// After [channel X] & Before [channel Y]
			//util.AddUpdatedField(e, "Position", old.Position, new.Position, true)
		}

		switch old.Type {
		case discord.GuildText, discord.GuildNews:
			if old.Topic != new.Topic {
				util.AddUpdatedField(e, "Topic", old.Topic, new.Topic, false)
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
			if old.DefaultAutoArchiveDuration != new.DefaultAutoArchiveDuration {
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

				util.AddUpdatedField(e, "Thread Auto-Archive", fmtArchiveDuration(old.DefaultAutoArchiveDuration), fmtArchiveDuration(new.DefaultAutoArchiveDuration), false)
			}

			oldNews := old.Type == discord.GuildNews
			newNews := new.Type == discord.GuildNews
			if oldNews != newNews {
				util.AddUpdatedField(e, "News", util.YesNoBool(oldNews), util.YesNoBool(newNews), false)
			}

		case discord.GuildVoice, discord.GuildStageVoice:
			if old.VoiceBitrate != new.VoiceBitrate {
				fmtBitrate := func(s uint) string {
					return fmt.Sprintf("%d kbps", s/1000)
				}
				util.AddUpdatedField(e, "Voice Bitrate", fmtBitrate(old.VoiceBitrate), fmtBitrate(new.VoiceBitrate), true)
			}

			if old.VideoQualityMode != new.VideoQualityMode {
				fmtVidQuality := func(s discord.VideoQualityMode) string {
					if s == discord.AutoVideoQuality {
						return "Auto"
					} else if s == discord.FullVideoQuality {
						return "720p"
					} else {
						return "Unknown"
					}
				}
				util.AddUpdatedField(e, "Video Quality", fmtVidQuality(old.VideoQualityMode), fmtVidQuality(new.VideoQualityMode), false)
			}
			if old.VoiceUserLimit != new.VoiceUserLimit {
				fmtUsers := func(s uint) string {
					if s == 0 {
						return "unlimited"
					} else {
						return english.Plural(int(s), "user", "")
					}
				}

				util.AddUpdatedField(e, "User Limit", fmtUsers(old.VoiceUserLimit), fmtUsers(new.VoiceUserLimit), false)
			}
			if old.RTCRegionID != new.RTCRegionID {
				fmtRegion := func(s string) string {
					if s == "automatic" {
						return "Automatic"
					}

					region, ok := util.Regions[s]
					if ok {
						return region.Emoji + " " + region.Name
					}
					return "Unknown"
				}
				util.AddUpdatedField(e, "Region Override", fmtRegion(old.RTCRegionID), fmtRegion(new.RTCRegionID), false)
			}
		}

		if len(e.Fields) == 0 {
			log.Debug().Interface("channelId", new.ID).Msg("Ignoring channel update as no known fields were changed")
			return
		}

		handleAuditError(s.SendEmbeds(auditChannel, *e))
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.ChannelUpdateEvent) {
			if !audit.AuditChannelUpdate.Check(&c.GuildID, &c.ID) {
				return
			}

			old, err := s.ChannelStore.Channel(c.ID)
			if err != nil {
				go handleError(
					audit.AuditChannelUpdate,
					err,
					fmt.Sprintf("Could not retrieve channel from cache: `%s` / `%s`", c.Channel.Name, c.Channel.ID),
					nil,
				)
				return
			}

			go handle(*old, c.Channel)
		})
	})
}
