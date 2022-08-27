package events

import (
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

//go:generate stringer -type=changeType -linecomment
type changeType uint

const (
	created changeType = iota
	deleted
	updated
	permissionsUpdated // permissions updated
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddHandler(func(c *gateway.ChannelCreateEvent) {
			if !AuditChannelCreate.check(&c.GuildID, &c.ID) {
				return
			}

			e := &discord.Embed{
				Color:     color.Green,
				Timestamp: discord.NewTimestamp(c.CreatedAt()),
				Footer:    &discord.EmbedFooter{Text: fmt.Sprintf("ID: %s", c.ID)},
			}

			e.Description = channelChangeHeader(created, c.Channel)

			// realistically, none of the fields should be filled at this point
			if c.Channel.Type == discord.GuildVoice {
				if region, ok := util.Regions[c.RTCRegionID]; ok {
					e.Fields = append(e.Fields, discord.EmbedField{Name: "Voice Region", Value: region.Emoji + " " + region.Name})
				}
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Voice Bitrate", Value: fmt.Sprintf("%d", c.VoiceBitrate)})
				if c.VoiceUserLimit != 0 {
					e.Fields = append(e.Fields, discord.EmbedField{Name: "Voice User Limit", Value: fmt.Sprintf("%d", c.VoiceUserLimit)})
				}
			}

			if c.Topic != "" {
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Topic", Value: c.Topic})
			}
			if c.UserRateLimit > 0 {
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Slowmode", Value: c.UserRateLimit.String()})
			}

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}

func channelChangeHeader(t changeType, c discord.Channel) string {
	emoji := ":pencil2:"
	if t == created {
		emoji = ":pencil:"
	} else if t == deleted {
		emoji = ":wastebasket:"
	} else if t == permissionsUpdated {
		emoji = ":crossed_swords:"
	}

	var chanType string
	var mention string
	switch c.Type {
	case discord.GuildText:
		chanType = "Text channel"
	case discord.GuildNews:
		chanType = "News channel"
	case discord.GuildVoice:
		chanType = "Voice channel"
	case discord.GuildCategory:
		chanType = "Channel category"
		mention = fmt.Sprintf("`%s`", c.Name)
	default:
		chanType = "Channel"
		mention = fmt.Sprintf("`%s`", c.Name)
	}
	if mention == "" {
		mention = fmt.Sprintf("`#%s`", c.Name)
	}

	return fmt.Sprintf("**%s %s %s: %s**", emoji, chanType, t.String(), mention)
}
