package schema

import (
	"audit/audit"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Guild struct {
	ID                 discord.GuildID   `bson:"guildID"`
	AuditChannelID     discord.ChannelID `bson:"auditChannelID"`
	LoggingDisabled    bool              `bson:"loggingDisabled"`
	DisabledAuditTypes []audit.Type      `bson:"disabledAuditTypes"`
}

func (g *Guild) DisabledAuditTypesSet() mapset.Set[audit.Type] {
	dat := mapset.NewSet[audit.Type]()
	for _, auditType := range g.DisabledAuditTypes {
		dat.Add(auditType)
	}
	return dat
}

func (g *Guild) ChangeAuditTypes(types mapset.Set[audit.Type], enable bool) {
	dat := g.DisabledAuditTypesSet()
	for t := range types.Iter() {
		if enable {
			dat.Remove(t)
		} else {
			dat.Add(t)
		}
	}
	g.DisabledAuditTypes = dat.ToSlice()
}
