# core

Handles essential IRC housekeeping on every connection: NickServ authentication, IRC operator login, and joining the admin and log channels.

## Module name

`core`

## Behaviour

On `CONNECTED`:

1. If `ns_user` and `ns_pass` are both set, sends `PRIVMSG NickServ LOGIN <user> <pass>` to identify the bot's nick.
2. If `oper_user` and `oper_pass` are both set, sends the `OPER` command to obtain IRC operator privileges.
3. Joins `admin_chan` and `log_chan` (as configured in the `bots` table).

## Configuration

All keys are stored in the `bot_config` table (`key` / `value` pairs).

| Key         | Required | Description                                      |
|-------------|----------|--------------------------------------------------|
| `ns_user`   | No       | NickServ account name used for `LOGIN`           |
| `ns_pass`   | No       | NickServ password used for `LOGIN`               |
| `oper_user` | No       | IRC operator username sent with the `OPER` command |
| `oper_pass` | No       | IRC operator password sent with the `OPER` command |

`admin_chan` and `log_chan` are columns on the `bots` table, not `bot_config` keys.

## Commands

This module sends no user-facing commands and exposes no chat interface.
