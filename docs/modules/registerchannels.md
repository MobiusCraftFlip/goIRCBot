# registerchannels

Registers and deregisters channels with ChanServ, automatically setting access flags for configured user groups. Supports single-channel and ranged multi-channel operations.

## Module name

`registerchannels`

## Behaviour

### Registration (`registerchan`, `registermultichan`)

For each channel being registered:

1. If the bot is not already in the channel, sends `JOIN`.
2. If the bot does not hold `+o` in the channel, sends `MODE <channel> +o <botnick>`.
3. Sends `PRIVMSG ChanServ :REGISTER <channel>`.
4. For each configured user group (coordinators, operators, supporters), sends `PRIVMSG ChanServ :FLAGS <channel> <nick> <flags>` for every nick in that group's list.

### Deregistration (`deregisterchan`, `deregistermultichan`)

Sends `PRIVMSG ChanServ :FDROP <channel>` for each target channel. No join or mode checks are performed.

### Multi-channel ranging

Multi-channel commands accept a base channel name and an inclusive numeric range. Channel numbers are zero-padded to at least two digits (e.g. `1`→`01`, `10`→`10`, `100`→`100`).

## Configuration

| Key | Description |
|-----|-------------|
| `registerchannels.coordinators` | Comma-separated list of nicks/accounts to receive coordinator flags |
| `registerchannels.operators` | Comma-separated list of nicks/accounts to receive operator flags |
| `registerchannels.supporters` | Comma-separated list of nicks/accounts to receive supporter flags |
| `registerchannels.coordinators_flags` | ChanServ FLAGS string applied to coordinators (e.g. `+AFRefiorstv`) |
| `registerchannels.operators_flags` | ChanServ FLAGS string applied to operators (e.g. `+AORefiorstv`) |
| `registerchannels.supporters_flags` | ChanServ FLAGS string applied to supporters (e.g. `+Vv`) |

Any group whose nick list or flags string is empty is silently skipped.

## Commands

| Command | Arguments | Description |
|---------|-----------|-------------|
| `registerchan` | `<#channel>` | Registers a single channel with ChanServ and sets group flags |
| `registermultichan` | `<#base> <start> <end>` | Registers a numbered range of channels |
| `deregisterchan` | `<#channel>` | Sends ChanServ FDROP for a single channel |
| `deregistermultichan` | `<#base> <start> <end>` | Sends ChanServ FDROP for a numbered range of channels |

**Syntax:**

```
.<botnick> registerchan <#channel>
.<botnick> registermultichan <#base> <start> <end>
.<botnick> deregisterchan <#channel>
.<botnick> deregistermultichan <#base> <start> <end>
```

**Examples:**

```
.<botnick> registerchan #general
```
Joins `#general` if needed, ensures op, registers with ChanServ, and sets flags for all configured groups.

```
.<botnick> registermultichan #room 1 5
```
Registers `#room01` through `#room05`.

```
.<botnick> deregistermultichan #room 1 5
```
Sends `FDROP` to ChanServ for `#room01` through `#room05`.
