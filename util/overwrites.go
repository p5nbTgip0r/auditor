package util

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type OverwritePermissions struct {
	Allowed []DiscordPermission
	Neutral []DiscordPermission
	Denied  []DiscordPermission
}

// OverwriteDiff takes in two discord.Overwrite structs and returns a OverwritePermissions struct showing the
// difference in permission states between the two.
//
// - Only one of the discord.Overwrite parameters is required, the other can be left as a default struct for channel creation/deletion events.
//
// For channel creation or deletion, permission overwrites will be passed through into the final OverwritePermissions struct
// (i.e. this function will return all allowed, neutral, and denied permissions specified in the overwrite)
//
// - Both discord.Overwrite parameters must reference the same target (ignored for channel creation/deletion)
// - (Optional) A discord.ChannelType can be specified to adjust the permission checking.
func OverwriteDiff(old, new discord.Overwrite, channelType *discord.ChannelType) (final OverwritePermissions) {
	final = OverwritePermissions{}
	// new or old overwrite target ID must be valid
	if !new.ID.IsValid() && !old.ID.IsValid() {
		return
	}

	// cannot diff overrides with differing targets
	if old.ID.IsValid() && new.ID.IsValid() && old.ID != new.ID {
		return
	}

	perms := &Permissions
	if channelType != nil {
		// swaps out the full permissions slice with one that's only applicable for that channel type.
		// this avoids needlessly checking permissions that can't be applied as an overwrite, and the
		// permission names are adjusted to how the discord client displays them.
		switch *channelType {
		case discord.GuildText:
			perms = &TextChannelPermissions
		case discord.GuildVoice:
			perms = &VoiceChannelPermissions
		case discord.GuildStageVoice:
			perms = &StageChannelPermissions
		}
	}

	for _, perm := range *perms {
		if !old.ID.IsValid() {
			// created, pass all permissions through
			switch {
			case new.Allow.Has(perm.Permission):
				final.Allowed = append(final.Allowed, perm)
			case new.Deny.Has(perm.Permission):
				final.Denied = append(final.Denied, perm)
			default:
				final.Neutral = append(final.Neutral, perm)
			}
		} else if !new.ID.IsValid() {
			// deleted, pass all permissions through
			switch {
			case old.Allow.Has(perm.Permission):
				final.Allowed = append(final.Allowed, perm)
			case old.Deny.Has(perm.Permission):
				final.Denied = append(final.Denied, perm)
			default:
				final.Neutral = append(final.Neutral, perm)
			}
		} else {
			var allow, neutral, deny bool
			switch {
			case old.Allow.Has(perm.Permission) && new.Deny.Has(perm.Permission):
				// was allowed, now denied
				deny = true

			case old.Deny.Has(perm.Permission) && new.Allow.Has(perm.Permission):
				// was denied, now allowed
				allow = true

			case old.Deny.Has(perm.Permission) && !(new.Allow.Has(perm.Permission) || new.Deny.Has(perm.Permission)):
				// was denied, now neutral
				neutral = true

			case old.Allow.Has(perm.Permission) && !(new.Allow.Has(perm.Permission) || new.Deny.Has(perm.Permission)):
				// was allowed, now neutral
				neutral = true

			case new.Allow.Has(perm.Permission) && !(old.Allow.Has(perm.Permission) || old.Deny.Has(perm.Permission)):
				// was neutral, now allowed
				allow = true

			case new.Deny.Has(perm.Permission) && !(old.Allow.Has(perm.Permission) || old.Deny.Has(perm.Permission)):
				// was neutral, now denied
				deny = true
			}

			switch {
			case allow:
				final.Allowed = append(final.Allowed, perm)
			case neutral:
				final.Neutral = append(final.Neutral, perm)
			case deny:
				final.Denied = append(final.Denied, perm)
			}
		}
	}

	return
}
