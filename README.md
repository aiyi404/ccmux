# ccmux — Claude Code Provider Multiplexer

[中文文档](README_ZH.md)

Manage multiple Claude Code API providers and switch between them instantly — per-session isolation, no config conflicts.

```
$ ccc use openrouter
▸ launching claude with profile 'openrouter'
```

## Why

If you use multiple API providers with Claude Code — Anthropic direct API, reverse proxies, OpenRouter, AWS Bedrock, Google Vertex — switching means manually editing `~/.claude/settings.json`, and there's no way to use different providers in different terminals simultaneously.

ccmux (`ccc`) fixes both:

- `ccc use <name>` — launch Claude Code with a specific provider **for this session only** (global config untouched)
- `ccc switch <name>` — switch the global default for all new sessions
- Run **different providers in different terminals simultaneously** — no conflicts

## Features

- Per-session provider isolation via Claude Code's native `--settings` overlay
- Global provider switching with automatic backup
- Interactive TUI with sidebar navigation, Dracula theme (powered by [bubbletea](https://github.com/charmbracelet/bubbletea))
- Fuzzy name matching (case-insensitive prefix)
- Import from [cc-switch](https://github.com/farion1231/cc-switch) database via `ccc import-all`
- Single binary, no runtime dependencies
- EN/ZH bilingual support

## Install

```bash
# Clone and install (requires Go 1.22+)
git clone https://github.com/aiyi404/ccmux.git
cd ccmux && ./install.sh
```

Or build manually:

```bash
go build -ldflags="-s -w" -o ccc .
mv ccc ~/.local/bin/
```

## Quick Start

```bash
# Import your current settings as a profile
ccc import my-proxy

# Create a new profile interactively
ccc add openrouter

# List all profiles
ccc list

# Use a profile in this terminal only (no global change)
ccc use openrouter

# Switch globally (writes to ~/.claude/settings.json)
ccc switch my-proxy

# Launch interactive TUI
ccc
```

### Migrating from cc-switch-cli

If you have [cc-switch-cli](https://github.com/SaladDay/cc-switch-cli) installed, import all providers in one command:

```bash
ccc import-all
```

Deduplication is by `ANTHROPIC_BASE_URL` + `ANTHROPIC_MODEL` — existing profiles won't be overwritten. On first TUI launch, ccmux will detect cc-switch and offer to import automatically.

## Commands

| Command | Description |
|---------|-------------|
| `ccc` | Launch interactive TUI |
| `ccc list` | List all providers, `→` marks the active one |
| `ccc use <name> [-- args]` | Launch Claude Code with a provider (session-scoped) |
| `ccc switch <name>` | Switch globally — writes to `settings.json` |
| `ccc current` | Show the active provider |
| `ccc show <name>` | Show provider config (tokens auto-masked) |
| `ccc add <name>` | Create a new profile interactively |
| `ccc edit <name>` | Edit a profile with `$EDITOR` |
| `ccc rm <name>` | Remove a profile |
| `ccc import [name]` | Import current `settings.json` as a profile |
| `ccc import-all` | Batch import all providers from cc-switch database |

## Per-Session Provider Isolation

The core feature. `ccc use` launches Claude Code with a specific provider for the current session only:

```bash
# Terminal A — fast model
ccc use sonnet-proxy

# Terminal B — opus for complex work
ccc use opus-proxy -- -c    # continue last conversation
```

Under the hood, `ccc use` creates a temp settings file and passes it via `--settings`. Your global `~/.claude/settings.json` stays untouched — hooks, permissions, MCP servers all remain as-is.

## Fuzzy Matching

Provider names support case-insensitive prefix matching:

```bash
ccc use op        # matches "openrouter"
ccc switch ki     # matches "kirors"
```

## Profile Format

Profiles are stored in `~/.config/ccc/profiles/<name>.json`:

```json
{
  "name": "my-proxy",
  "env": {
    "ANTHROPIC_BASE_URL": "http://proxy.example.com:8990",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "ANTHROPIC_MODEL": "claude-opus-4-6-thinking"
  },
  "model": "opus[1m]"
}
```

## Configuration

Config file at `~/.config/ccc/config.json`:

```json
{
  "lang": "zh",
  "current": "my-proxy"
}
```

| Field | Description |
|-------|-------------|
| `lang` | UI language: `"en"` or `"zh"` |
| `current` | Currently active profile name |

## Requirements

- Go 1.22+ (build time only)
- Claude Code CLI (`claude`)

## License

MIT
