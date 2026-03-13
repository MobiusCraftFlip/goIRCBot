package db

import (
	"context"
	"fmt"
)

// Config holds all settings for a bot instance.
type Config struct {
	BotID     int
	Nick      string
	Username  string
	Realname  string
	Server    string
	Port      int
	SSL       bool
	Channels  []string
	AdminChan string
	LogChan   string
	Modules   []string
	// Extra holds module-specific and other dynamic settings keyed by name.
	Extra map[string]string
}

// Load reads the bot configuration from the database for the given bot ID.
// The database schema is expected to have:
//   - bots(id, nick, username, realname, server, port, ssl, admin_chan, log_chan)
//   - bot_channels(bot_id, channel)
//   - bot_config(bot_id, key, value)
func Load(database *DB, botID int) (*Config, error) {
	pool := database.Pool()
	ctx := context.Background()

	cfg := &Config{
		BotID: botID,
		Extra: make(map[string]string),
	}

	row := pool.QueryRow(ctx,
		`SELECT nick, username, realname, server, port, ssl, admin_chan, log_chan, modules
		   FROM bots WHERE id = $1`, botID)

	if err := row.Scan(
		&cfg.Nick, &cfg.Username, &cfg.Realname,
		&cfg.Server, &cfg.Port, &cfg.SSL,
		&cfg.AdminChan, &cfg.LogChan, &cfg.Modules,
	); err != nil {
		return nil, fmt.Errorf("bot %d not found: %w", botID, err)
	}

	rows, err := pool.Query(ctx,
		`SELECT channel FROM bot_channels WHERE bot_id = $1`, botID)
	if err != nil {
		return nil, fmt.Errorf("loading channels: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ch string
		if err := rows.Scan(&ch); err != nil {
			return nil, err
		}
		cfg.Channels = append(cfg.Channels, ch)
	}

	rows2, err := pool.Query(ctx,
		`SELECT key, value FROM bot_config WHERE bot_id = $1`, botID)
	if err != nil {
		return nil, fmt.Errorf("loading extra config: %w", err)
	}
	defer rows2.Close()
	for rows2.Next() {
		var k, v string
		if err := rows2.Scan(&k, &v); err != nil {
			return nil, err
		}
		cfg.Extra[k] = v
	}

	return cfg, nil
}

// SetExtra persists a module-specific config key/value to the database.
func SetExtra(database *DB, botID int, key, value string) error {
	_, err := database.Pool().Exec(context.Background(),
		`INSERT INTO bot_config (bot_id, key, value)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (bot_id, key) DO UPDATE SET value = EXCLUDED.value`,
		botID, key, value)
	return err
}
