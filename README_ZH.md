# ccmux — Claude Code 多服务商复用器

[English](README.md)

管理多个 Claude Code API 服务商，按会话独立使用——秒切，零冲突。

> Claude Code 多服务商管理 | 切换 API 服务商 | 会话级隔离 | Profile 管理 | 支持 OpenRouter / AWS Bedrock / Anthropic API / 代理中转

```
$ ccc
  NAME               BASE_URL                    MODEL
→ my-proxy           proxy.example.com:8990      claude-opus-4-6-thinking
  openrouter         openrouter.ai/api           openrouter/pony-alpha
  bedrock            bedrock.us-east-1           claude-sonnet-4-6

$ ccc use openrouter
▸ launching claude with profile 'openrouter'
```

## 为什么需要

用 Claude Code 对接多个 API 服务商（Anthropic 直连、反向代理、OpenRouter、AWS Bedrock、Google Vertex、社区中转等）时，有两个痛点：每次切换都要手动改 `~/.claude/settings.json`，而且没法在不同终端同时使用不同的服务商。

ccmux（命令：`ccc`）同时解决这两个问题：

- `ccc use <name>` — 在**当前会话**使用指定服务商启动 Claude Code（不改全局配置）
- `ccc switch <name>` — 全局切换，所有新会话生效
- **不同会话同时使用不同服务商**——哪怕在同一个终端里，互不干扰

## 特性

