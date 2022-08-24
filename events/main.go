package events

import (
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
)

var (
	s *state.State

	handler []func()
)

const auditChannel = 670908023388241931

//go:generate go run golang.org/x/tools/cmd/stringer -type=AuditType -trimprefix=Audit
type AuditType uint

const (
	AuditMessageDelete AuditType = iota
	AuditMessageUpdate
	AuditMessagePurge
	AuditMemberNickname
	AuditMemberAvatar
	AuditMemberRoles
	AuditMemberTimeout
	AuditMemberScreening
	AuditMemberJoin
	AuditMemberLeave
	AuditMemberBan
	AuditMemberUnban
	AuditMemberKick
	AuditRoleCreate
	AuditRoleUpdate
	AuditRoleDelete
	AuditServerEdited
	AuditServerEmoji
	AuditUserName
	AuditUserAvatar
	AuditChannelCreate
	AuditChannelUpdate
	AuditChannelDelete
	AuditInviteSend
	AuditInviteCreate
	AuditInviteDelete
	AuditVoiceState
)

func InitEventHandlers(state *state.State) {
	s = state
	for _, event := range handler {
		event()
	}
}

func handleAuditError(msg *discord.Message, err error, embeds ...discord.Embed) {
	if err != nil {
		log.Err(err).
			Interface("embeds", embeds).
			Msg("Could not send log message")
	} else {
		log.Debug().
			Interface("message", msg).
			Msgf("Successfully sent log message")
	}
}

func handleError(auditType AuditType, err error, msg string, user *discord.User) {
	embed := errorEmbed(auditType, msg, user)
	handleAuditError(s.SendEmbeds(auditChannel, *embed))
}

func errorEmbed(auditType AuditType, msg string, user *discord.User) *discord.Embed {
	var e *discord.Embed
	if user != nil {
		e = userBaseEmbed(*user, "", false)
	} else {
		e = &discord.Embed{}
	}

	e.Description = fmt.Sprintf("**:warning: Error when creating audit message for __%s__**\n\n:pencil: %s", auditType, msg)
	e.Color = color.Red

	return e
}

func userBaseEmbed(user discord.User, url string, userUpdate bool) *discord.Embed {
	e := &discord.Embed{}
	e.Author = &discord.EmbedAuthor{
		Name: user.Tag(),
		URL:  url,
		Icon: user.AvatarURL(),
	}
	e.Timestamp = discord.NowTimestamp()
	if userUpdate {
		e.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("User ID: %s", user.ID)}
		e.Thumbnail = &discord.EmbedThumbnail{URL: user.AvatarURL()}
	}
	return e
}

func (t AuditType) check(g *discord.GuildID, c *discord.ChannelID) bool {
	// todo: actually implement this
	return true
}
