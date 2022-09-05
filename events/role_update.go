package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(old, new discord.Role, guildID discord.GuildID) {
		if !check(audit.RoleUpdate, &guildID, nil) {
			return
		}

		e := &discord.Embed{
			Description: fmt.Sprintf("**:pencil: Role updated: %s**", new.Name),
			Color:       color.Gold,
			Timestamp:   discord.NowTimestamp(),
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("ID: %s", new.ID)},
		}

		if old.Name != new.Name {
			e.Fields = append(e.Fields,
				discord.EmbedField{
					Name:  "Name",
					Value: fmt.Sprintf("`%s` -> `%s`", old.Name, new.Name),
				},
			)
		}
		if old.Color != new.Color {
			e.Fields = append(e.Fields,
				discord.EmbedField{
					Name:  "Color",
					Value: color.ColorViewerLink(old.Color, old.Color.String()) + " -> " + color.ColorViewerLink(new.Color, new.Color.String()),
				},
			)
		}

		if old.Hoist != new.Hoist {
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Hoisted", Value: util.YesNoBool(old.Hoist) + " -> " + util.YesNoBool(new.Hoist)})
		}

		if old.Mentionable != new.Mentionable {
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Mentionable", Value: util.YesNoBool(old.Mentionable) + " -> " + util.YesNoBool(new.Mentionable)})
		}

		if old.Icon != new.Icon {
			// todo archive icon
			format := "[__Link__](%s)"
			var oldIcon string
			if icon := old.IconURL(); icon != "" {
				oldIcon = fmt.Sprintf(format, icon)
			}
			var newIcon string
			if icon := new.IconURL(); icon != "" {
				newIcon = fmt.Sprintf(format, icon)
			}
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Icon", Value: oldIcon + " -> " + newIcon})
		} else if old.UnicodeEmoji != new.UnicodeEmoji {
			oldEmoji := old.UnicodeEmoji
			if oldEmoji == "" {
				oldEmoji = "none"
			}
			newEmoji := new.UnicodeEmoji
			if newEmoji == "" {
				newEmoji = "none"
			}
			e.Fields = append(e.Fields, discord.EmbedField{Name: "Emoji", Value: oldEmoji + " -> " + newEmoji})
		}

		added, removed := util.SliceDiff(util.GetPermissions(old.Permissions), util.GetPermissions(new.Permissions))
		if added.Cardinality() != 0 {
			e.Fields = append(e.Fields, discord.EmbedField{Name: "✓ Allowed permissions", Value: util.PermissionSetString(added)})
		}
		if removed.Cardinality() != 0 {
			e.Fields = append(e.Fields, discord.EmbedField{Name: "✘ Denied permissions", Value: util.PermissionSetString(removed)})
		}

		if len(e.Fields) == 0 {
			return
		}

		bot.QueueEmbed(audit.RoleUpdate, guildID, *e)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildRoleUpdateEvent) {
			role, err := s.RoleStore.Role(c.GuildID, c.Role.ID)
			if err != nil {
				go handleError(
					audit.RoleUpdate,
					c.GuildID,
					err,
					fmt.Sprintf("Could not retrieve role from cache: `%s` / `%s`", c.Role.Name, c.Role.ID),
					discord.User{},
				)
				return
			}

			log.Debug().Interface("old_role", role).Interface("new_role", c.Role).Msgf("Got role update for %s", c.Role.Name)
			go handle(*role, c.Role, c.GuildID)
		})
	})
}
