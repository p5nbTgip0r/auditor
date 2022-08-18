package events

import (
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
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(e *gateway.GuildMemberUpdateEvent) {
			old, err := s.MemberStore.Member(e.GuildID, e.User.ID)
			if err != nil {
				log.Warn().
					Err(err).
					Interface("event", e).
					Msgf("Member was updated, but could not retrieve previous state: %s", e.User.Tag())
			} else {
				fields := determineUpdatedFields(old, e)
				log.Info().Msgf("Fields: %s", fields.fields.String())

				log.Debug().
					Interface("event", e).
					Interface("member", old).
					Msgf("Member updated: %s", e.User.Tag())

				diffMember(e, *old, fields)
				diffUser(e, old.User, fields)
			}
		})
	})
}

func determineUpdatedFields(old *discord.Member, new *gateway.GuildMemberUpdateEvent) userDiff {
	var u updatedUserFields

	// Member nickname
	if old.Nick != new.Nick {
		u |= fieldMemberNickname
	}
	// Member roles
	addRoles, remRoles := roleDiff(old.RoleIDs, new.RoleIDs)
	if addRoles.Cardinality() != 0 || remRoles.Cardinality() != 0 {
		u |= fieldMemberRoles
	}
	// Member Avatar
	if old.Avatar != new.Avatar {
		u |= fieldMemberAvatar
	}

	// Member Pending (todo, requires arikawa to add the corresponding field in the Member Update event payload)

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

func diffMember(new *gateway.GuildMemberUpdateEvent, old discord.Member, diff userDiff) {
	getEmbed := func(desc string) *discord.Embed {
		c := userBaseEmbed(new.User, "", true)
		// https://github.com/Rapptz/discord.py/blob/de941ababe9da898dd62d2b2a2d21aaecac6bd09/discord/colour.py#L295
		c.Color = 0xf1c40f
		c.Description = desc
		return c
	}

	if diff.fields.Has(fieldMemberNickname) {
		c := getEmbed(fmt.Sprintf("**:pencil: %s nickname edited**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old nickname", Value: old.Nick},
			discord.EmbedField{Name: "New nickname", Value: new.Nick},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	if diff.fields.Has(fieldMemberTimeout) {
		var c *discord.Embed
		if new.CommunicationDisabledUntil.IsValid() {
			c = getEmbed(fmt.Sprintf("**:zipper_mouth: %s was timed out**", new.User.Mention()))
			c.Fields = append(c.Fields,
				discord.EmbedField{
					Name:  "Timeout Expiry",
					Value: fmt.Sprintf("<t:%d:R>", new.CommunicationDisabledUntil.Time().Unix()),
				},
			)
			// it's sometimes possible to remove a timeout after it's expired.
			// this `else` statement will only be reached if the timeout was removed, so this just ensures the old timeout
			// didn't expire already.
		} else if old.CommunicationDisabledUntil.Time().After(time.Now()) {
			c = getEmbed(fmt.Sprintf("**:zipper_mouth: %s's timeout was removed**", new.User.Mention()))
		}

		if c != nil {
			msg, err := s.SendEmbeds(auditChannel, *c)
			handleAuditError(msg, err, *c)
		}
	}

	if diff.fields.Has(fieldMemberAvatar) {
		c := getEmbed(fmt.Sprintf("**:frame_photo: %s changed their __guild__ avatar**", new.User.Mention()))
		c.Fields = append(c.Fields,
			discord.EmbedField{Name: "Old avatar", Value: old.AvatarURL(new.GuildID)},
			discord.EmbedField{Name: "New avatar", Value: new.Avatar},
		)
		msg, err := s.SendEmbeds(auditChannel, *c)
		handleAuditError(msg, err, *c)
	}

	addRoles, remRoles := diff.addedRoles, diff.removedRoles
	if diff.fields.Has(fieldMemberRoles) {
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
	iter := r.Iterator()
	for id := range iter.C {
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

func roleDiff(old, new []discord.RoleID) (added, removed mapset.Set[discord.RoleID]) {
	oldSet := mapset.NewSet[discord.RoleID]()
	newSet := mapset.NewSet[discord.RoleID]()
	for _, id := range old {
		oldSet.Add(id)
	}
	for _, id := range new {
		newSet.Add(id)
	}

	add := newSet.Difference(oldSet)
	rem := oldSet.Difference(newSet)
	return add, rem
}
