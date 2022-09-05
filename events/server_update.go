package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/dustin/go-humanize/english"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildUpdateEvent) {
			g, err := s.GuildStore.Guild(c.ID)
			if err != nil {
				go handleError(
					audit.ServerEdited,
					c.ID,
					err,
					"Could not retrieve guild from cache: `"+c.ID.String()+"`",
					discord.User{},
				)
				return
			}

			go handleServerUpdate(*g, c.Guild)
		})
	})
}

func handleServerUpdate(old, new discord.Guild) {
	if !check(audit.ServerEdited, &new.ID, nil) {
		return
	}

	var embed discord.Embed
	embed.Description = "**:pencil: Server information updated!**"
	embed.Color = color.Gold
	embed.Timestamp = discord.NowTimestamp()
	var fields []discord.EmbedField

	if old.Name != new.Name {
		fields = append(fields, discord.EmbedField{Name: "Name", Value: fmt.Sprintf("`%s` -> `%s`", old.Name, new.Name)})
	}

	if old.AFKTimeout != new.AFKTimeout {
		oldAfkText := english.Plural(int(old.AFKTimeout/60), "minute", "")
		newAfkText := english.Plural(int(new.AFKTimeout/60), "minute", "")

		fields = append(fields, discord.EmbedField{
			Name:  "AFK Timeout",
			Value: fmt.Sprintf("%s -> %s", oldAfkText, newAfkText),
		})
	}

	if old.AFKChannelID != new.AFKChannelID {
		oldChan := old.AFKChannelID.Mention()
		if old.AFKChannelID == discord.NullChannelID {
			oldChan = "__none__"
		}
		newChan := new.AFKChannelID.Mention()
		if new.AFKChannelID == discord.NullChannelID {
			newChan = "__none__"
		}

		fields = append(fields, discord.EmbedField{
			Name:  "AFK Channel",
			Value: fmt.Sprintf("%s -> %s", oldChan, newChan),
		})
	}

	if old.SystemChannelID != new.SystemChannelID {
		oldChan := old.SystemChannelID.Mention()
		if old.SystemChannelID == discord.NullChannelID {
			oldChan = "__none__"
		}
		newChan := new.SystemChannelID.Mention()
		if new.SystemChannelID == discord.NullChannelID {
			newChan = "__none__"
		}

		fields = append(fields, discord.EmbedField{
			Name:  "System Messages Channel",
			Value: fmt.Sprintf("%s -> %s", oldChan, newChan),
		})
	}
	// todo maybe one for system channel flags?
	// todo icon, banner, splash

	if old.Verification != new.Verification {
		verificationName := func(filter discord.Verification) string {
			switch filter {
			case discord.VeryHighVerification:
				return "Very High"
			case discord.HighVerification:
				return "High"
			case discord.MediumVerification:
				return "Medium"
			case discord.LowVerification:
				return "Low"
			case discord.NoVerification:
				return "None"
			default:
				return "Unknown"
			}
		}

		oldVerif := verificationName(old.Verification)
		newVerif := verificationName(new.Verification)

		fields = append(fields, discord.EmbedField{
			Name:  "Verification level",
			Value: fmt.Sprintf("__%s__ -> __%s__", oldVerif, newVerif),
		})
	}

	if old.ExplicitFilter != new.ExplicitFilter {
		filterName := func(filter discord.ExplicitFilter) string {
			switch filter {
			case discord.AllMembers:
				return "All members"
			case discord.MembersWithoutRoles:
				return "Members without roles"
			case discord.NoContentFilter:
				return "No content filter"
			default:
				return "Unknown"
			}
		}

		oldFilter := filterName(old.ExplicitFilter)
		newFilter := filterName(new.ExplicitFilter)

		fields = append(fields, discord.EmbedField{
			Name:  "Explicit content filter",
			Value: fmt.Sprintf("__%s__ -> __%s__", oldFilter, newFilter),
		})
	}

	if old.MFA != new.MFA {
		fields = append(fields, discord.EmbedField{
			Name:  "Admin 2FA required",
			Value: util.YesNoBool(new.MFA == discord.ElevatedMFA),
		})
	}

	if len(fields) == 0 {
		return
	}

	embed.Fields = fields

	handleAuditError(s.SendEmbeds(auditChannel, embed))
}
