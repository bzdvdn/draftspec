# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-03-31

### Added

- Initial release of the Draftspec CLI
- File-based project context system for development agents and humans
- Eight-phase workflow: `constitution → spec → inspect → plan → tasks → implement → verify → archive`
- Bilingual support for English (`en`) and Russian (`ru`) templates and prompts
- Agent integration for Claude, Codex, Copilot, Cursor, Kilocode, and Trae
- Shell support for `sh` and `powershell`
- `init` and `refresh` commands to manage the `.draftspec/` workspace
- `doctor` command with comprehensive workspace health checks
- `status` and `features` commands for feature lifecycle visibility
- `feature`, `feature repair`, `list-specs`, `show-spec`, `migrate` commands
- `add-agent`, `remove-agent`, `list-agents`, `cleanup-agents` commands for agent artifact management
- Internal CLI (`__internal`) for script delegation without exposing plumbing commands
- Phase readiness scripts for each workflow phase
- Stable IDs for traceability: `AC-*`, `RQ-*`, `DEC-*`, `T*`
- Phase Contract headers in all agent prompts with concrete inputs, outputs, and stop conditions
- Self-Check sections in inspect and verify prompts
- `--version` flag reporting the build version
- Multi-platform CI builds for linux/amd64, linux/arm64, windows/amd64, windows/arm64

[Unreleased]: https://github.com/bzdv/draftspec/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/bzdv/draftspec/releases/tag/v0.1.0