- 会话级服务商隔离，基于 Claude Code 原生 `--settings` 覆盖机制
- 全局切换自动备份
- 模糊名称匹配（大小写不敏感、前缀匹配）
- 自动集成 [CC-Switch](https://github.com/farion1231/cc-switch) GUI，也可完全独立使用
- 零配置启动——`ccc import` 一键快照当前配置
- 单文件 shell 脚本，无需编译，唯一依赖 `jq`

## 安装

```bash
# 方式一：一键安装
curl -fsSL https://raw.githubusercontent.com/aiyi404/ccmux/main/install.sh | bash

# 方式二：克隆安装
git clone https://github.com/aiyi404/ccmux.git
cd ccmux && ./install.sh

# 方式三：直接软链接
ln -sf /path/to/ccmux/ccc ~/.local/bin/ccc
```

依赖 `jq`（`brew install jq` / `apt install jq`）。

## 快速上手

### 独立模式（无需其他依赖）

```bash
# 从当前配置导入为 profile
ccc import my-proxy

# 交互式创建新 profile
ccc add openrouter

# 列出所有 profile
ccc

# 在当前终端使用指定 profile（不改全局）
ccc use openrouter

# 全局切换（写入 settings.json）
ccc switch my-proxy
```

### 配合 CC-Switch 使用

如果你安装了 [CC-Switch](https://github.com/farion1231/cc-switch)，`ccc` 会自动检测其数据库并直接读取服务商配置——无需额外设置。你在 GUI 中配置的所有服务商，命令行里立即可用。

```bash
# 列出 CC-Switch 中的所有服务商
ccc

# 使用任意服务商
ccc use kirors

# 全局切换（同步回 CC-Switch）
ccc switch CPA
```

## 命令参考

| 命令 | 说明 |
|------|------|
| `ccc` / `ccc list` | 列出所有服务商，`→` 标记当前活跃的 |
| `ccc use <name> [-- args]` | 用指定服务商启动 Claude Code（仅当前终端，不改全局） |
| `ccc switch <name>` | 全局切换，写入 `settings.json` |
| `ccc current` | 显示当前活跃的服务商 |
| `ccc show <name>` | 查看服务商配置详情（token 自动脱敏） |

### 仅独立模式可用

| 命令 | 说明 |
|------|------|
| `ccc add <name>` | 交互式创建新 profile |
| `ccc edit <name>` | 用 `$EDITOR` 编辑 profile |
| `ccc rm <name>` | 删除 profile |
| `ccc import [name]` | 将当前 `settings.json` 导入为 profile |

### 选项

| 标志 | 说明 |
|------|------|
| `--standalone` | 强制使用独立模式（忽略 CC-Switch） |
| `--cc-switch` | 强制使用 CC-Switch 模式 |
| `-h, --help` | 显示帮助 |
| `-v, --version` | 显示版本 |

## 会话级服务商隔离

核心能力。每次 `ccc use` 启动一个独立的 Claude Code 会话，使用各自的服务商——哪怕在同一个终端里，前后两次会话也可以用不同的服务商：

```bash
# 会话 1 — 快速模型处理简单任务
ccc use sonnet-proxy

# （退出后，在同一个终端启动新会话）

# 会话 2 — opus 处理复杂架构工作
ccc use opus-proxy

# 也可以在多个终端并行运行不同会话：
# 终端 A
ccc use openrouter -- -c       # 继续上次对话

# 终端 B
ccc use bedrock -- -p "hello"  # print 模式
```

底层原理：`ccc use` 通过 Claude Code 原生的 `--settings` 选项注入配置。全局 `~/.claude/settings.json` 保持不变——hooks、权限、插件、MCP 服务器等配置不受影响。

## 模糊匹配

服务商名称支持大小写不敏感的前缀匹配：

```bash
ccc use op        # 匹配 "openrouter"
ccc switch cpa    # 匹配 "CPA"
ccc show ki       # 有歧义 → 列出候选：kirors, Kimi2.5, Kiro...
```

## Profile 格式（独立模式）

Profile 文件存放在 `~/.config/ccc/profiles/<name>.json`：

```json
{
  "name": "my-proxy",
  "description": "我的 API 代理",
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

只有 `env` 和 `model` 会传给 Claude Code。`name` 和 `description` 是 `ccc` 自用的元数据。

最简配置——只填你要覆盖的字段：

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

## 配置文件

可选配置文件 `~/.config/ccc/config.json`：

```json
{
  "mode": "auto",
  "default_profile": "my-proxy"
}
```

| 字段 | 可选值 | 说明 |
|------|--------|------|
| `mode` | `"auto"` / `"standalone"` / `"ccswitch"` | 覆盖模式检测 |
| `default_profile` | profile 名称 | `ccc use` 不带参数时的默认 profile |

环境变量：`CCC_MODE=standalone` 或 `CCC_MODE=ccswitch`。

## 工作原理

### `ccc use`（会话级）

利用 Claude Code 原生的 `--settings` 选项注入服务商配置作为覆盖层。全局 `~/.claude/settings.json` 保持不变——hooks、权限、插件等配置不受影响。

### `ccc switch`（全局）

将服务商配置写入 `~/.claude/settings.json`（自动备份到 `~/.claude/backups/`）。在 CC-Switch 模式下，还会同步数据库的 `is_current` 标记，使 GUI 界面反映变更。

## 双模式架构

| | CC-Switch 模式 | 独立模式 |
|---|---|---|
| 数据来源 | `~/.cc-switch/cc-switch.db` | `~/.config/ccc/profiles/*.json` |
| 自动启用条件 | CC-Switch 数据库存在 | 未检测到 CC-Switch |
| 服务商管理 | 在 CC-Switch GUI 中操作 | `ccc add/edit/rm/import` |
| 额外依赖 | `sqlite3`（macOS/Linux 自带） | 无 |

## 与 CC-Switch 生态的关系

- [CC-Switch](https://github.com/farion1231/cc-switch) — GUI 桌面应用，管理 AI 编程工具配置
- [CC-Switch CLI](https://github.com/SaladDay/cc-switch-cli) — 全功能 Rust CLI（服务商 + MCP + 代理 + 技能 + TUI）
- **ccmux** — 轻量 shell 脚本，专注快速切换服务商和会话级隔离

ccmux 是 CC-Switch 的命令行伴侣，不是替代品。需要 MCP 管理、代理路由或技能同步？用 CC-Switch CLI。只想快速切换服务商、不同会话用不同服务商？用 ccmux。

ccmux 也可以完全独立使用，无需安装 CC-Switch。

## 系统要求

- `bash` 3.2+（macOS 默认版本）或 `zsh`
- `jq` — JSON 处理工具
- `sqlite3` — 仅 CC-Switch 模式需要（macOS/Linux 自带）
- Claude Code CLI（`claude`）

## 许可证

MIT
