package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lrstanley/girc"
)

// m_RegisterChannels registers a channel with ChanServ and applies configured
// flags to the coordinator, operator, and supporter user groups.
//
// Usage: .<botnick> registerchannel <#channel>
// Only accepted from the configured admin channel.
//
// Configuration keys in bot_config:
//
//	registerchannels.coordinators       – comma-separated list of nicks/accounts
//	registerchannels.operators          – comma-separated list of nicks/accounts
//	registerchannels.supporters         – comma-separated list of nicks/accounts
//	registerchannels.coordinators_flags – ChanServ FLAGS string (e.g. "+AFRefiorstv")
//	registerchannels.operators_flags    – ChanServ FLAGS string (e.g. "+AORefiorstv")
//	registerchannels.supporters_flags   – ChanServ FLAGS string (e.g. "+Vv")
type m_RegisterChannels struct {
	b *Bot
}

func m_registerchannels_new() *m_RegisterChannels {
	return &m_RegisterChannels{}
}

func (m *m_RegisterChannels) Name() string { return "registerchannels" }

func (m *m_RegisterChannels) Init(b *Bot) error {
	m.b = b
	b.RegisterCommand("registerchan", m.handleRegister)
	b.RegisterCommand("registermultichan", m.handleRegisterMulti)
	b.RegisterCommand("deregisterchan", m.handleDeregister)
	b.RegisterCommand("deregistermultichan", m.handleDeregisterMulti)
	return nil
}

func (m *m_RegisterChannels) HandleEvent(c *girc.Client, e girc.Event) {}

func (m *m_RegisterChannels) Shutdown() {}

func (m *m_RegisterChannels) handleRegister(c *girc.Client, e girc.Event, args []string) {
	if len(args) == 0 {
		c.Cmd.Reply(e, "Usage: registerchannel <#channel>")
		return
	}
	m.registerChannel(c, args[0])
	c.Cmd.Reply(e, fmt.Sprintf("Registration commands sent for %s.", args[0]))
}

// handleRegisterMulti registers a numbered range of channels.
// Usage: .<botnick> registermultichannels <#base> <start> <end>
// Example: registermultichannels #channel 1 5  →  #channel01 … #channel05
func (m *m_RegisterChannels) handleRegisterMulti(c *girc.Client, e girc.Event, args []string) {
	if len(args) < 3 {
		c.Cmd.Reply(e, "Usage: registermultichannels <#base> <start> <end>")
		return
	}

	base := args[0]
	start, err1 := strconv.Atoi(args[1])
	end, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil || start < 1 || end < start {
		c.Cmd.Reply(e, "start and end must be positive integers with start <= end")
		return
	}

	// Pad width is at least 2, growing to match the digits in end.
	width := max(2, len(strconv.Itoa(end)))
	format := fmt.Sprintf("%s%%0%dd", base, width)

	for i := start; i <= end; i++ {
		m.registerChannel(c, fmt.Sprintf(format, i))
	}

	c.Cmd.Reply(e, fmt.Sprintf("Registration commands sent for %s%0*d through %s%0*d.",
		base, width, start, base, width, end))
}

// registerChannel ensures the bot is in the channel with +o, registers it
// with ChanServ, and applies the configured flags for all user groups.
func (m *m_RegisterChannels) registerChannel(c *girc.Client, channel string) {
	// Ensure the bot is in the channel.
	if c.LookupChannel(channel) == nil {
		c.Cmd.Join(channel)
	}

	// Ensure the bot has +o; request op if not.
	if !m.botIsOp(c, channel) {
		c.Cmd.Mode(channel, "+o", c.GetNick())
	}

	// Register the channel with ChanServ.
	c.Cmd.Message("ChanServ", fmt.Sprintf("REGISTER %s", channel))

	// Apply configured flags for each group.
	extra := m.b.Config().Extra
	m.setGroupFlags(c, channel, extra["registerchannels.coordinators"], extra["registerchannels.coordinators_flags"])
	m.setGroupFlags(c, channel, extra["registerchannels.operators"], extra["registerchannels.operators_flags"])
	m.setGroupFlags(c, channel, extra["registerchannels.supporters"], extra["registerchannels.supporters_flags"])
}

func (m *m_RegisterChannels) handleDeregister(c *girc.Client, e girc.Event, args []string) {
	if len(args) == 0 {
		c.Cmd.Reply(e, "Usage: deregisterchannel <#channel>")
		return
	}
	m.deregisterChannel(c, args[0])
	c.Cmd.Reply(e, fmt.Sprintf("FDROP sent for %s.", args[0]))
}

// handleDeregisterMulti runs FDROP on a numbered range of channels.
// Usage: .<botnick> deregistermultichannels <#base> <start> <end>
func (m *m_RegisterChannels) handleDeregisterMulti(c *girc.Client, e girc.Event, args []string) {
	if len(args) < 3 {
		c.Cmd.Reply(e, "Usage: deregistermultichannels <#base> <start> <end>")
		return
	}

	base := args[0]
	start, err1 := strconv.Atoi(args[1])
	end, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil || start < 1 || end < start {
		c.Cmd.Reply(e, "start and end must be positive integers with start <= end")
		return
	}

	width := max(2, len(strconv.Itoa(end)))
	format := fmt.Sprintf("%s%%0%dd", base, width)

	for i := start; i <= end; i++ {
		m.deregisterChannel(c, fmt.Sprintf(format, i))
	}

	c.Cmd.Reply(e, fmt.Sprintf("FDROP sent for %s%0*d through %s%0*d.",
		base, width, start, base, width, end))
}

// deregisterChannel sends a ChanServ FDROP for the given channel.
func (m *m_RegisterChannels) deregisterChannel(c *girc.Client, channel string) {
	c.Cmd.Message("ChanServ", fmt.Sprintf("FDROP %s", channel))
}

// botIsOp returns true if the bot currently holds +o in the given channel.
func (m *m_RegisterChannels) botIsOp(c *girc.Client, channel string) bool {
	user := c.LookupUser(c.GetNick())
	if user == nil {
		return false
	}
	perms, ok := user.Perms.Lookup(channel)
	return ok && perms.Op
}

// setGroupFlags sends ChanServ FLAGS commands for each nick in the
// comma-separated list, skipping empty entries or a missing flags string.
func (m *m_RegisterChannels) setGroupFlags(c *girc.Client, channel, nicks, flags string) {
	if nicks == "" || flags == "" {
		return
	}
	for nick := range strings.SplitSeq(nicks, ",") {
		nick = strings.TrimSpace(nick)
		if nick == "" {
			continue
		}
		c.Cmd.Message("ChanServ", fmt.Sprintf("FLAGS %s %s %s", channel, nick, flags))
	}
}
