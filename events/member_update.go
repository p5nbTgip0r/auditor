package events

import (
	"audit/audit"
	"audit/util"
	"audit/util/color"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type updatedUserFields uint

const (
	fieldMemberNickname updatedUserFields = 1 << iota
	fieldMemberRoles
	fieldMemberAvatar
	fieldMemberPending
	fieldMemberTimeout
	fieldUserName
	fieldUserDiscriminator
	fieldUserAvatar
	numFields = 8
)

// Has returns true if u has the given fields set.
func (u updatedUserFields) Has(fields updatedUserFields) bool {
	return discord.HasFlag(uint64(u), uint64(fields))
}

// Name returns the name of the field. If the index of the field is known, that can be specified to skip a loop
func (u updatedUserFields) Name(index *uint) string {
	names := [numFields]string{
		"MemberNickname",
		"MemberRoles",
		"MemberAvatar",
		"MemberPending",
		"MemberTimeout",
		"UserName",
		"UserDiscriminator",
		"UserAvatar",
	}

	if index != nil && *index < numFields {
		return names[*index]
	}

	for i := 0; i < numFields; i++ {
		resolved := updatedUserFields(1 << i)
		if resolved == u {
			return names[i]
		}
	}

	return ""
}

// String returns a string representation of the field(s) updated
func (u updatedUserFields) String() string {
	var b strings.Builder
	var i uint
	b.WriteString("{")
	for i = 0; i < numFields; i++ {
		resolved := updatedUserFields(1 << i)
		if u.Has(resolved) {
			if b.Len() != 1 {
				b.WriteString(", ")
			}
			b.WriteString(resolved.Name(&i))
			b.WriteString(" (0b")
			b.WriteString(strconv.FormatInt(int64(resolved), 2))
			b.WriteString(")")
		}
	}
	b.WriteString("}")
	return b.String()
}

type userDiff struct {
	fields       updatedUserFields
	addedRoles   mapset.Set[discord.RoleID]
	removedRoles mapset.Set[discord.RoleID]
}

func init() {
	handle := func(old discord.Member, event gateway.GuildMemberUpdateEvent) {
		fields := determineUpdatedFields(old, event)
		log.Info().Msgf("Fields: %s", fields.fields.String())

		log.Debug().
			Interface("event", event).
			Interface("member", old).
			Msgf("Member updated: %s", event.User.Tag())

		diffMember(event, old, fields)
		diffUser(event, old.User, fields)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(e *gateway.GuildMemberUpdateEvent) {
			old, err := s.MemberStore.Member(e.GuildID, e.User.ID)
			if err != nil {
				go handleError(
					audit.AuditMemberJoin,
					err,
					"Could not retrieve member from cache: "+util.UserTag(e.User),
					&e.User,
				)
				return
			}

			go handle(*old, *e)
		})
	})
}

func determineUpdatedFields(old discord.Member, new gateway.GuildMemberUpdateEvent) userDiff {
	var u updatedUserFields

	// Member nickname
	if old.Nick != new.Nick {
		u |= fieldMemberNickname
	}
	// Member roles
	addRoles, remRoles := util.SliceDiff(old.RoleIDs, new.RoleIDs)
	if addRoles.Cardinality() != 0 || remRoles.Cardinality() != 0 {
		u |= fieldMemberRoles
	}
	// Member Avatar
	if old.Avatar != new.Avatar {
		u |= fieldMemberAvatar
	}

	// Member Pending
	if old.IsPending != new.IsPending {
		u |= fieldMemberPending
	}

	// Member Timeout
	if !old.CommunicationDisabledUntil.Time().Equal(new.CommunicationDisabledUntil.Time()) {
		u |= fieldMemberTimeout
	}

	// User Name
	if old.User.Username != new.User.Username {
		u |= fieldUserName
	}
	// User Discriminator
	if old.User.Discriminator != new.User.Discriminator {
		u |= fieldUserDiscriminator
	}
	// User Avatar
	if old.User.Avatar != new.User.Avatar {
		u |= fieldUserAvatar
	}

	return userDiff{
		fields:       u,
		addedRoles:   addRoles,
		removedRoles: remRoles,
	}
}

