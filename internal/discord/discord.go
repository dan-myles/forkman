package discord

import (
	"errors"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/avvo-na/forkman/internal/discord/qna"
	"github.com/avvo-na/forkman/internal/discord/verification"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Discord struct {
	session      *discordgo.Session
	db           *gorm.DB
	log          *zerolog.Logger
	cfg          *config.ForkConfig
	email        *ses.Client
	bedrock      *bedrockagentruntime.Client
	moderation   map[string]*moderation.Moderation     /* GuildID -> module*/
	verification map[string]*verification.Verification /* GuildID -> module */
	qna          map[string]*qna.QNA                   /* GuildID -> module */
}

var ErrModuleNotFound = errors.New("module not found")

func New(cfg *config.ForkConfig, log *zerolog.Logger, db *gorm.DB, acfg aws.Config) *Discord {
	d := &Discord{
		db:      db,
		log:     log,
		cfg:     cfg,
		email:   ses.NewFromConfig(acfg),
		bedrock: bedrockagentruntime.NewFromConfig(acfg),
	}

	s, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		panic(err)
	}

	// Settings
	s.Identify.Intents = discordgo.IntentsAll // What do we need permission for?
	s.SyncEvents = false                      // Launch goroutines for handlers
	s.StateEnabled = true
	d.session = s

	// Module stores
	d.moderation = make(map[string]*moderation.Moderation)
	d.verification = make(map[string]*verification.Verification)
	d.qna = make(map[string]*qna.QNA)

	// Global handlers
	s.AddHandler(d.onReadyNotify)
	s.AddHandler(d.onGuildCreateGuildUpdate)
	s.AddHandler(d.onInteractionCreate)
	s.AddHandler(d.onMessageCreate)

	// Open the session
	log.Info().Msg("Opening discord session")
	err = s.Open()
	if err != nil {
		panic(err)
	}

	return d
}

func (d *Discord) Open() error {
	err := d.session.Open()
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) Close() error {
	err := d.session.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) GetSession() *discordgo.Session {
	return d.session
}

func (d *Discord) GetQNAModule(guildSnowflake string) (*qna.QNA, error) {
	mod, ok := d.qna[guildSnowflake]
	if !ok {
		return nil, ErrModuleNotFound
	}

	return mod, nil
}

func (d *Discord) GetModerationModule(guildSnowflake string) (*moderation.Moderation, error) {
	mod, ok := d.moderation[guildSnowflake]
	if !ok {
		return nil, ErrModuleNotFound
	}

	return mod, nil
}

func (d *Discord) GetVerificationModule(guildSnowflake string) (*verification.Verification, error) {
	mod, ok := d.verification[guildSnowflake]
	if !ok {
		return nil, ErrModuleNotFound
	}

	return mod, nil
}

func (d *Discord) GetUserAdminServers(userID string) ([]discordgo.UserGuild, error) {
	// loop through all guilds and check if that user is an admin
	// if so, add the guild to the list
	ret := []discordgo.UserGuild{}
	guilds, err := d.session.UserGuilds(200, "", "", false)
	if err != nil {
		return nil, err
	}

	for _, guild := range guilds {
		member, err := d.session.GuildMember(guild.ID, userID)
		if err != nil {
			return nil, err
		}

		roles, err := d.session.GuildRoles(guild.ID)
		if err != nil {
			return nil, err
		}

		for _, userRole := range member.Roles {
			for _, guildRole := range roles {
				if userRole == guildRole.ID {
					if guildRole.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
						ret = append(ret, *guild)
					}
				}
			}
		}
	}

	return ret, nil
}

func (d *Discord) onReadyNotify(s *discordgo.Session, r *discordgo.Ready) {
	d.log.Info().Msgf("Session has connected to Discord as %s", r.User.String())
}

func (d *Discord) onGuildCreateGuildUpdate(s *discordgo.Session, g *discordgo.GuildCreate) {
	log := d.log.With().Str("event", "onGuildCreate").
		Str("guild_snowflake", g.Guild.ID).
		Str("guild_name", g.Guild.Name).
		Logger()

		// Create guild repo
	repo := database.NewGuildRepository(d.db)

	// Read or create guild
	_, err := repo.ReadGuild(g.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("critical error reading guild")
		return
	}

	if err == gorm.ErrRecordNotFound {
		if _, err := repo.CreateGuild(g.Guild); err != nil {
			log.Error().Err(err).Msg("critical error creating guild")
			return
		}
		log.Info().Msg("guild creation complete")
		return
	}

	_, err = repo.UpdateGuild(g.Guild)
	if err != nil {
		log.Error().Err(err).Msg("critical error updating guild")
		return
	}

	m := moderation.New(g.Name, g.ID, d.cfg.DiscordAppID, d.session, d.db, d.log)
	if err := m.Load(); err != nil {
		log.Error().Err(err).Msg("critical error init moderation module")
		return
	}

	v := verification.New(g.Name, g.ID, d.cfg.DiscordAppID, d.session, d.db, d.email, d.log)
	if err := v.Load(); err != nil {
		log.Error().Err(err).Msg("critical error init verification module")
		return
	}

	q := qna.New(g.Name, g.ID, d.cfg.DiscordAppID, d.session, d.bedrock, d.cfg.FORUM_CHANNEL_ID, d.cfg.AWS_BEDROCK_KBI, d.db, d.log)
	if err := q.Load(); err != nil {
		log.Error().Err(err).Msg("critical error init qna module")
		return
	}

	d.moderation[g.ID] = m
	d.verification[g.ID] = v
	d.qna[g.ID] = q

	log.Debug().Msg("guild instantiation complete")
}

func (d *Discord) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	go d.moderation[i.GuildID].OnInteractionCreate(s, i)
	go d.verification[i.GuildID].OnInteractionCreate(s, i)
	go d.qna[i.GuildID].OnInteractionCreate(s, i)

	log := d.log.With().
		Str("guild_id", i.GuildID).
		Str("channel_id", i.ChannelID).
		Str("interaction_id", i.ID).
		Str("interaction_type", i.Type.String()).
		Str("user_id", i.Member.User.ID).
		Str("user_name", i.Member.User.GlobalName).
		Str("guild_locale", i.GuildLocale.String()).
		Int("version", i.Version).
		Logger()

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		log = log.With().
			Str("command_name", i.ApplicationCommandData().Name).
			Str("command_id", i.ApplicationCommandData().ID).
			Str("command_type", GetApplicationCommandType(i.ApplicationCommandData().CommandType)).
			Logger()
	case discordgo.InteractionModalSubmit:
		log = log.With().
			Str("modal_type", i.ModalSubmitData().Type().String()).
			Str("custom_id", i.ModalSubmitData().CustomID).
			Logger()
	case discordgo.InteractionMessageComponent:
		log = log.With().
			Str("custom_id", i.MessageComponentData().CustomID).
			Logger()
	default:
		log.Error().Msg("critical: unhandled application interaction")
		return
	}

	guild, _ := d.session.Guild(i.GuildID)
	log = log.With().Str("guild_name", guild.Name).Logger()

	log.Info().Msg("interaction request received")
}

func (d *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go d.qna[m.GuildID].OnMessageCreate(s, m)
}
