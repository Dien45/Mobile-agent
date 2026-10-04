# Arka Agent

Arka is a local-first personal AI agent written in Go. The browser is only the interface: provider calls, sessions, skills, files, commands, Git, and every production tool call run in the local daemon.

> The repository currently contains a runnable foundation, not every feature in the PRD. See [implementation status](docs/IMPLEMENTATION.md).

## Quick start

### macOS, Linux, Termux, or a Linux distro in Android PRoot

```bash
git clone https://github.com/Dien45/Mobile-agent.git
cd Mobile-agent
./install.sh
arka start
```

### Windows PowerShell

```powershell
git clone https://github.com/Dien45/Mobile-agent.git
cd Mobile-agent
.\install.ps1
arka start
```

The installer obtains a pinned Go toolchain when needed, builds one local binary, and does not require Node.js. Arka listens on `127.0.0.1:7331` and opens an authenticated local URL.

## Current capabilities

- Embedded responsive claymorphism web UI—no CDN or separate frontend server.
- Local sessions with rename, trash, provider/model selection, and background execution.
- Custom OpenAI-compatible provider profiles for OmniRoute and similar gateways.
- Model discovery through `/models` and manual per-session model selection.
- Real local agent tool loop with up to 1,000 tool calls per task.
- Workspace-confined file tools, local commands, and Git status.
- Local terminal command panel.
- Local skill install from Git, directories, or ZIP archives.
- Loopback-only server with authenticated bootstrap cookie.
- Windows, macOS, Linux x64/arm64, Termux, and PRoot-aware installers.

## Product direction

- **Web-first and local:** storage, credentials, memory, skills, artifacts, terminal, browser automation, and tools belong in the daemon—not JavaScript mocks.
- **Provider-flexible:** common provider presets plus Custom OpenAI-compatible for OmniRoute.
- **Long-running:** more than 500 tool calls, checkpoints, and sessions that continue when the user navigates away.
- **Extensible:** installable Agent Skills-style packages and local MCP tools.
- **Safe autonomy:** persistent workspace grants, risk-based approval, audit records, and explicit remote Git permissions.
- **No messaging gateway:** WhatsApp, Telegram, Discord, and Slack are intentionally out of scope.

## Development

Requires Go 1.22 or newer (the installer pins Go 1.24.2):

```bash
make test
make build
./bin/arka start
```

- [Product requirements](docs/PRD.md)
- [Implementation status](docs/IMPLEMENTATION.md)

The product name is **Arka**.
