# pmproxy

Forwards private messages (PMs) sent directly to the bot to the configured admin channel, so operators can monitor and respond to them.

## Module name

`pmproxy`

## Behaviour

On every `PRIVMSG` event:

1. Checks whether the message target is the bot's own nick (i.e. it is a PM, not a channel message).
2. If `admin_chan` is set, relays the message to that channel in the format:

```
[PM from <nick>] <message text>
```

Messages addressed to channels are silently ignored.

## Configuration

This module has no `bot_config` keys. It reads `admin_chan` directly from the `bots` table row.

| Setting      | Source      | Description                              |
|--------------|-------------|------------------------------------------|
| `admin_chan` | `bots` table | Channel where proxied PMs are posted    |

## Commands

This module sends no commands and provides no chat interface. It is purely passive.
