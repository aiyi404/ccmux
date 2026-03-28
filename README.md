# ccmux — Claude Code Provider Multiplexer

[中文文档](README_ZH.md)

Manage multiple Claude Code API providers and use them per-session — switch instantly, no config conflicts.

> Claude Code multi-provider manager | Switch API provider | Per-session isolation | Profile management | OpenRouter / AWS Bedrock / Anthropic API / proxy support

```
$ ccc
  NAME               BASE_URL                    MODEL
→ my-proxy           proxy.example.com:8990      claude-opus-4-6-thinking
  openrouter         openrouter.ai/api           openrouter/pony-alpha
  bedrock            bedrock.us-east-1           claude-sonnet-4-6

$ ccc use openrouter
▸ launching claude with profile 'openrouter'
```

## Why

If you use multiple API providers with Claude Code — Anthropic direct API, reverse proxies, OpenRouter, AWS Bedrock, Google Vertex, or community relays — you know the pain: manually editing `~/.claude/settings.json` every time you want to switch, and no way to use different providers in different terminals at the same time.

ccmux (command: `ccc`) fixes both:

- `ccc use <name>` — launch Claude Code with a specific provider **for this session only** (global config untouched)
- `ccc switch <name>` — switch the global default for all new sessions
- Run **different providers in different sessions simultaneously** — no conflicts

## Features

