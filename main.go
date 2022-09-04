package main

import (
	"audit/bot"
	"audit/database"
	"audit/events"
	"audit/logging"
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	logging.Initialize()
	_, err := database.Initialize(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Failed connecting to MongoDB")
	}
	defer func(ctx context.Context) {
		err := database.Disconnect(ctx)
		if err != nil {
			log.Panic().Err(err).Msg("Could not disconnect from MongoDB")
		}
	}(ctx)

	s, err := bot.Initialize(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Failed opening Discord bot")
	}
	defer func() {
		err := bot.Close()
		if err != nil {
			log.Panic().Err(err).Msg("Could not close Discord bot")
		}
	}()

	events.InitEventHandlers(s)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.Info().Msg("Bot is ready")

	<-ctx.Done() // block until Ctrl+C
}
