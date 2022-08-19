package events

func init() {
	//handler = append(handler, func() {
	//	s.PreHandler.AddSyncHandler(func(c *gateway.GuildMemberRemoveEvent) {
	//		m, _ := s.MemberStore.Member(c.GuildID, c.User.ID)
	//		_, bannerMember, _ := getAuditActioner(c.GuildID, discord.Snowflake(c.User.ID), api.AuditLogData{ActionType: discord.MemberKick})
	//
	//		var joined string
	//		if m != nil {
	//			joined = fmt.Sprintf("<t:%d:R>", m.Joined.Time().Unix())
	//		} else {
	//			joined = "Could not retrieve"
	//		}
	//
	//		embed := userBaseEmbed(c.User, "", true)
	//		embed.Color = 0x992D22
	//
	//		embed.Description = fmt.Sprintf("**:boot: %s was kicked**", c.User.Mention())
	//		embed.Fields = append(embed.Fields,
	//			discord.EmbedField{
	//				Name:  "Joined server",
	//				Value: joined,
	//			},
	//		)
	//		if bannerMember != nil {
	//			embed.Fields = append(embed.Fields,
	//				discord.EmbedField{
	//					Name:  "Kicked by",
	//					Value: util.FullTag(bannerMember.User),
	//				},
	//			)
	//		}
	//
	//		handleAuditError(s.SendMessage(auditChannel, "", *embed))
	//	})
	//})
}
