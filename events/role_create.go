package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strings"
)

func init() {
	handler = append(handler, func() {
		s.AddHandler(func(c *gateway.GuildRoleCreateEvent) {
			if !check(audit.AuditRoleCreate, &c.GuildID, nil) {
				return
			}

			e := &discord.Embed{
				Description: fmt.Sprintf("**:crossed_swords: Role created: %s**", c.Role.Name),
				Color:       color.Green,
				Timestamp:   discord.NewTimestamp(c.Role.CreatedAt()),
				Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("ID: %s", c.Role.ID)},
			}

			e.Fields = append(e.Fields,
				discord.EmbedField{
					Name:  "Color",
					Value: color.ColorViewerLink(c.Role.Color, c.Role.Color.String()),
				},
			)

			if c.Role.Icon != "" {
				// todo archive icon
				link := fmt.Sprintf("[__Link__](%s)", c.Role.IconURL())
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Icon", Value: link})
			} else if c.Role.UnicodeEmoji != "" {
				e.Fields = append(e.Fields, discord.EmbedField{Name: "Emoji", Value: c.Role.UnicodeEmoji})
			}
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Hoisted", Value: util.YesNoBool(c.Role.Hoist)})
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Mentionable", Value: util.YesNoBool(c.Role.Mentionable)})

			perms := util.PermissionString(c.Role.Permissions)
			perms = strings.ReplaceAll(perms, "Administrator", "**Administrator**")
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Permissions", Value: perms})

			handleAuditError(s.SendEmbeds(auditChannel, *e))
		})
	})
}
