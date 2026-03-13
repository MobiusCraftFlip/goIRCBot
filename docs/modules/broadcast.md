# broadcast

Sends a message to every channel the bot is currently in, announced with the sender's nick.

## Module name

`broadcast`

## Behaviour

On receiving the `broadcast` command:

1. Checks that the command originated from `admin_chan`. Commands from any other source are silently ignored.
2. Sends the following message to every channel returned by the IRC client's channel list:

```
!!! Broadcast from <nick> <message text>
```

If no message text is supplied, replies to the sender with usage instructions.

## Configuration

This module has no `bot_config` keys. It reads `admin_chan` directly from the `bots` table row.

| Setting      | Source       | Description                                        |
|--------------|--------------|----------------------------------------------------|
| `admin_chan` | `bots` table | The only channel from which the command is accepted |

## Commands

| Command     | Arguments   | Description                              |
|-------------|-------------|------------------------------------------|
| `broadcast` | `<message>` | Sends `<message>` to every joined channel |

**Syntax:** `.<botnick> broadcast <message>`

**Example:**

```
.<botnick> broadcast Server restarting in 5 minutes
```

Sends `!!! Broadcast from <nick> Server restarting in 5 minutes` to all channels the bot is in.
