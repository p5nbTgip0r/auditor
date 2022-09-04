package events

import (
	"audit/audit"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func diffUser(new gateway.GuildMemberUpdateEvent, old discord.User, diff userDiff) {
	getEmbed := func(desc string) *discord.Embed {
		c := userBaseEmbed(new.User, "", false)
		// https://github.com/Rapptz/discord.py/blob/de941ababe9da898dd62d2b2a2d21aaecac6bd09/discord/colour.py#L295
		c.Color = color.Gold
		c.Description = desc
		return c
	}

	if (diff.fields.Has(fieldUserName) || diff.fields.Has(fieldUserDiscriminator)) && audit.AuditUserName.Check(&new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:pencil: %s Discord tag changed**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old tag", Value: old.Tag()},
			discord.EmbedField{Name: "New tag", Value: new.User.Tag()},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	if diff.fields.Has(fieldUserAvatar) && audit.AuditUserAvatar.Check(&new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:frame_photo: %s changed their __user__ avatar**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old avatar", Value: old.AvatarURL()},
			discord.EmbedField{Name: "New avatar", Value: new.User.AvatarURL()},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}
}
