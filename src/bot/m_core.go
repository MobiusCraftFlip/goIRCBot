package bot

import (
	"log"

	"github.com/lrstanley/girc"
)

// Module handles m_Core IRC housekeeping: gaining oper status and ensuring the
// bot is present in the admin and log channels.
type m_Core struct {
	b *Bot
}

// New returns a new m_Core Module.
func m_core_new() *m_Core {
	return &m_Core{}
}

func (m *m_Core) Name() string { return "m_Core" }

func (m *m_Core) Init(b *Bot) error {
	m.b = b
	return nil
}

// HandleEvent watches for CONNECTED so we can send OPER and join the
// admin/log channels immediately after the bot registers on the server.
func (m *m_Core) HandleEvent(c *girc.Client, e girc.Event) {
	if e.Command != girc.CONNECTED {
		return
	}

	cfg := m.b.Config()
	nsUser := cfg.Extra["ns_user"]
	nsPass := cfg.Extra["ns_pass"]

	if nsUser != "" && nsPass != "" {
		c.Cmd.Message("NickServ", "LOGIN "+nsUser+" "+nsPass)
		log.Printf("[m_Core] sent NickServ LOGIN as %s", nsUser)
	}
	log.SetOutput(m.b)
	// Gain oper privileges if credentials are configured.
	operUser := cfg.Extra["oper_user"]
	operPass := cfg.Extra["oper_pass"]
	log.Printf("[m_Core] might OPER as %s", operUser)
	log.Println(cfg.Extra)
	if operUser != "" && operPass != "" {
		c.Cmd.Oper(operUser, operPass)
		log.Printf("[m_Core] sent OPER as %s", operUser)
	}

	// Ensure the bot is in the admin and log channels.
	for _, ch := range []string{cfg.AdminChan, cfg.LogChan} {
		if ch != "" {
			c.Cmd.Join(ch)
		}
	}
}

func (m *m_Core) Shutdown() {}
