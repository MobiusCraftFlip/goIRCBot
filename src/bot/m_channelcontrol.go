package bot

import (
	"fmt"
	"strconv"

	"github.com/lrstanley/girc"
)

type m_ChannelControl struct {
	b *Bot
}

func m_channelcontrol_new() *m_ChannelControl {
	return &m_ChannelControl{}
}

func (m *m_ChannelControl) Name() string { return "channelcontrol" }

func (m *m_ChannelControl) Init(b *Bot) error {
	m.b = b
	b.RegisterCommand("openchan", m.handleOpen)
	b.RegisterCommand("closechan", m.handleClose)
	b.RegisterCommand("openmultichans", m.handleOpenMulti)
	b.RegisterCommand("closemultichans", m.handleCloseMulti)
	return nil
}

func (m *m_ChannelControl) HandleEvent(c *girc.Client, e girc.Event) {}

func (m *m_ChannelControl) Shutdown() {}

// handleOpen applies the open topic and modes to the given channel.
// Usage: .<botnick> openchan <#channel>
func (m *m_ChannelControl) handleOpen(c *girc.Client, e girc.Event, args []string) {
	if len(args) == 0 {
		c.Cmd.Reply(e, "Usage: openchan <#channel>")
		return
	}
	if ok := m.applyState(c, args[0], "open"); ok {
		c.Cmd.Reply(e, fmt.Sprintf("Channel %s is now open.", args[0]))
	} else {
		c.Cmd.Reply(e, "No open configuration found.")
	}
}

// handleClose applies the close topic and modes to the given channel.
// Usage: .<botnick> closechan <#channel>
func (m *m_ChannelControl) handleClose(c *girc.Client, e girc.Event, args []string) {
	if len(args) == 0 {
		c.Cmd.Reply(e, "Usage: closechan <#channel>")
		return
	}
	if ok := m.applyState(c, args[0], "close"); ok {
		c.Cmd.Reply(e, fmt.Sprintf("Channel %s is now closed.", args[0]))
	} else {
		c.Cmd.Reply(e, "No close configuration found.")
	}
}

// handleOpenMulti opens a numbered range of channels.
// Usage: .<botnick> openmultichans <#base> <start> <end>
func (m *m_ChannelControl) handleOpenMulti(c *girc.Client, e girc.Event, args []string) {
	base, start, end, ok := m.parseMultiArgs(c, e, args, "openmultichans")
	if !ok {
		return
	}
	width, format := multiFormat(base, end)
	for i := start; i <= end; i++ {
		m.applyState(c, fmt.Sprintf(format, i), "open")
	}
	c.Cmd.Reply(e, fmt.Sprintf("Opened %s%0*d through %s%0*d.", base, width, start, base, width, end))
}

// handleCloseMulti closes a numbered range of channels.
// Usage: .<botnick> closemultichans <#base> <start> <end>
func (m *m_ChannelControl) handleCloseMulti(c *girc.Client, e girc.Event, args []string) {
	base, start, end, ok := m.parseMultiArgs(c, e, args, "closemultichans")
	if !ok {
		return
	}
	width, format := multiFormat(base, end)
	for i := start; i <= end; i++ {
		m.applyState(c, fmt.Sprintf(format, i), "close")
	}
	c.Cmd.Reply(e, fmt.Sprintf("Closed %s%0*d through %s%0*d.", base, width, start, base, width, end))
}

// applyState sets the topic and modes for state ("open" or "close").
// Returns false if no configuration exists for that state.
// Topic is a global config key: channelcontrol.<state>_topic
// Modes are a global config key: channelcontrol.<state>_modes
func (m *m_ChannelControl) applyState(c *girc.Client, channel, state string) bool {
	extra := m.b.Config().Extra

	topic, hasTopic := extra[fmt.Sprintf("channelcontrol.%s_topic", state)]
	modes, hasModes := extra[fmt.Sprintf("channelcontrol.%s_modes", state)]

	if !hasTopic && !hasModes {
		return false
	}

	if hasModes && modes != "" {
		c.Cmd.Mode(channel, modes)
	}
	if hasTopic {
		c.Cmd.Topic(channel, topic)
	}
	return true
}

func (m *m_ChannelControl) parseMultiArgs(c *girc.Client, e girc.Event, args []string, cmd string) (base string, start, end int, ok bool) {
	if len(args) < 3 {
		c.Cmd.Reply(e, fmt.Sprintf("Usage: %s <#base> <start> <end>", cmd))
		return
	}
	var err1, err2 error
	base = args[0]
	start, err1 = strconv.Atoi(args[1])
	end, err2 = strconv.Atoi(args[2])
	if err1 != nil || err2 != nil || start < 1 || end < start {
		c.Cmd.Reply(e, "start and end must be positive integers with start <= end")
		return
	}
	ok = true
	return
}

// multiFormat returns the zero-pad width and fmt format string for a numbered
// channel range where numbers are padded to at least 2 digits.
func multiFormat(base string, end int) (width int, format string) {
	width = max(2, len(strconv.Itoa(end)))
	format = fmt.Sprintf("%s%%0%dd", base, width)
	return
}
