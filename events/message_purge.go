package events

import (
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
	handler = append(handler, func() {
		s.PreHandler.AddSyncHandler(func(c *gateway.MessageDeleteBulkEvent) {
			if !AuditMessagePurge.check(&c.GuildID, &c.ChannelID) {
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

			go handlePurgedMessages(c, messages, unrecoverableMessages)
		})
	})
}

func handlePurgedMessages(c *gateway.MessageDeleteBulkEvent, msgs []discord.Message, unrecMsgs []discord.MessageID) {
	desc := fmt.Sprintf("**:wastebasket: Messages purged from %s:**\n\nTotal deleted messages: %d", c.ChannelID.Mention(), len(c.IDs))
	embed := discord.Embed{
		Description: desc,
		Color:       0xFF0000,
		Footer:      &discord.EmbedFooter{Text: fmt.Sprintf("Channel ID: %s", c.ChannelID.String())},
		Timestamp:   discord.NowTimestamp(),
	}

	files := make([]sendpart.File, 0)
	if len(msgs) != 0 {
		// todo: create a separate file with human-readable formatting for purged messages
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

	// todo: use the archive system to save the raw message data and hyperlink to it
	handleAuditError(s.SendMessageComplex(auditChannel, api.SendMessageData{
		Embeds: []discord.Embed{embed},
		Files:  files,
	}))
}
