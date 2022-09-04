package audit

import "strings"

//go:generate stringer -type=Type
type Type uint

const (
	Unknown Type = iota
	MessageDelete
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

// UnString is adapted from:
// https://stackoverflow.com/a/61310395
func UnString(s string) Type {
	s = strings.ToLower(s)
	l := strings.ToLower(_Type_name)
	for i := 0; i < len(_Type_index)-1; i++ {
		//goland:noinspection GoRedundantConversion
		i2 := Type(i)
		output := l[_Type_index[i2]:_Type_index[i2+1]]
		if s == output {
			return i2
		}
	}
	return 0
}