- Per-session provider isolation via Claude Code's native `--settings` overlay
- Global provider switching with automatic backup
- Interactive TUI with visual menus, fuzzy search, and styled forms (powered by [gum](https://github.com/charmbracelet/gum))
- Fuzzy name matching (case-insensitive prefix)
- Auto-integration with [CC-Switch](https://github.com/farion1231/cc-switch) GUI, or fully standalone
- Zero config to start — `ccc import` snapshots your current setup
- Single shell script, no compilation, no runtime dependencies beyond `jq`

## Install

```bash
# Option 1: curl one-liner
curl -fsSL https://raw.githubusercontent.com/aiyi404/ccmux/main/install.sh | bash

# Option 2: clone and install
git clone https://github.com/aiyi404/ccmux.git
cd ccmux && ./install.sh

# Option 3: just symlink it
ln -sf /path/to/ccmux/ccc ~/.local/bin/ccc
```

Requires `jq` (`brew install jq` / `apt install jq`).

Optional: install [gum](https://github.com/charmbracelet/gum) for the interactive TUI experience:

```bash
# macOS
brew install gum

# Go
go install github.com/charmbracelet/gum@latest
```

When `gum` is installed, `ccc` (no args) launches an interactive TUI menu. Without `gum`, all commands fall back to plain text — no functionality is lost.

## Quick Start

### Standalone mode (no extra dependencies)

```bash
# Import your current settings as a profile
ccc import my-proxy

# Create a new profile interactively
ccc add openrouter

# List all profiles
ccc

# Use a profile in this terminal only (no global change)
ccc use openrouter

# Switch globally (writes to ~/.claude/settings.json)
ccc switch my-proxy
```

### With CC-Switch

If you have [CC-Switch](https://github.com/farion1231/cc-switch) installed, `ccc` auto-detects its database and reads providers directly — no setup needed. Everything you configured in the GUI is instantly available from the command line.

```bash
# List all providers from CC-Switch
ccc

# Use any provider
ccc use kirors

# Switch globally (syncs back to CC-Switch)
ccc switch CPA
```

## Commands

| Command | Description |
|---------|-------------|
| `ccc` | Interactive TUI menu (with `gum`) or list providers (without `gum`) |
| `ccc list` | List all providers, `→` marks the active one |
| `ccc use [name] [-- args]` | Launch Claude Code with a provider (session-scoped, no global change) |
| `ccc switch [name]` | Switch globally — writes to `settings.json` |
| `ccc current` | Show the active provider |
| `ccc show [name]` | Show provider config (tokens auto-masked) |
| `ccc tui` | Launch interactive TUI menu (requires `gum`) |

When `gum` is available, commands that take `[name]` will show a fuzzy search picker if the name is omitted.

### Standalone mode only

| Command | Description |
|---------|-------------|
| `ccc add [name]` | Create a new profile interactively |
| `ccc edit [name]` | Edit a profile with `$EDITOR` |
| `ccc rm [name]` | Remove a profile |
| `ccc import [name]` | Import current `settings.json` as a profile |

### Options

| Flag | Description |
|------|-------------|
| `--standalone` | Force standalone mode (ignore CC-Switch) |
| `--cc-switch` | Force CC-Switch mode |
| `-h, --help` | Show help |
| `-v, --version` | Show version |

## Interactive TUI

When [gum](https://github.com/charmbracelet/gum) is installed, `ccc` provides a visual interactive experience:

- `ccc` (no args) — launches a full TUI menu with styled header and action picker
- `ccc add` — styled form inputs with password masking, JSON preview, and confirmation
- `ccc use` / `ccc switch` / `ccc show` / `ccc edit` / `ccc rm` — fuzzy search provider picker when name is omitted
- `ccc list` — bordered table output
- `ccc rm` — styled confirmation dialog

All features gracefully fall back to plain text when `gum` is not available.

## Per-Session Provider Isolation

The killer feature. Each `ccc use` launches an independent Claude Code session with its own provider — even in the same terminal, consecutive sessions can use different providers:

```bash
# Session 1 — fast model for quick tasks
ccc use sonnet-proxy

# (exit, then start another session in the same terminal)

# Session 2 — opus for complex architecture work
ccc use opus-proxy

# Or run multiple sessions in parallel across terminals:
# Terminal A
ccc use openrouter -- -c    # continue last conversation

# Terminal B
ccc use bedrock -- -p "hello"  # print mode
```

Under the hood, `ccc use` injects config via Claude Code's native `--settings` flag. Your global `~/.claude/settings.json` stays untouched — hooks, permissions, plugins, MCP servers all remain as-is.

## Fuzzy Matching

Provider names support case-insensitive prefix matching:

```bash
ccc use op        # matches "openrouter"
ccc switch cpa    # matches "CPA"
ccc show ki       # ambiguous → shows candidates: kirors, Kimi2.5, Kiro...
```

## Profile Format (Standalone Mode)

Profiles are stored in `~/.config/ccc/profiles/<name>.json`:

```json
{
  "name": "my-proxy",
  "description": "My API proxy",
  "env": {
    "ANTHROPIC_BASE_URL": "http://proxy.example.com:8990",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "ANTHROPIC_MODEL": "claude-opus-4-6-thinking",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "claude-haiku-4-5-20251001",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "claude-opus-4-6-thinking",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "claude-sonnet-4-6-thinking",
    "ANTHROPIC_REASONING_MODEL": "claude-opus-4-6-thinking"
  },
  "model": "opus[1m]"
}
```

Only `env` and `model` are sent to Claude Code. `name` and `description` are metadata for `ccc` only.

Minimal profile — only include what you need to override:

```json
{
  "name": "minimal",
  "env": {
    "ANTHROPIC_BASE_URL": "https://api.example.com",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "ANTHROPIC_MODEL": "claude-sonnet-4-6"
  }
}
```

## Configuration

Optional config file at `~/.config/ccc/config.json`:

```json
{
  "mode": "auto",
  "default_profile": "my-proxy"
}
```

| Field | Values | Description |
|-------|--------|-------------|
| `mode` | `"auto"` / `"standalone"` / `"ccswitch"` | Override mode detection |
| `default_profile` | profile name | Default for `ccc use` without arguments |

Environment variable: `CCC_MODE=standalone` or `CCC_MODE=ccswitch`.

## How It Works

### `ccc use` (session-scoped)

Uses Claude Code's native `--settings` flag to inject provider config as an overlay. Your global `~/.claude/settings.json` stays untouched — hooks, permissions, plugins all remain as-is.

### `ccc switch` (global)

Writes the provider config into `~/.claude/settings.json` (with automatic backup to `~/.claude/backups/`). In CC-Switch mode, also syncs the database `is_current` flag so the GUI reflects the change.

## Dual Mode Architecture

| | CC-Switch Mode | Standalone Mode |
|---|---|---|
| Data source | `~/.cc-switch/cc-switch.db` | `~/.config/ccc/profiles/*.json` |
| Auto-detected when | CC-Switch database exists | No CC-Switch found |
| Provider management | Use CC-Switch GUI | `ccc add/edit/rm/import` |
| Extra dependency | `sqlite3` (pre-installed on macOS/Linux) | None |

## Relationship to CC-Switch Ecosystem

- [CC-Switch](https://github.com/farion1231/cc-switch) — GUI desktop app for managing AI coding tool configs
- [CC-Switch CLI](https://github.com/SaladDay/cc-switch-cli) — Full-featured Rust CLI (providers + MCP + proxy + skills + TUI)
- **ccmux** — Lightweight shell script focused on fast provider switching and per-session isolation

ccmux is a companion to CC-Switch, not a replacement. Need MCP management, proxy routing, or skill sync? Use CC-Switch CLI. Just want to switch providers fast and use different ones per session? Use ccmux.

ccmux also works completely standalone without CC-Switch installed.

## Requirements

- `bash` 3.2+ (macOS default) or `zsh`
- `jq` — JSON processing
- `gum` — optional, for interactive TUI ([install](https://github.com/charmbracelet/gum#installation))
- `sqlite3` — only in CC-Switch mode (pre-installed on macOS/Linux)
- Claude Code CLI (`claude`)

## License

MIT
