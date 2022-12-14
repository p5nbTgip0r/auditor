package events

import (
	"audit/audit"
	"audit/bot"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/rs/zerolog/log"
)

func init() {
	handle := func(c gateway.MessageDeleteBulkEvent, msgs []discord.Message, unrecMsgs []discord.MessageID) {
		if !check(audit.MessagePurge, &c.GuildID, &c.ChannelID) {
			return
		}
		desc := fmt.Sprintf("**:wastebasket: Messages purged from %s:**\n\nTotal deleted messages: %d", c.ChannelID.Mention(), len(c.IDs))
		embed := discord.Embed{
			Description: desc,
			Color:       0xFF0000,
			Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Channel ID: %s", c.ChannelID.String())},
			Timestamp:   discord.NowTimestamp(),
		}

		files := make([]sendpart.File, 0)
		if len(msgs) != 0 {
			// TODO create a separate file with human-readable formatting for purged messages
			msgsJson, err := json.Marshal(msgs)
			if err == nil {
				files = append(files, sendpart.File{Name: "messages.json", Reader: bytes.NewReader(msgsJson)})
			}
		}
		if len(unrecMsgs) != 0 {
			desc = desc + fmt.Sprintf("\n\nTotal unrecoverable messages: %d", len(unrecMsgs))

			log.Warn().
				Interface("event", c).
				Interface("messageIds", unrecMsgs).
				Msgf("Could not recover some message IDs from bulk message delete")

			unrecMsgsJson, err := json.Marshal(unrecMsgs)
			if err == nil {
				files = append(files, sendpart.File{Name: "unrecoverable_messages.json", Reader: bytes.NewReader(unrecMsgsJson)})
			}
		}

		// TODO use the archive system to save the raw message data and hyperlink to it
		err := bot.ProcessMessage(bot.AuditMessage{
			AuditType: audit.MessageDelete,
			GuildID:   c.GuildID,
			SendMessageData: api.SendMessageData{
				Embeds: []discord.Embed{embed},
				Files:  files,
			},
		})
		if err == nil {
			return
		}
		// fallback to the normal message handling if sending attachments failed
		bot.QueueEmbed(audit.MessageDelete, c.GuildID, embed)
	}

	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.MessageDeleteBulkEvent) {
			if !c.GuildID.IsValid() {
				return
			}
			var unrecoverableMessages []discord.MessageID
			var messages []discord.Message
			for _, mID := range c.IDs {
				msg, _ := s.Cabinet.Message(c.ChannelID, mID)
				if msg == nil {
					unrecoverableMessages = append(unrecoverableMessages, mID)
				} else {
					messages = append(messages, *msg)
				}
			}

			go handle(*c, messages, unrecoverableMessages)
		})
	})
}