func diffMember(new gateway.GuildMemberUpdateEvent, old discord.Member, diff userDiff) {
	getEmbed := func(desc string) *discord.Embed {
		c := userBaseEmbed(new.User, "", true)
		c.Color = color.Gold
		c.Description = desc
		return c
	}

	if diff.fields.Has(fieldMemberNickname) && check(audit.AuditMemberNickname, &new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:pencil: %s nickname edited**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old nickname", Value: old.Nick},
			discord.EmbedField{Name: "New nickname", Value: new.Nick},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	if diff.fields.Has(fieldMemberTimeout) && check(audit.AuditMemberTimeout, &new.GuildID, nil) {
		var c *discord.Embed
		if new.CommunicationDisabledUntil.IsValid() {
			c = getEmbed(fmt.Sprintf("**:zipper_mouth: %s was timed out**", new.User.Mention()))
			c.Fields = append(c.Fields,
				discord.EmbedField{
					Name:  "Timeout Expiry",
					Value: util.Timestamp(new.CommunicationDisabledUntil.Time(), util.Relative),
				},
			)
		} else if old.CommunicationDisabledUntil.Time().After(time.Now()) {
			// it's sometimes possible to remove a timeout after it's expired.
			// this `else` statement will only be reached if the timeout was removed, so this just ensures the old timeout
			// didn't expire already.
			c = getEmbed(fmt.Sprintf("**:zipper_mouth: %s's timeout was removed**", new.User.Mention()))
		}

		if c != nil {
			msg, err := s.SendEmbeds(auditChannel, *c)
			handleAuditError(msg, err, *c)
		}
	}

	if diff.fields.Has(fieldMemberPending) && check(audit.AuditMemberScreening, &new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:clipboard: %s completed membership screening**", new.User.Mention()))
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	if diff.fields.Has(fieldMemberAvatar) && check(audit.AuditMemberAvatar, &new.GuildID, nil) {
		c := getEmbed(fmt.Sprintf("**:frame_photo: %s changed their __guild__ avatar**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old avatar", Value: old.AvatarURL(new.GuildID)},
			discord.EmbedField{Name: "New avatar", Value: new.Avatar},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	if diff.fields.Has(fieldMemberRoles) && check(audit.AuditMemberRoles, &new.GuildID, nil) {
		addRoles, remRoles := diff.addedRoles, diff.removedRoles
		roleNames := make(map[discord.RoleID]string)
		guild, err := s.Guild(new.GuildID)
		if err == nil {
			for _, role := range guild.Roles {
				roleNames[role.ID] = role.Name
			}
		} else {
			log.Err(err).
				Interface("event", new).
				Interface("oldMember", old).
				Msg("couldn't retrieve guild for role information")
		}

		embed := getEmbed(fmt.Sprintf("**:crossed_swords: %s roles have changed**", new.User.Mention()))
		if addRoles.Cardinality() != 0 {
			embed.Fields = append(embed.Fields, discord.EmbedField{
				Name:   "Added roles",
				Value:  formatRoles(addRoles, roleNames),
				Inline: false,
			})
		}

		if remRoles.Cardinality() != 0 {
			embed.Fields = append(embed.Fields, discord.EmbedField{
				Name:   "Removed roles",
				Value:  formatRoles(remRoles, roleNames),
				Inline: false,
			})
		}

		handleAuditError(s.SendEmbeds(auditChannel, *embed))
	}
}

func formatRoles(r mapset.Set[discord.RoleID], names map[discord.RoleID]string) string {
	var msg strings.Builder
	for id := range r.Iter() {
		if msg.Len() != 0 {
			msg.WriteString("; ")
		}
		msg.WriteString("`")
		if name, ok := names[id]; ok {
			msg.WriteString(name)
		} else {
			// fallback to role ID
			msg.WriteString(id.String())
		}
		msg.WriteString("`")
	}
	return msg.String()
}
