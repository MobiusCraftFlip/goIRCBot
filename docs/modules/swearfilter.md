# swearfilter

Monitors channel messages and kicks users whose messages match any configured regex pattern.

## Module name

`swearfilter`

## Behaviour

On every channel `PRIVMSG`:

1. Ignores private messages (targets not beginning with `#`).
2. Ignores messages sent by the bot itself.
3. Tests the message text against each compiled pattern in order.
4. On the first match, kicks the sender from the channel with the configured kick message and stops processing further patterns.

Patterns are compiled once at startup as case-insensitive regular expressions. Invalid patterns are logged and skipped; they do not prevent the module from loading.

## Configuration

| Key | Description |
|-----|-------------|
| `swearfilter.patterns` | Newline-separated list of Go regex patterns to match against message text |
| `swearfilter.kick_message` | Reason string sent with the kick |

**Pattern format:** one regex per line. Each pattern is automatically wrapped with `(?i)` for case-insensitive matching. Standard Go [`regexp`](https://pkg.go.dev/regexp/syntax) syntax applies.

**Example `swearfilter.patterns` value:**

```
\bbadword\b
offensive+phrase
another.*pattern
```

## Commands

This module has no commands. It operates passively on all incoming channel messages.
