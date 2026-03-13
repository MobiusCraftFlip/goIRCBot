package bot

import (
	"log"
	"regexp"
	"strings"

	"github.com/lrstanley/girc"
)

// m_SwearFilter kicks users whose messages match any configured regex pattern.
//
// Configuration keys in bot_config:
//
//	swearfilter.patterns     – newline-separated list of regex patterns
//	swearfilter.kick_message – reason sent with the kick
type m_SwearFilter struct {
	b        *Bot
	patterns []*regexp.Regexp
}

func m_swearfilter_new() *m_SwearFilter {
	return &m_SwearFilter{}
}

func (m *m_SwearFilter) Name() string { return "swearfilter" }

func (m *m_SwearFilter) Init(b *Bot) error {
	m.b = b
	m.compilePatterns()
	return nil
}

func (m *m_SwearFilter) Shutdown() {}

func (m *m_SwearFilter) HandleEvent(c *girc.Client, e girc.Event) {
	if e.Command != girc.PRIVMSG {
		return
	}

	// Only act on channel messages, not PMs.
	target := e.Params[0]
	if !strings.HasPrefix(target, "#") {
		return
	}

	// Don't filter the bot's own messages.
	if e.Source == nil || e.Source.Name == c.GetNick() {
		return
	}

	// Dont filter moderators:

	message := e.Last()
	for _, re := range m.patterns {
		if re.MatchString(message) {
			kickMsg := m.b.Config().Extra["swearfilter.kick_message"]
			c.Cmd.Kick(target, e.Source.Name, kickMsg)
			return
		}
	}
}

// compilePatterns reads swearfilter.patterns from config and compiles each
// newline-separated entry into a case-insensitive regexp.
func (m *m_SwearFilter) compilePatterns() {
	raw := m.b.Config().Extra["swearfilter.patterns"]
	m.patterns = nil
	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		re, err := regexp.Compile("(?i)" + line)
		if err != nil {
			log.Printf("[swearfilter] invalid pattern %q: %v", line, err)
			continue
		}
		m.patterns = append(m.patterns, re)
	}
	log.Printf("[swearfilter] loaded %d pattern(s)", len(m.patterns))
}
