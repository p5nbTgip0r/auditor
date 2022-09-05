package bot

import (
	"audit/util"
	"context"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	s *state.State
)

func Initialize() (*state.State, error) {
	log.Debug().Msg("Initialize Discord bot..")

	token, err := util.GetEnvRequired("DISCORD_TOKEN")
	if err != nil {
		return nil, err
	}

	maxMsgs := util.GetEnvIntDefault("DISCORD_MAX_MESSAGES", 1000)
	store := defaultstore.New()
	store.MessageStore = defaultstore.NewMessage(maxMsgs)

	s = state.NewWithStore("Bot "+token, store)

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

	return s, err
}

func Open(ctx context.Context) error {
	log.Debug().Msg("Connecting to Discord..")
	err := s.Open(ctx)
	if err == nil {
		log.Debug().Msg("Connected to Discord")
		err = errors.Wrap(updateCommands(), "update slash commands")
	}
	return err
}

func Close() error {
	log.Debug().Msg("Bot disconnect..")
	return s.Close()
}
