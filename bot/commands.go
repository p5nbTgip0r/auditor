package bot

import (
	"audit/audit"
	"audit/commands"
	"audit/database"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
)

//var (
//	cmds = []commands.Command{
//		*commands.BuildCommand(
//			"enableaudittype",
//			"Enables the specified audit types for this guild",
//			discord.NewStringOption("audit types", "space-separated list of audit types", true),
//		),
//	}
//)

type commandHandlerPair struct {
	data    api.CreateCommandData
	handler func(*commands.Interaction)
}

func createCommandData(handlers []commandHandlerPair) []api.CreateCommandData {
	d := make([]api.CreateCommandData, len(handlers))
	for i, handler := range handlers {
		d[i] = handler.data
	}
	return d
}

func handlerMap(handlers []commandHandlerPair) map[string]func(*commands.Interaction) {
	m := make(map[string]func(*commands.Interaction), len(handlers))
	for _, handler := range handlers {
		m[handler.data.Name] = handler.handler
	}
	return m
}

func getCmds() []commandHandlerPair {
	parseTypes := func(types string) mapset.Set[audit.Type] {
		set := mapset.NewSet[audit.Type]()
		split := strings.Split(types, " ")
		for _, t := range split {
			auditType := audit.UnString(t)
			if auditType == audit.Unknown {
				continue
			}
			set.Add(auditType)
		}
		return set
	}
	toggleAuditTypes := func(c *commands.Interaction, enable bool) {
		_ = c.DeferResponse(true)
		opt := c.CommandInteraction.Options.Find("audit_types")
		guild, _ := database.Collections.Guilds.GetGuild(c.GuildID)
		types := opt.String()
		parsedTypes := parseTypes(types)
		guild.ChangeAuditTypes(parsedTypes, enable)

		err := database.Collections.Guilds.SetGuild(c.GuildID, *guild)
		if err != nil {
			log.Err(err).Msg("Failed to update guild")
			_, _ = c.EditTextResponse("Failed to update audit types")
			return
		}
		var builder strings.Builder
		for _, t := range guild.DisabledAuditTypes {
			if builder.Len() != 0 {
				builder.WriteString("\n")
			}
			builder.WriteString("- ")
			builder.WriteString(t.String())
		}
		if builder.Len() == 0 {
			builder.WriteString("-no disabled audit types-")
		}

		_, _ = c.EditTextResponse(fmt.Sprintf("New list of disabled audit types:\n```%s```", builder.String()))
	}
	adminPerm := discord.Permissions(0)

	auditTypes := discord.NewStringOption("audit_types", "space-separated list of audit types", true)
	tempCmds := []commandHandlerPair{
		{
			data: api.CreateCommandData{
				Name:                     "enableaudittype",
				Description:              "Enable the specified audit types",
				NoDMPermission:           true,
				DefaultMemberPermissions: &adminPerm,
				Options:                  []discord.CommandOption{auditTypes},
			},
			handler: func(c *commands.Interaction) {
				toggleAuditTypes(c, true)
			},
		},
		{
			data: api.CreateCommandData{
				Name:                     "disableaudittype",
				Description:              "Disable the specified audit types",
				NoDMPermission:           true,
				DefaultMemberPermissions: &adminPerm,
				Options:                  []discord.CommandOption{auditTypes},
			},
			handler: func(c *commands.Interaction) {
				toggleAuditTypes(c, false)
			},
		},
		{
			data: api.CreateCommandData{
				Name:                     "auditchannel",
				Description:              "Set the channel where audit logs are sent",
				NoDMPermission:           true,
				DefaultMemberPermissions: &adminPerm,
				Options: []discord.CommandOption{
					discord.NewChannelOption("channel", "Audit log channel", true),
				},
			},
			handler: func(c *commands.Interaction) {
				_ = c.DeferResponse(true)
				channelID := discord.NullChannelID
				opt := c.CommandInteraction.Options.Find("channel")
				guild, _ := database.Collections.Guilds.GetGuild(c.GuildID)
				v, _ := opt.SnowflakeValue()
				channelID = discord.ChannelID(v)
				guild.AuditChannelID = channelID

				err := database.Collections.Guilds.SetGuild(c.GuildID, *guild)
				if err != nil {
					log.Err(err).Msg("Failed to update guild")
					_, _ = c.EditTextResponse("Failed to update audit channel")
					return
				}

				_, _ = c.EditTextResponse(fmt.Sprintf("Audit channel was set to %s", channelID.Mention()))
			},
		},
		{
			data: api.CreateCommandData{
				Name:                     "logging",
				Description:              "Toggle logging on or off",
				NoDMPermission:           true,
				DefaultMemberPermissions: &adminPerm,
				Options: []discord.CommandOption{
					discord.NewBooleanOption("enabled", "Logging enabled status", false),
				},
			},
			handler: func(c *commands.Interaction) {
				_ = c.DeferResponse(true)
				opt := c.CommandInteraction.Options.Find("enabled")
				guild, _ := database.Collections.Guilds.GetGuild(c.GuildID)
				if opt.Name != "" {
					e, _ := opt.BoolValue()
					guild.LoggingDisabled = !e
				} else {
					guild.LoggingDisabled = !guild.LoggingDisabled
				}
				err := database.Collections.Guilds.SetGuild(c.GuildID, *guild)
				if err != nil {
					log.Err(err).Msg("Failed to update guild")
					_, _ = s.EditInteractionResponse(c.AppID, c.Token, api.EditInteractionResponseData{
						Content: option.NewNullableString("Failed to update logging status"),
					})
					return
				}

				_, _ = s.EditInteractionResponse(c.AppID, c.Token, api.EditInteractionResponseData{
					Content: option.NewNullableString(fmt.Sprintf("Logging was set to %t", !guild.LoggingDisabled)),
				})
			},
		},
	}

	return tempCmds
}

