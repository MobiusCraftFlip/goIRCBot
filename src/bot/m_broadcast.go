package bot

import (
	"fmt"
	"strings"

	"github.com/lrstanley/girc"
)

// m_Broadcast sends a message to every channel the bot is in.
// Usage: .<botnick> broadcast <message>
// Only accepted from the configured admin channel.
type m_Broadcast struct {
	b *Bot
}

func m_broadcast_new() *m_Broadcast {
	return &m_Broadcast{}
}

func (m *m_Broadcast) Name() string { return "broadcast" }

func (m *m_Broadcast) Init(b *Bot) error {
	m.b = b
	b.RegisterCommand("broadcast", m.handleBroadcast)
	return nil
}

func (m *m_Broadcast) handleBroadcast(c *girc.Client, e girc.Event, args []string) {
	// Only allow from admin channel.
	source := e.Params[0]
	adminChan := m.b.Config().AdminChan
	if adminChan == "" || !strings.EqualFold(source, adminChan) {
		return
	}

	if len(args) == 0 {
		c.Cmd.Reply(e, "Usage: broadcast <message>")
		return
	}

	text := strings.Join(args, " ")
	for _, ch := range m.b.IRC().ChannelList() {
		if ch != "" {
			m.b.Privmsg(ch, fmt.Sprintf("!!! Broadcast from %s, %s", e.Source.Name, text))
		}
	}
}

func (m *m_Broadcast) HandleEvent(c *girc.Client, e girc.Event) {}

func (m *m_Broadcast) Shutdown() {}
