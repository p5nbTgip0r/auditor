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
	"golang.org/x/exp/slices"
	"strings"
)

type overwritePair struct {
	old discord.Overwrite
	new discord.Overwrite
}

func init() {
	handle := func(old, new discord.Channel) {
		if !check(audit.ChannelUpdate, &new.GuildID, &new.ID) {
			return
		}
		msg := generateOverwriteMessage(old, new)
		if msg == "" {
			log.Debug().Msgf("Channel %s didn't have any permission updates", new.ID)
			return
		}
		// TODO somehow check if this permission update was a part of the parent category updating
		// this isn't very easy to do because the update for the parent category happens before all the child channels
		// and by the time that we compare permissions to the parent, the old permissions will have been updated to the new ones.
		// this basically means we have to get a list of all child channels on category update, then selectively
		// exclude permission updates from them for a short period of time (until the permission update events pass).
		sync, _ := isPermissionSync(new)
		if sync {
			//log.Debug().Msgf("Ignoring permission update for channel %s because permissions are synced", new.ID)
			msg = ":arrows_counterclockwise: Channel permissions synced with parent category"
		}

		e := discord.Embed{
			Description: channelChangeHeader(permissionsUpdated, new) + "\n\n" + msg,
			Timestamp:   discord.NowTimestamp(),
			Color:       color.Gold,
		}
		bot.QueueEmbed(audit.ChannelUpdate, new.GuildID, e)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.ChannelUpdateEvent) {
			old, err := s.ChannelStore.Channel(c.ID)
			if err != nil {
				go handleError(
					audit.ChannelUpdate,
					c.GuildID,
					err,
					fmt.Sprintf("Could not retrieve channel from cache: `%s` / `%s`", c.Channel.Name, c.Channel.ID),
					discord.User{},
				)
				return
			}

			go handle(*old, c.Channel)
		})
	})
}

func isPermissionSync(c discord.Channel) (bool, error) {
	// check for permission sync with category
	switch c.Type {
	case discord.GuildPublicThread, discord.GuildPrivateThread:
		// threads can't have permissions
		return false, fmt.Errorf("can't check permission sync for thread %s", c.ID)
	case discord.GuildCategory:
		return false, fmt.Errorf("can't check permission sync for category %s", c.ID)
	}

	if !c.ParentID.IsValid() {
		return false, fmt.Errorf("can't check permission sync for channel %s due to missing parent", c.ID)
	}

	parent, err := s.Channel(c.ParentID)
	if err != nil {
		return false, err
	}

	if parent.Type != discord.GuildCategory {
		return false, fmt.Errorf("can't check permission sync for channel (%s) parent (%s) as their type (%s) is not a category", c.ID, parent.ID, parent.Type)
	}

	return slices.Equal(c.Overwrites, parent.Overwrites), nil
}

func generateOverwriteMessage(old, new discord.Channel) string {
	pairs := make(map[discord.Snowflake]*overwritePair, 0)

	for _, overwrite := range old.Overwrites {
		pairs[overwrite.ID] = &overwritePair{old: overwrite}
	}
	for _, overwrite := range new.Overwrites {
		pair, exists := pairs[overwrite.ID]
		if exists {
			pair.new = overwrite
		} else {
			pairs[overwrite.ID] = &overwritePair{new: overwrite}
		}
	}

	var b strings.Builder
	appendPermissions := func(msg string, perms []util.DiscordPermission) {
		b.WriteString("\n\n")
		b.WriteString(msg)
		b.WriteString(":\n")
		for i, permission := range perms {
			if i != 0 {
				b.WriteString("\n")
			}
			b.WriteString("- ")
			b.WriteString(permission.Name)
		}
	}

	for id, pair := range pairs {
		diff := util.OverwriteDiff(pair.old, pair.new, &new.Type)
		created := !pair.old.ID.IsValid() && pair.new.ID.IsValid()
		deleted := pair.old.ID.IsValid() && !pair.new.ID.IsValid()
		// allow empty diffs for channel creation/deletion (even though this first condition check shouldn't be necessary? idk)
		// TODO maybe remove the created/deleted condition check
		if (!created || !deleted) && len(diff.Allowed) == 0 && len(diff.Neutral) == 0 && len(diff.Denied) == 0 {
			continue
		}
		if b.Len() != 0 {
			b.WriteString("\n\n")
		}

		if deleted {
			b.WriteString(":wastebasket: __**Deleted**__ ")
		} else if created {
			b.WriteString(":inbox_tray: __**Created**__ ")
		} else {
			b.WriteString(":pencil: __**Edited**__ ")
		}

		var t discord.OverwriteType
		if pair.old.ID.IsValid() {
			t = pair.old.Type
		} else if pair.new.ID.IsValid() {
			t = pair.new.Type
		} else {
			log.Warn().
				Interface("pair", pair).
				Interface("oldChannel", old).
				Interface("newChannel", new).
				Msgf("Could not find the type of permission overwrite. This should not happen.")
			t = discord.OverwriteRole
		}

		if t == discord.OverwriteMember {
			b.WriteString("**User:** ")
			uid := discord.UserID(id)
			m, _ := s.Member(new.GuildID, uid)
			if m != nil {
				b.WriteString(util.UserTag(m.User))
			} else {
				b.WriteString(fmt.Sprintf("%s (`%s`)", uid.Mention(), uid.String()))
			}
		} else if t == discord.OverwriteRole {
			b.WriteString("**Role:** ")
			if id == discord.Snowflake(new.GuildID) {
				// @everyone role
				b.WriteString("`@everyone`")
			} else {
				rid := discord.RoleID(id)
				r, _ := s.Role(new.GuildID, rid)
				if r != nil {
					b.WriteString(fmt.Sprintf("`@%s` (`%s`)", r.Name, r.ID))
				} else {
					b.WriteString(fmt.Sprintf("%s (`%s`)", rid.Mention(), rid))
				}
			}
		}

		if len(diff.Allowed) != 0 {
			//✓
			//✅
			appendPermissions("✓ **Allowed permissions**", diff.Allowed)
		}
		if len(diff.Neutral) != 0 && !created && !deleted {
			//⧄
			//🔘
			appendPermissions("⧄ **Neutral permissions**", diff.Neutral)
		}
		if len(diff.Denied) != 0 {
			//✘
			//⛔
			appendPermissions("✘ **Denied permissions**", diff.Denied)
		}
		log.Debug().Interface("diff", diff).Msgf("Got diff for %s", id)
	}

	return b.String()
}
