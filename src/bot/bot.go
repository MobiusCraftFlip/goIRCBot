package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lrstanley/girc"

	"scoutdev.org/m/v2/goIrcBot/src/db"
)

// Module is the interface every bot module must implement.
type Module interface {
	// Name returns the unique identifier for the module.
	Name() string
	// Init is called once when the bot starts, before joining channels.
	Init(b *Bot) error
	// HandleEvent is called for every IRC event received.
	HandleEvent(c *girc.Client, e girc.Event)
	// Shutdown is called when the bot is stopping.
	Shutdown()
}

// Bot is the central coordinator: IRC connection, database, and modules.
type Bot struct {
	cfg      *db.Config
	db       *db.DB
	client   *Client
	modules  []Module
	commands map[string]func(*girc.Client, girc.Event, []string)
}

// New creates a Bot with the given config and database handle.
// Modules listed in cfg.Modules are instantiated from the factory registry.
func New(cfg *db.Config, database *db.DB) *Bot {
	b := &Bot{cfg: cfg, db: database, commands: make(map[string]func(*girc.Client, girc.Event, []string))}

	for _, name := range cfg.Modules {
		switch name {
		case "core":
			b.Register(m_core_new())
		case "pmproxy":
			b.Register(m_pmproxy_new())
		case "broadcast":
			b.Register(m_broadcast_new())
		case "channelcontrol":
			b.Register(m_channelcontrol_new())
		case "registerchannels":
			b.Register(m_registerchannels_new())
		case "swearfilter":
			b.Register(m_swearfilter_new())
		default:
			log.Printf("[bot] unknown module in config: %s", name)
		}
	}

	return b
}

// Register adds a module. Must be called before Run.
func (b *Bot) Register(m Module) {
	b.modules = append(b.modules, m)
}

func (b *Bot) RegisterCommand(cmd string, handler func(*girc.Client, girc.Event, []string)) {
	b.commands[cmd] = handler
}

// Config returns the bot's current configuration.
func (b *Bot) Config() *db.Config {
	return b.cfg
}

// DB returns the database handle.
func (b *Bot) DB() *db.DB {
	return b.db
}

// Privmsg sends a PRIVMSG to target.
func (b *Bot) Privmsg(target, text string) {
	if b.client != nil {
		b.client.Privmsg(target, text)
	}
}

// Notice sends a NOTICE to target.
func (b *Bot) Notice(target, text string) {
	if b.client != nil {
		b.client.Notice(target, text)
	}
}

// IRC returns the underlying girc.Client for direct access when needed.
func (b *Bot) IRC() *girc.Client {
	if b.client == nil {
		return nil
	}
	return b.client.Inner()
}

// Run initialises modules then connects, reconnecting automatically on errors.
func (b *Bot) Run() error {
	for _, m := range b.modules {
		if err := m.Init(b); err != nil {
			return err
		}
		log.Printf("[bot] module %q initialised", m.Name())
	}

	for {
		if err := b.connect(); err != nil {
			log.Printf("[bot] connection error: %v — retrying in 30s", err)
			time.Sleep(30 * time.Second)
			continue
		}
		break
	}

	for _, m := range b.modules {
		m.Shutdown()
	}
	return nil
}

func (b *Bot) connect() error {
	client, err := newClient(b.cfg, b.dispatch)
	if err != nil {
		return err
	}
	b.client = client
	return client.Run()
}

func (b *Bot) HandleCommand(c *girc.Client, e girc.Event) {
	if !b.isModerator(e) {
		return
	}

	message := e.Params[len(e.Params)-1]
	fields := strings.Fields(message)
	if strings.ToLower(fields[0]) == "."+strings.ToLower(b.cfg.Nick) {
		commandName := strings.ToLower(fields[1])

		cmd, exists := b.commands[commandName]
		if exists {
			cmd(c, e, fields[2:])
		} else {
			c.Cmd.Reply(e, fmt.Sprintf("Unknown command: %s", commandName))
		}
	}
}

// dispatch fans each IRC event out to every registered module.
func (b *Bot) dispatch(c *girc.Client, e girc.Event) {
	if e.Command == girc.PRIVMSG {
		b.HandleCommand(c, e)
	}

	log.Println(e.String())
	for _, m := range b.modules {
		m.HandleEvent(c, e)
	}
}
