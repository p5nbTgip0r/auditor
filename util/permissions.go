package util

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/diamondburned/arikawa/v3/discord"
	"strings"
)

type DiscordPermission struct {
	Permission discord.Permissions
	Name       string
}

var (
	Permissions = []DiscordPermission{
		{discord.PermissionCreateInstantInvite, "Create Instant Invite"},
		{discord.PermissionKickMembers, "Kick Members"},
		{discord.PermissionBanMembers, "Ban Members"},
		{discord.PermissionAdministrator, "Administrator"},
		{discord.PermissionManageChannels, "Manage Channels"},
		{discord.PermissionManageGuild, "Manage Guild"},
		{discord.PermissionAddReactions, "Add Reactions"},
		{discord.PermissionViewAuditLog, "View Audit Log"},
		{discord.PermissionPrioritySpeaker, "Priority Speaker"},
		{discord.PermissionStream, "Stream"},
		{discord.PermissionViewChannel, "View Channel"},
		{discord.PermissionSendMessages, "Send Messages"},
		{discord.PermissionSendTTSMessages, "Send TTS Messages"},
		{discord.PermissionManageMessages, "Manage Messages"},
		{discord.PermissionEmbedLinks, "Embed Links"},
		{discord.PermissionAttachFiles, "Attach Files"},
		{discord.PermissionReadMessageHistory, "Read Message History"},
		{discord.PermissionMentionEveryone, "Mention Everyone"},
		{discord.PermissionUseExternalEmojis, "Use External Emojis"},
		// arikawa is missing this one
		{1 << 19, "View Guild Insights"},
		{discord.PermissionConnect, "Connect"},
		{discord.PermissionSpeak, "Speak"},
		{discord.PermissionMuteMembers, "Mute Members"},
		{discord.PermissionDeafenMembers, "Deafen Members"},
		{discord.PermissionMoveMembers, "Move Members"},
		{discord.PermissionUseVAD, "Use VAD"},
		{discord.PermissionChangeNickname, "Change Nickname"},
		{discord.PermissionManageNicknames, "Manage Nicknames"},
		{discord.PermissionManageRoles, "Manage Roles"},
		{discord.PermissionManageWebhooks, "Manage Webhooks"},
		{discord.PermissionManageEmojisAndStickers, "Manage Emojis And Stickers"},
		{discord.PermissionUseSlashCommands, "Use Slash Commands"},
		{discord.PermissionRequestToSpeak, "Request to Speak"},
		{discord.PermissionManageEvents, "Manage Events"},
		{discord.PermissionManageThreads, "Manage Threads"},
		{discord.PermissionCreatePublicThreads, "Create Public Threads"},
		{discord.PermissionCreatePrivateThreads, "Create Private Threads"},
		{discord.PermissionUseExternalStickers, "Use External Stickers"},
		{discord.PermissionSendMessagesInThreads, "Send Messages in Threads"},
		{discord.PermissionStartEmbeddedActivities, "Start Embedded Activities"},
		{discord.PermissionModerateMembers, "Moderate Members"},
	}

	TextChannelPermissions = []DiscordPermission{
		{discord.PermissionCreateInstantInvite, "Create Instant Invite"},
		{discord.PermissionManageChannels, "Manage Channel"},
		{discord.PermissionAddReactions, "Add Reactions"},
		{discord.PermissionViewChannel, "View Channel"},
		{discord.PermissionSendMessages, "Send Messages"},
		{discord.PermissionSendTTSMessages, "Send TTS Messages"},
		{discord.PermissionManageMessages, "Manage Messages"},
		{discord.PermissionEmbedLinks, "Embed Links"},
		{discord.PermissionAttachFiles, "Attach Files"},
		{discord.PermissionReadMessageHistory, "Read Message History"},
		{discord.PermissionMentionEveryone, "Mention Everyone"},
		{discord.PermissionUseExternalEmojis, "Use External Emojis"},
		{discord.PermissionManageRoles, "Manage Permissions"},
		{discord.PermissionManageWebhooks, "Manage Webhooks"},
		{discord.PermissionUseSlashCommands, "Use Slash Commands"},
		{discord.PermissionManageThreads, "Manage Threads"},
		{discord.PermissionCreatePublicThreads, "Create Public Threads"},
		{discord.PermissionCreatePrivateThreads, "Create Private Threads"},
		{discord.PermissionUseExternalStickers, "Use External Stickers"},
		{discord.PermissionSendMessagesInThreads, "Send Messages in Threads"},
	}

	VoiceChannelPermissions = []DiscordPermission{
		{discord.PermissionCreateInstantInvite, "Create Instant Invite"},
		{discord.PermissionManageChannels, "Manage Channels"},
		{discord.PermissionAddReactions, "Add Reactions"},
		{discord.PermissionPrioritySpeaker, "Priority Speaker"},
		{discord.PermissionStream, "Stream"},
		{discord.PermissionViewChannel, "View Channel"},
		{discord.PermissionSendMessages, "Send Messages"},
		{discord.PermissionSendTTSMessages, "Send TTS Messages"},
		{discord.PermissionManageMessages, "Manage Messages"},
		{discord.PermissionEmbedLinks, "Embed Links"},
		{discord.PermissionAttachFiles, "Attach Files"},
		{discord.PermissionReadMessageHistory, "Read Message History"},
		{discord.PermissionMentionEveryone, "Mention Everyone"},
		{discord.PermissionUseExternalEmojis, "Use External Emojis"},
		{discord.PermissionConnect, "Connect"},
		{discord.PermissionSpeak, "Speak"},
		{discord.PermissionMuteMembers, "Mute Members"},
		{discord.PermissionDeafenMembers, "Deafen Members"},
		{discord.PermissionMoveMembers, "Move Members"},
		{discord.PermissionUseVAD, "Use Voice Activity"},
		{discord.PermissionManageRoles, "Manage Permissions"},
		{discord.PermissionManageWebhooks, "Manage Webhooks"},
		{discord.PermissionUseSlashCommands, "Use Slash Commands"},
		{discord.PermissionManageEvents, "Manage Events"},
		{discord.PermissionUseExternalStickers, "Use External Stickers"},
		{discord.PermissionStartEmbeddedActivities, "Start Embedded Activities"},
	}

	StageChannelPermissions = []DiscordPermission{
		{discord.PermissionCreateInstantInvite, "Create Instant Invite"},
		{discord.PermissionViewChannel, "View Channel"},
		{discord.PermissionMentionEveryone, "Mention Everyone"},
		{discord.PermissionConnect, "Connect"},
		{discord.PermissionMuteMembers, "Mute Members"},
		{discord.PermissionDeafenMembers, "Deafen Members"},
		{discord.PermissionMoveMembers, "Move Members"},
		{discord.PermissionManageRoles, "Manage Permissions"},
		{discord.PermissionRequestToSpeak, "Request to Speak"},
		{discord.PermissionManageEvents, "Manage Events"},
	}
)

