package schema

import (
	"audit/audit"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Guild struct {
	ID                 discord.GuildID   `bson:"guildID"`
	AuditChannelID     discord.ChannelID `bson:"auditChannelID"`
	LoggingDisabled    bool              `bson:"loggingDisabled"`
	DisabledAuditTypes []audit.AuditType `bson:"disabledAuditTypes"`
}
