package audit

import "github.com/diamondburned/arikawa/v3/discord"

//go:generate stringer -type=AuditType -trimprefix=Audit
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
	AuditVoiceConnection
	AuditVoiceAudioState
)

func (t AuditType) Check(g *discord.GuildID, c *discord.ChannelID) bool {
	// todo: actually implement this
	return true
}
