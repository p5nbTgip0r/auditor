package events

func init() {
	//handler = append(handler, func() {
	//	s.PreHandler.AddSyncHandler(func(c *gateway.GuildBanAddEvent) {
	//		m, _ := s.MemberStore.Member(c.GuildID, c.User.ID)
	//		// todo: we can goroutine after this point, since we've grabbed all we need from the cache
	//		_, bannerMember, _ := getAuditActioner(c.GuildID, discord.Snowflake(c.User.ID), api.AuditLogData{ActionType: discord.MemberBanAdd})
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
	//		embed.Description = fmt.Sprintf("**:rotating_light: %s was banned**", c.User.Mention())
	//		embed.Fields = append(embed.Fields,
	//			discord.EmbedField{
	//				Name:  "Joined server",
	//				Value: joined,
	//			},
	//		)
	//		if bannerMember != nil {
	//			embed.Fields = append(embed.Fields,
	//				discord.EmbedField{
	//					Name:  "Banned by",
	//					Value: util.FullTag(bannerMember.User),
	//				},
	//			)
	//		}
	//
	//		handleAuditError(s.SendMessage(auditChannel, "", *embed))
	//	})
	//})
}
