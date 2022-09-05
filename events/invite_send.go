package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize/english"
	"regexp"
	"strings"
)

func init() {
	inviteRegex, _ := regexp.Compile(`(?:https?://)?(?:www\.)?(?:discord\.(?:gg|io|me|li)|(?:discordapp|discord)\.com/invite)/[\w]+`)

	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.MessageCreateEvent) {
			if !c.GuildID.IsValid() || c.Author.Bot || !check(audit.InviteSend, &c.GuildID, &c.ChannelID) {
				return
			}

			invites := inviteRegex.FindAllString(c.Content, -1)
			for _, embed := range c.Embeds {
				invites = append(invites, inviteRegex.FindAllString(embed.Description, -1)...)
				for _, field := range embed.Fields {
					invites = append(invites, inviteRegex.FindAllString(field.Value, -1)...)
				}
			}

			if len(invites) == 0 {
				return
			}

			jump, _ := util.DiscordJumpLink(c.GuildID, c.ChannelID, c.ID)
			e := userBaseEmbed(c.Author, jump, false)
			e.Color = color.Red
			e.Footer = &discord.EmbedFooter{
				Text: fmt.Sprintf("Message ID: %s | User ID: %s", c.ID.String(), c.Author.ID.String()),
			}
			var desc strings.Builder
			desc.WriteString(fmt.Sprintf("**:envelope_with_arrow: %s sent %s in %s**\n\n```",
				c.Author.Mention(),
				english.PluralWord(len(invites), "an invite", "multiple invites"),
				c.ChannelID.Mention(),
			))
			for i, invite := range invites {
				if i != 0 {
					desc.WriteString("\n")
				}
				desc.WriteString("- ")
				desc.WriteString(invite)
			}
			desc.WriteString("```")
			e.Description = desc.String()

			bot.QueueEmbed(audit.InviteSend, c.GuildID, *e)
		})
	})
}