func PermissionString(permissions discord.Permissions) string {
	if permissions == 0 {
		return "No permissions"
	} else if permissions == discord.PermissionAll {
		return "All permissions"
	}

	var builder strings.Builder

	for _, perm := range Permissions {
		if !permissions.Has(perm.Permission) {
			continue
		}
		if builder.Len() != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(perm.Name)
	}

	return builder.String()
}

func PermissionSliceString(permissions []DiscordPermission) string {
	if len(permissions) == 0 {
		return "No permissions"
	} else if len(permissions) == len(Permissions) {
		// duplicate entries can mess this up, but hopefully there's no duplicates in the first place
		return "All permissions"
	}

	var builder strings.Builder

	for _, perm := range permissions {
		if builder.Len() != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(perm.Name)
	}

	return builder.String()
}

func PermissionSetString(permissions mapset.Set[DiscordPermission]) string {
	if permissions.Cardinality() == 0 {
		return "No permissions"
	} else if permissions.Cardinality() == len(Permissions) {
		return "All permissions"
	}

	var builder strings.Builder

	for perm := range permissions.Iter() {
		if builder.Len() != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(perm.Name)
	}

	return builder.String()
}

func GetPermissions(permissions discord.Permissions) []DiscordPermission {
	if permissions == 0 {
		return make([]DiscordPermission, 0)
	} else if permissions == discord.PermissionAll {
		dupe := make([]DiscordPermission, len(Permissions))
		copy(dupe, Permissions)
		return dupe
	}

	perms := make([]DiscordPermission, 0)

	for _, perm := range Permissions {
		if !permissions.Has(perm.Permission) {
			continue
		}
		perms = append(perms, perm)
	}

	return perms
}
