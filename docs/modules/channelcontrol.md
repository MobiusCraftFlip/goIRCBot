# channelcontrol

Opens and closes channels by applying a configured topic and mode string. Supports single-channel and ranged multi-channel operations.

## Module name

`channelcontrol`

## Behaviour

Each command applies a pre-configured **topic** and **mode string** to the target channel(s). The topic and modes are global — the same values are applied regardless of which channel is targeted.

- If the `open` state is requested, the open topic and open modes are applied.
- If the `close` state is requested, the close topic and close modes are applied.
- If neither a topic nor a mode string is found in config for the requested state, the command replies with an error and takes no action.

Multi-channel commands (`openmultichans`, `closemultichans`) accept a base channel name and an inclusive numeric range. Channel numbers are zero-padded to at least two digits (e.g. `1`→`01`, `10`→`10`, `100`→`100`).

## Configuration

| Key | Description |
|-----|-------------|
| `channelcontrol.open_topic` | Topic to set when opening a channel |
| `channelcontrol.open_modes` | Mode string to apply when opening (e.g. `-im`) |
| `channelcontrol.close_topic` | Topic to set when closing a channel |
| `channelcontrol.close_modes` | Mode string to apply when closing (e.g. `+im`) |

Both `open_topic` and `open_modes` (and their `close_` equivalents) are optional individually, but at least one of the pair must be present for the command to act.

## Commands

| Command | Arguments | Description |
|---------|-----------|-------------|
| `openchan` | `<#channel>` | Opens a single channel |
| `closechan` | `<#channel>` | Closes a single channel |
| `openmultichans` | `<#base> <start> <end>` | Opens a numbered range of channels |
| `closemultichans` | `<#base> <start> <end>` | Closes a numbered range of channels |

**Syntax:**

```
.<botnick> openchan <#channel>
.<botnick> closechan <#channel>
.<botnick> openmultichans <#base> <start> <end>
.<botnick> closemultichans <#base> <start> <end>
```

**Examples:**

```
.<botnick> openchan #general
```
Sets the open topic and applies the open mode string to `#general`.

```
.<botnick> closemultichans #room 1 5
```
Closes `#room01`, `#room02`, `#room03`, `#room04`, and `#room05`.
