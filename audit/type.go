package audit

//go:generate stringer -type=Type
type Type uint

const (
	MessageDelete Type = iota
	MessageUpdate
	MessagePurge
	MemberNickname
	MemberAvatar
	MemberRoles
	MemberTimeout
	MemberScreening
	MemberJoin
	MemberLeave
	MemberBan
	MemberUnban
	MemberKick
	RoleCreate
	RoleUpdate
	RoleDelete
	ServerEdited
	ServerEmoji
	UserName
	UserAvatar
	ChannelCreate
	ChannelUpdate
	ChannelDelete
	InviteSend
	InviteCreate
	InviteDelete
	VoiceConnection
	VoiceAudioState
)
