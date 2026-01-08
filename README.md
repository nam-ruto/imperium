# Imperium

<p align="center">
  <img src="assets/imp-trans-1.svg" alt="Imperium logo" width="580">
</p>

> *"Manage your commands."*

**Imperium** is a terminal-first command notebook for storing, discovering, and safely reusing shell commands.

In Latin, *imperium* means command, authority, and the power to direct action — fitting for a CLI tool that gives you mastery over your terminal knowledge.

## ✨ Features

- **Fuzzy Search** — Find commands instantly by typing any part of the command, description, or tags
- **Instant Startup** — Single Go binary, ~5ms launch time
- **Beautiful TUI** — Clean, fzf-inspired interface with transparent background
- **Copy to Clipboard** — Selected commands are automatically copied, ready to paste
- **Tag-Based Organization** — Categorize commands with tags for quick filtering
- **Local-First** — Your data stays on your machine in human-readable JSON

## 🚀 Installation

### From Source (Go)

```bash
# Clone and build
git clone https://github.com/youruser/imperium.git
cd imperium/imperium
go build -o bin/imp .

# Add to PATH
sudo mv bin/imp /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/youruser/imperium@latest
```

### Homebrew (coming soon)

```bash
brew tap youruser/tap
brew install imperium
```

## 🛠 Usage

### Search & Recall

Launch the interactive search. Type to filter, use **↑/↓** to navigate, and **Enter** to copy.

```bash
imp              # Launch search mode
imp search git   # Search with pre-filled query
```

### Add a Command

Add a new command through sequential prompts.

```bash
imp add
```

### Browse All Commands

```bash
imp browse
```

### Show Version

```bash
imp version
```

## ⌨️ Keybindings

| Key | Action |
|-----|--------|
| `↑` / `Ctrl+P` | Move selection up |
| `↓` / `Ctrl+N` | Move selection down |
| `Enter` | Copy selected command and exit |
| `Esc` | Clear query; if empty, exit |
| `Ctrl+C` / `q` | Exit immediately |
| `Ctrl+U` / `Ctrl+D` | Page up / Page down |

## 📂 Storage

Commands are stored in a human-readable JSON file:

```
~/.config/imperium/commands.json
```

### Data Format

```json
{
  "commands": [
    {
      "command": "git log --oneline -10",
      "usage": "Show last 10 commits",
      "tags": ["git", "log"],
      "platform": "macos",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 🏗 Project Structure

```
imperium/
├── main.go                 # Entry point
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── tui/               # Bubble Tea TUI
│   ├── storage/           # JSON storage
│   ├── search/            # Fuzzy search
│   └── models/            # Data models
└── Makefile               # Build commands
```

## 🔧 Development

```bash
cd imperium

# Build
make build

# Run directly
go run . search

# Build for all platforms
make build-all
```

## 📜 Philosophy

- **Fast Recall** — Commands are discoverable in seconds
- **Safe Reuse** — No automatic execution; commands are copied for you to inspect
- **Keyboard-First** — An fzf-inspired interface that stays out of your way
- **Zero Dependencies** — Single binary, no runtime required

## ⚖️ License

MIT
