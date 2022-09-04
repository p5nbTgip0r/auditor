package bot

import (
	"context"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	s *state.State
)

func Initialize(ctx context.Context) (*state.State, error) {
	log.Debug().Msg("Initialize Discord bot..")
	s = state.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentGuildMembers)
	s.AddIntents(gateway.IntentGuildPresences)
	s.AddIntents(gateway.IntentGuildVoiceStates)
	s.AddIntents(gateway.IntentGuildBans)
	s.AddIntents(gateway.IntentGuildEmojis)
	s.AddIntents(gateway.IntentGuildInvites)
	s.AddIntents(gateway.IntentGuildWebhooks)
	// Make a pre-handler
	s.PreHandler = handler.New()

	log.Debug().Msg("Connecting to Discord..")
	err := s.Open(ctx)
	if err == nil {
		log.Debug().Msg("Connected to Discord")
	}
	return s, err
}

func Close() error {
	log.Debug().Msg("Bot disconnect..")
	return s.Close()
}
