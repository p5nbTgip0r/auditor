package util

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strings"
)

type DiscordPermission struct {
	Permission discord.Permissions
	Name       string
}

var (
	PermissionName = []DiscordPermission{
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
)

func PermissionString(permissions discord.Permissions) string {
	if permissions == 0 {
		return "No permissions"
	} else if permissions == discord.PermissionAll {
		return "All permissions"
	}

	var builder strings.Builder

	for _, perm := range PermissionName {
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