func updateCommands() error {
	app, _ := s.CurrentApplication()
	commandData := getCmds()
	datas := createCommandData(commandData)
	handlers := handlerMap(commandData)

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG_ENABLED"))
	guildIdEnv := os.Getenv("DEBUG_GUILD_ID")
	var guildID discord.GuildID

	if atoi, err := strconv.Atoi(guildIdEnv); err == nil {
		guildID = discord.GuildID(atoi)
	}

	var cmds []discord.Command
	var err error
	if debug {
		cmds, err = s.BulkOverwriteGuildCommands(app.ID, guildID, datas)
	} else {
		cmds, err = s.BulkOverwriteCommands(app.ID, datas)
		if err == nil && guildID.IsValid() {
			// delete all the debug guild commands
			guildCmds, err := s.BulkOverwriteGuildCommands(app.ID, guildID, []api.CreateCommandData{})
			if err != nil {
				// this step isn't crucial to the bot functioning, so we just warn about the error
				log.Warn().Err(err).Msg("Could not delete guild commands for non-debug environment")
			} else {
				log.Debug().Interface("commands", guildCmds).Msg("Deleted guild commands successfully")
			}
		}
	}
	if err != nil {
		return err
	}

	log.Debug().Interface("commands", cmds).Msg("Updated commands successfully")

	s.Handler.AddHandler(func(c *gateway.InteractionCreateEvent) {
		switch data := c.Data.(type) {
		case *discord.CommandInteraction:
			h, ok := handlers[data.Name]
			if ok {
				log.Debug().Interface("interaction_data", data).Msg("Received known interaction")
				interaction := &commands.Interaction{
					InteractionCreateEvent: c,
					CommandInteraction:     data,
					State:                  s,
				}
				h(interaction)
			} else {
				log.Warn().Interface("interaction", data).Msg("Received unknown interaction")
			}
		}
	})

	return nil
}
