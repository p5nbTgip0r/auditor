package util

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"time"
)

func GetAuditActioner(
	state *state.State,
	guild discord.GuildID,
	target discord.Snowflake,
	data api.AuditLogData,
	filter func(discord.AuditLogEntry) bool,
) (
	entry discord.AuditLogEntry,
	actionerID discord.UserID,
	actionerMember *discord.Member,
	err error,
) {
	audit, err := state.AuditLog(guild, data)

	if err != nil {
		return
	}

	for _, e := range audit.Entries {
		// only consider entries within 30 seconds
		if e.CreatedAt().Sub(time.Now()).Abs().Seconds() > 30 {
			continue
		}
		if !filter(e) {
			continue
		}
		if e.TargetID == target {
			entry = e
			actionerID = e.UserID
			actionerMember, _ = state.Member(guild, e.UserID)
			break
		}
	}

	return
}
