package schema

import (
	"audit/events"
)

type Guild struct {
	ID                 uint64             `bson:"guildID"`
	AuditChannelID     uint64             `bson:"auditChannelID"`
	LoggingDisabled    bool               `bson:"loggingDisabled"`
	DisabledAuditTypes []events.AuditType `bson:"disabledAuditTypes"`
}
