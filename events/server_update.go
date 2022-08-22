package events

import (
	"audit/util"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.GuildUpdateEvent) {
			if !AuditServerEdited.check(&c.ID, nil) {
				return
			}
			g, err := s.GuildStore.Guild(c.ID)
			if err != nil {
				// todo: proper warning for this
				log.Warn().Err(err).Interface("event", c).Msg("Could not retrieve guild from cache for server update")
				return
			}

			go handleServerUpdate(*g, c.Guild)
		})
	})
}

func handleServerUpdate(old, new discord.Guild) {
	var embed discord.Embed
	embed.Description = "**:pencil: Server information updated!**"
	embed.Color = 0xF1C40F
	embed.Timestamp = discord.NowTimestamp()
	var fields []discord.EmbedField

	if old.Name != new.Name {
		fields = append(fields, discord.EmbedField{Name: "Name", Value: fmt.Sprintf("`%s` -> `%s`", old.Name, new.Name)})
	}

	if old.AFKTimeout != new.AFKTimeout {
		oldMins := old.AFKTimeout / 60
		oldAfkText := fmt.Sprintf("%d minute%s", oldMins, util.Plural(int(oldMins)))
		newMins := new.AFKTimeout / 60
		newAfkText := fmt.Sprintf("%d minute%s", newMins, util.Plural(int(newMins)))

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
		mfaReq := "No"
		if new.MFA == discord.ElevatedMFA {
			mfaReq = "Yes"
		}

		fields = append(fields, discord.EmbedField{
			Name:  "Admin 2FA required",
			Value: mfaReq,
		})
	}

	if len(fields) == 0 {
		return
	}

	embed.Fields = fields

	handleAuditError(s.SendEmbeds(auditChannel, embed))
}
