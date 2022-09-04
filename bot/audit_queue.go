package bot

import (
	"audit/audit"
	"audit/database"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize/english"
	"github.com/gammazero/deque"
	"github.com/rs/zerolog/log"
	"sync"
)

type AuditMessage struct {
	api.SendMessageData

	AuditType audit.Type
	GuildID   discord.GuildID

	attempts uint8

	// errorMessage signifies this message is a wrapper for a failed audit message [currently unimplemented]
	errorMessage *api.SendMessageData
}

var (
	mu           = sync.Mutex{}
	messageQueue = deque.Deque[*AuditMessage]{}
)

// incrementAttempts increments the number of attempts, and returns a boolean for whether the attempt limit was reached.
func (m *AuditMessage) incrementAttempts() bool {
	if m.attempts >= 3 {
		log.Error().Interface("message", m).Msg("Audit message failed to be sent after 3 retries, aborting message")
		return true
	}
	m.attempts += 1
	return false
}

// QueueEmbed creates and queues a new audit message with the given embeds.
func QueueEmbed(auditType audit.Type, guildID discord.GuildID, embeds ...discord.Embed) {
	QueueMessageRaw(AuditMessage{
		AuditType:       auditType,
		GuildID:         guildID,
		SendMessageData: api.SendMessageData{Embeds: embeds},
	})
}

// QueueMessage creates and queues a new audit message with the given message data.
func QueueMessage(auditType audit.Type, guildID discord.GuildID, data api.SendMessageData) {
	QueueMessageRaw(AuditMessage{
		AuditType:       auditType,
		GuildID:         guildID,
		SendMessageData: data,
	})
}

// QueueMessageRaw queues the given audit message.
// Currently, this also starts ProcessMessageQueue as a goroutine after queueing.
func QueueMessageRaw(message AuditMessage) {
	mu.Lock()
	defer mu.Unlock()
	messageQueue.PushBack(&message)

	// TODO do this at a better position
	go ProcessMessageQueue()
}

// ProcessMessageQueue goes through the message queue and attempts to send them to Discord.
// This function essentially wraps ProcessMessage.
func ProcessMessageQueue() {
	mu.Lock()
	defer mu.Unlock()
	var retries []*AuditMessage
	retryMsg := func(msg *AuditMessage) {
		retries = append(retries, msg)
	}

	for messageQueue.Len() != 0 {
		log.Debug().Msgf("Processing message queue: %s to go", english.Plural(messageQueue.Len(), "message", ""))
		msg := messageQueue.PopFront()

		err := ProcessMessage(*msg)
		if err != nil {
			if !msg.incrementAttempts() {
				retryMsg(msg)
			} else {
				// TODO wrap the message and retry as a new message
			}
		}
	}

	for _, message := range retries {
		messageQueue.PushBack(message)
	}
}

// ProcessMessage handles verifying an AuditMessage and attempting to send it to the appropriate log channel.
// If sending the message failed due to an error, it will be returned.
//
// This function does not perform retries or fallbacks; you will have to handle these situations manually based on the
// returned error.
func ProcessMessage(msg AuditMessage) error {
	guild, err := database.Collections.Guilds.GetGuild(msg.GuildID)
	if err != nil {
		// this shouldn't happen normally because the event handlers check for enabled audit types before creating messages
		log.Err(err).
			Interface("guildID", msg.GuildID).
			Interface("auditType", msg.AuditType.String()).
			Msg("Could not process audit message as the guild could not be found in cache")

		return err
	}

	if guild.LoggingDisabled {
		log.Debug().
			Interface("guildID", msg.GuildID).
			Interface("auditType", msg.AuditType.String()).
			Msg("Skipping audit message because logging is disabled")
		return nil
	}

	if !guild.AuditChannelID.IsValid() {
		log.Warn().
			Interface("guildID", msg.GuildID).
			Interface("auditType", msg.AuditType.String()).
			Msg("Could not process audit message as the guild channel ID is not set")
		return nil
	}

	sentMsg, err := s.SendMessageComplex(guild.AuditChannelID, msg.SendMessageData)
	if err != nil {
		log.Err(err).
			Interface("guildID", msg.GuildID).
			Interface("auditType", msg.AuditType.String()).
			Msg("Audit message failed to be sent")

		return err
	}

	log.Debug().
		Interface("msgID", sentMsg.ID).
		Interface("auditType", msg.AuditType.String()).
		Msg("Successfully sent log message")

	return nil
}
