package events

import (
	"audit/audit"
	"audit/bot"
	"audit/util/color"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

func init() {
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.MessageDeleteEvent) {
			if !c.GuildID.IsValid() {
				return
			}
			m, err := s.Message(c.ChannelID, c.ID)
			if err != nil {
				log.Warn().
					Err(err).
					Interface("event", c).
					Msgf("Message was deleted, but not found in cache: %s", c.ID)

				go func() {
					if !check(audit.MessageDelete, &c.GuildID, &c.ChannelID) {
						return
					}
					desc := fmt.Sprintf("**:wastebasket: Message deleted from %s:**\n\n:warning: Message details could not be retrieved from cache.", c.ChannelID.Mention())
					embeds := deletedMessageEmbeds(desc, c.ID, c.ChannelID, nil, nil, color.Gold)
					bot.QueueEmbed(audit.MessageDelete, c.GuildID, embeds...)
				}()
			} else {
				if m.Author.Bot {
					// ignore bot messages
					return
				}
				go deletedMessageLogs(m)
			}
		})
	})
}

func httpDown(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//defer res.Body.Close()

	// todo: check the size first, then dump to a temporary file on the disk if it's large enough
	//d, err := io.ReadAll(res.Body)
	return res.Body, err
}

func deletedMessageLogs(m *discord.Message) {
	if !check(audit.MessageDelete, &m.GuildID, &m.ChannelID) {
		return
	}

	log.Debug().
		Interface("msg", m).
		Msgf("[Deleted] %s: %s", m.Author.Username, m.Content)

	mContent := m.Content
	if m.Content == "" {
		mContent = "Message has no content."
	}
	desc := fmt.Sprintf("**:wastebasket: Message deleted from %s:**\n\n%s", m.ChannelID.Mention(), mContent)

	// todo: this is BAD. if a 100mb file is sent, the whole file is gonna be loaded into memory *and* this whole handler
	//thing is gonna block until the download is finished. i gotta clean this up, but it's a 1:1 implementation from the
	//python version at the moment
	// todo update: this is a little better now, but still kind of fucky. i've modified it so i'm just passing the reader
	// into the file, then it won't load the whole thing into ram first
	var attachmentsField strings.Builder
	var files []sendpart.File
	var fields []discord.EmbedField
	for _, att := range m.Attachments {
		if attachmentsField.Len() != 0 {
			attachmentsField.WriteString("\n")
		}
		attachmentsField.WriteString(fmt.Sprintf("[%s](%s) [**`Alt Link`**](%s)", att.Filename, att.URL, att.Proxy))

		res, err := http.Get(att.URL)
		if err != nil {
			log.Err(err).Msg("failed downloading attachment")
			continue
		}

		// there's not really an easy way to avoid this warning afaik. we're using the reader for this whole function
		defer res.Body.Close()

		files = append(files, sendpart.File{
			Name:   att.Filename,
			Reader: res.Body,
		})
	}

	if len(m.Attachments) != 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Attachments",
			Value: attachmentsField.String(),
		})
	}

	embeds := deletedMessageEmbeds(desc, m.ID, m.ChannelID, &m.Author, fields, color.DarkerGray)

	if len(m.Attachments) != 0 {
		err := bot.ProcessMessage(bot.AuditMessage{
			AuditType: audit.MessageDelete,
			GuildID:   m.GuildID,
			SendMessageData: api.SendMessageData{
				Embeds: embeds,
				Files:  files,
			},
		})
		if err == nil {
			return
		}
		// fallback to the normal message handling if sending attachments failed
	}
	bot.QueueEmbed(audit.MessageDelete, m.GuildID, embeds...)
}

func deletedMessageEmbeds(
	desc string,
	messageID discord.MessageID,
	channelID discord.ChannelID,
	author *discord.User,
	fields []discord.EmbedField,
	color discord.Color,
) []discord.Embed {
	var eAuthor *discord.EmbedAuthor
	if author != nil {
		eAuthor = &discord.EmbedAuthor{
			Name: fmt.Sprintf("%s#%s", author.Username, author.Discriminator),
			Icon: author.AvatarURL(),
		}
	}

	embeds := make([]discord.Embed, 2)
	embeds[0] = discord.Embed{
		Description: desc,
		Color:       color,
		Footer: &discord.EmbedFooter{
			// timestamp will be displayed after the footer text
			Text: fmt.Sprintf("Message ID: %s & sent on", messageID),
		},
		Timestamp: discord.Timestamp(messageID.Time()),
		Author:    eAuthor,
		Fields:    fields,
	}
	embeds[1] = discord.Embed{
		Color: color,
		Footer: &discord.EmbedFooter{
			// timestamp will be displayed after the footer text
			Text: fmt.Sprintf("Channel ID: %s & deleted on", channelID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	return embeds
}
