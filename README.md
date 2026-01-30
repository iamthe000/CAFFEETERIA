# CAFFEETERIA

**CAFFEE_Editorの軽量版です**
<a href="https://github.com/iamthe000/CAFFEE_Editor">CAFFEE_Editor</a>

[![License](https://img.shields.io/badge/license-GNU-blue.svg)](LICENSE)

**CAFFEETERIA** is a ultra-lightweight, TUI-based text editor written in Go. Designed for developers who need a simple, fast, and no-frills editing environment within their terminal.

## 🚀 Features

- **Minimalist TUI**: A distraction-free interface that stays out of your way.
- **Efficient Command System**: Access file operations and utilities via a familiar colon-prefixed command line.
- **Project Exploration**: Built-in directory tree generator to visualize your project structure.
- **Standard Keybindings**: Supports common terminal editor shortcuts for a smooth transition.
- **Cross-Platform**: Leverage Go's portability to run CAFFEETERIA on Linux, macOS, and Windows.

## 🛠 Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.24.4 or later)

### Build from Source

Clone the repository and build the executable:

```bash
git clone https://github.com/yourusername/caffeeteria.git
cd caffeeteria
go build -o caffeeteria main.go
```

Move the binary to your path (optional):

```bash
mv caffeeteria /usr/local/bin/
```

## ⌨️ Usage

Launch the editor by running:

```bash
./caffeeteria
```

### Keybindings

| Key | Action |
|-----|--------|
| `Ctrl + S` | Save the current file |
| `Ctrl + O` | Prompt to open a file |
| `Ctrl + X` | Exit CAFFEETERIA |
| `:` | Enter Command Mode |
| `Arrows` | Navigate the cursor |
| `Backspace` | Delete character / Join lines |
| `Enter` | New line |

### Command Reference

Press `:` to enter command mode, then type one of the following:

- `open <filename>`: Load a file into the buffer.
- `new <filename>`: Create a new file buffer.
- `save [filename]`: Save the current buffer. If a filename is provided, it saves as that file.
- `file_txt`: Generate a visual tree of the current directory.
- `Esc`: Exit command mode.

## 📁 Project Structure

CAFFEETERIA is designed with simplicity in mind:

- `main.go`: The core logic of the editor, handling TUI rendering and event loops.
- `go.mod` / `go.sum`: Dependency management.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Made with ☕ and Go.*
