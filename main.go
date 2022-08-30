package main

import (
	"audit/events"
	"audit/logging"
	"context"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

func main() {
	logging.Initialize()

	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))

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

	events.InitEventHandlers(s)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Open(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot open")
	}

	defer func(s *state.State) {
		err := s.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot close")
		}
	}(s)

	<-ctx.Done() // block until Ctrl+C
}
