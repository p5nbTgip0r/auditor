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

	AuditType audit.AuditType
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
func QueueEmbed(auditType audit.AuditType, guildID discord.GuildID, embeds ...discord.Embed) {
	QueueMessageRaw(AuditMessage{
		AuditType:       auditType,
		GuildID:         guildID,
		SendMessageData: api.SendMessageData{Embeds: embeds},
	})
}

// QueueMessage creates and queues a new audit message with the given message data.
func QueueMessage(auditType audit.AuditType, guildID discord.GuildID, data api.SendMessageData) {
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

		guild, err := database.Collections.Guilds.GetGuild(msg.GuildID)
		if err != nil {
			// this shouldn't happen normally because the event handlers check for enabled audit types before creating messages
			log.Err(err).
				Interface("guildID", msg.GuildID).
				Interface("auditType", msg.AuditType.String()).
				Msg("Could not process audit message as the guild could not be found in cache")

			if !msg.incrementAttempts() {
				retryMsg(msg)
			} else {
				// TODO wrap the message and retry as a new message
			}
			continue
		}

		if guild.LoggingDisabled {
			log.Debug().
				Interface("guildID", msg.GuildID).
				Interface("auditType", msg.AuditType.String()).
				Msg("Skipping audit message because logging is disabled")
			continue
		}

		if !guild.AuditChannelID.IsValid() {
			log.Warn().
				Interface("guildID", msg.GuildID).
				Interface("auditType", msg.AuditType.String()).
				Msg("Could not process audit message as the guild channel ID is not set")
			continue
		}

		sentMsg, err := s.SendMessageComplex(guild.AuditChannelID, msg.SendMessageData)
		if err != nil {
			log.Err(err).
				Interface("guildID", msg.GuildID).
				Interface("auditType", msg.AuditType.String()).
				Msg("Audit message failed to be sent")

			if !msg.incrementAttempts() {
				retryMsg(msg)
			} else {
				// TODO wrap the message and retry as a new message
			}
			continue
		}

		log.Debug().
			Interface("msgID", sentMsg.ID).
			Interface("auditType", msg.AuditType.String()).
			Msg("Successfully sent log message")
	}

	for _, message := range retries {
		messageQueue.PushBack(message)
	}
}