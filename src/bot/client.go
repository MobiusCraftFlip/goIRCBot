package bot

import (
	"crypto/tls"
	"fmt"

	"github.com/lrstanley/girc"

	"scoutdev.org/m/v2/goIrcBot/src/db"
)

// Client wraps girc.Client and manages the IRC connection lifecycle.
type Client struct {
	inner *girc.Client
}

// newClient creates a Client from bot config and wires the dispatch handler.
func newClient(cfg *db.Config, dispatch func(*girc.Client, girc.Event)) (*Client, error) {
	gircCfg := girc.Config{
		Server:     cfg.Server,
		Port:       cfg.Port,
		Nick:       cfg.Nick,
		User:       cfg.Username,
		Name:       cfg.Realname,
		AllowFlood: true,
	}

	if cfg.SSL {
		gircCfg.SSL = true
		gircCfg.TLSConfig = &tls.Config{ServerName: cfg.Server}
	}

	inner := girc.New(gircCfg)

	// Join configured channels once the bot is registered on the server.
	inner.Handlers.AddBg(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		for _, ch := range cfg.Channels {
			if ch != "" {
				c.Cmd.Join(ch)
			}
		}
	})

	// Fan every event out to the module dispatcher.
	inner.Handlers.Add(girc.ALL_EVENTS, dispatch)

	return &Client{inner: inner}, nil
}

// Run connects to the IRC server and blocks until the connection is closed.
func (c *Client) Run() error {
	if err := c.inner.Connect(); err != nil {
		return fmt.Errorf("irc connect: %w", err)
	}
	return nil
}

// Privmsg sends a PRIVMSG to target.
func (c *Client) Privmsg(target, text string) {
	c.inner.Cmd.Message(target, text)
}

// Notice sends a NOTICE to target.
func (c *Client) Notice(target, text string) {
	c.inner.Cmd.Notice(target, text)
}

// Inner returns the underlying girc.Client for direct access.
func (c *Client) Inner() *girc.Client {
	return c.inner
}
