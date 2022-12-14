package events

import (
	"audit/audit"
	"audit/bot"
	"audit/database"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

var (
	s *state.State

	handler []func()
)

func InitEventHandlers(state *state.State) {
	s = state
	for _, event := range handler {
		event()
	}
}

func handleError(auditType audit.Type, g discord.GuildID, err error, msg string, user discord.User) {
	embed := errorEmbed(auditType, msg, user)
	log.Warn().
		Err(err).
		Interface("auditType", auditType.String()).
		Interface("guildID", g).
		Interface("user", user).
		Msgf("Encountered error from audit type '%s': %s", auditType.String(), msg)
	bot.QueueEmbed(auditType, g, *embed)
}

func errorEmbed(auditType audit.Type, msg string, user discord.User) *discord.Embed {
	var e *discord.Embed
	if user.ID.IsValid() {
		e = userBaseEmbed(user, "", false)
	} else {
		e = &discord.Embed{}
	}

	e.Description = fmt.Sprintf("**:warning: Error when creating audit message for __%s__**\n\n:pencil: %s", auditType, msg)
	e.Color = color.Red

	return e
}

func userBaseEmbed(user discord.User, url string, userUpdate bool) *discord.Embed {
	e := &discord.Embed{}
	e.Author = &discord.EmbedAuthor{
		Name: user.Tag(),
		URL:  url,
		Icon: user.AvatarURL(),
	}
	e.Timestamp = discord.NowTimestamp()
	if userUpdate {
		e.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("User ID: %s", user.ID)}
		e.Thumbnail = &discord.EmbedThumbnail{URL: user.AvatarURL()}
	}
	return e
}

// check looks whether the given audit type is enabled for the given guild and channel IDs.
// This function blocks to access the database.
func check(a audit.Type, g *discord.GuildID, c *discord.ChannelID) bool {
	if g == nil || !g.IsValid() {
		return false
	}
	sg, err := database.Collections.Guilds.GetGuild(*g)
	if err != nil || sg.LoggingDisabled || !sg.AuditChannelID.IsValid() {
		return false
	}

	return !slices.Contains(sg.DisabledAuditTypes, a)
}
