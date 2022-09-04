package audit

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
