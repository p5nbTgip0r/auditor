package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strings"
)

func init() {
	handle := func(c gateway.GuildRoleDeleteEvent, role discord.Role) {
		if !check(audit.RoleDelete, &c.GuildID, nil) {
			return
		}

		e := &discord.Embed{
			Description: fmt.Sprintf("**:wastebasket: Role deleted: %s**", role.Name),
			Color:       color.Red,
			Timestamp:   discord.NowTimestamp(),
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("ID: %s", role.ID)},
		}

		e.Fields = append(e.Fields,
			discord.EmbedField{
				Name:  "Color",
				Value: color.ColorViewerLink(role.Color, role.Color.String()),
			},
		)

		if role.Icon != "" {
			// todo archive icon
			link := fmt.Sprintf("[__Link__](%s)", role.IconURL())
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Icon", Value: link})
		} else if role.UnicodeEmoji != "" {
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Emoji", Value: role.UnicodeEmoji})
		}
		e.Fields = append(e.Fields, discord.EmbedField{Name: "Hoisted", Value: util.YesNoBool(role.Hoist)})
		e.Fields = append(e.Fields, discord.EmbedField{Name: "Mentionable", Value: util.YesNoBool(role.Mentionable)})

		perms := util.PermissionString(role.Permissions)
		perms = strings.ReplaceAll(perms, "Administrator", "**Administrator**")
		e.Fields = append(e.Fields, discord.EmbedField{Name: "Permissions", Value: perms})

		bot.QueueEmbed(audit.RoleDelete, c.GuildID, *e)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildRoleDeleteEvent) {
			role, err := s.RoleStore.Role(c.GuildID, c.RoleID)
			if err != nil {
				go handleError(
					audit.RoleDelete,
					c.GuildID,
					err,
					fmt.Sprintf("Could not retrieve role from cache: `%s`", c.RoleID),
					discord.User{},
				)
				return
			}

			go handle(*c, *role)
		})
	})
}
