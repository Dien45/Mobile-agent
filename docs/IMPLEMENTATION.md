# Implementation status

This repository now contains the first runnable Arka foundation. The PRD remains the product contract; this file records what the current code actually implements.

## Implemented

- Pure-Go local daemon and embedded claymorphism web UI.
- Loopback-only Host validation and authenticated local cookie bootstrap.
- Persistent local sessions, rename, soft delete, restore API, and purge.
- Background agent turns that continue when the browser changes sessions or disconnects.
- Custom OpenAI-compatible profiles suitable for OmniRoute.
- `/models` discovery and per-session provider/model selection.
- Real local tool loop with a 1,000-call ceiling.
- Workspace-confined file list/read/atomic-write tools.
- Local terminal command plus Git status, diff, explicit-path staging, and local commit tools (never implicit push).
- Local web terminal command panel (command mode).
- Skill listing/install from Git, local directory, or ZIP; uninstall and archive safety limits.
- Cross-platform bootstrap scripts for Unix/Termux and Windows.
- Standard-library-only runtime: no Node.js frontend or CDN.

## Next milestones before claiming PRD completion

- Replace the atomic JSON state backend with pure-Go SQLite + migrations/FTS while preserving the repository interface.
- Add durable event cursors, fair multi-session queue, checkpoints, and crash resume.
- Add SSE streaming; the current UI polls durable local state.
- Add native Anthropic/Gemini adapters and provider presets beyond the OpenAI-compatible transport.
- Implement risk-profile approvals and workspace grants in the UI/runtime.
- Replace terminal command mode with PTY/ConPTY WebSocket transport.
- Add local Chromium/CDP browser automation.
- Complete skill update/diff/rollback/catalog and dependency inspection.
- Move provider secrets from the mode-0600 state file to OS keyring/encrypted fallback.
- Add GitHub Device Flow, remote Git capabilities, and durable commit audit records.
- Add Android Termux/PRoot device smoke tests and capability probes.

No browser-side mock tool executor exists. Tool and provider actions in production go through the local Go daemon.
