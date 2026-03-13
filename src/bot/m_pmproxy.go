package bot

import (
	"fmt"
	"log"

	"github.com/lrstanley/girc"
)

// m_PmProxy forwards private messages sent to the bot to the admin channel.
type m_PmProxy struct {
	b *Bot
}

func m_pmproxy_new() *m_PmProxy {
	return &m_PmProxy{}
}

func (m *m_PmProxy) Name() string { return "pmproxy" }

func (m *m_PmProxy) Init(b *Bot) error {
	m.b = b
	return nil
}

// HandleEvent listens for PRIVMSG events addressed directly to the bot and
// relays them to the configured admin channel.
func (m *m_PmProxy) HandleEvent(c *girc.Client, e girc.Event) {

	if e.Command != girc.PRIVMSG {
		return
	}

	target := e.Params[0]

	// Ignore messages sent to channels; we only want PMs.
	if target != c.GetNick() {
		return
	}

	adminChan := m.b.Config().AdminChan
	if adminChan == "" {
		return
	}

	sender := e.Source.Name
	text := e.Last()
	log.Println("[PM from %s], %s", sender, text, adminChan)
	c.Cmd.Message(adminChan, fmt.Sprintf("[PM from %s] %s", sender, text))
}

func (m *m_PmProxy) Shutdown() {}
