package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util"
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

	if (diff.fields.Has(fieldUserName) || diff.fields.Has(fieldUserDiscriminator)) && check(audit.UserName, &new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:pencil: %s Discord tag changed**", new.User.Mention()))
		util.AddField(c, "Old tag", old.Tag(), false)
		util.AddField(c, "New tag", new.User.Tag(), false)
		bot.QueueEmbed(audit.UserName, new.GuildID, *c)
	}

	if diff.fields.Has(fieldUserAvatar) && check(audit.UserAvatar, &new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:frame_photo: %s changed their __user__ avatar**", new.User.Mention()))
		util.AddField(c, "Old avatar", old.AvatarURL(), false)
		util.AddField(c, "New avatar", new.User.AvatarURL(), false)
		bot.QueueEmbed(audit.UserAvatar, new.GuildID, *c)
	}
}
