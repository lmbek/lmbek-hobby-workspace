# Install Docker

One-shot installer scripts for Docker across different platforms.

## Variants

| Directory | Platform | What it installs |
|---|---|---|
| `windows/` | Windows (native) | Docker Desktop via PowerShell — no WSL required |
| `windows-wsl/` | Windows + WSL | Docker Engine inside a WSL Ubuntu distribution |
| `linux/` | Linux (Ubuntu/Debian) | Docker Engine natively via `apt` |

## Which One Should I Use?

- **Windows user who wants a GUI and the simplest setup** → `windows/`
- **Windows user who wants a lightweight, license-free Docker** → `windows-wsl/`
- **Linux user (bare metal or VM, no WSL)** → `linux/`

See each subdirectory's README for detailed usage and post-install instructions.
