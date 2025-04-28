# gid

> **gid** — Tiny Git Identity toolkit for Ed25519 keys ⚡

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![CI](https://github.com/yourusername/gid/actions/workflows/ci.yml/badge.svg)

## Features

- 🔐 Generate new Ed25519 key pairs (private + public keys)
- 📂 Load existing identities from disk
- 🌎 Fetch public keys remotely from GitHub
- ⚡ Lightweight, no external crypto dependencies

## Installation

```bash
go install github.com/mynextid/gid/cmd/gid@latest
```

Or clone manually:

```bash
git clone https://github.com/mynextid/gid.git
cd cli
go build
```

## Usage

```bash
gid [command] [flags]
```

Available commands:

| Command                  | Description                                          |
| :----------------------- | :--------------------------------------------------- |
| `generate [path]`        | Generate a new Ed25519 identity and save to `[path]` |
| `load [path]`            | Load an existing identity and display the public key |
| `fetch [GitHub handler]` | Fetch a public key from a GitHub repository          |

### Generate an Identity

```bash
gid generate ./keys/identity.pem
```

This will:

- Create a new Ed25519 private key at `./keys/identity.pem`
- Save the corresponding public key at `./keys/gid.pem`

### Load an Identity

```bash
gid load ./keys/identity.pem
```

Outputs the **public key** from your saved identity.

### Fetch a Public Key from GitHub

```bash
gid fetch user
```

Fetches and displays the public key from a GitHub repo, branch, and file path.

Example:

```bash
gid fetch mynextid
```

## Security Notes

- Private keys are saved with **0600 permissions** (`rw-------`) by default.
- Only Ed25519 keys are supported (simple, modern, safe).

## Project Structure

```plaintext
/cmd/         CLI commands (Cobra-based)
/identity/    Core identity logic (pure Go library)
/main.go      CLI entry point
```

## 🛠️ Development

```bash
# Install dependencies
go mod tidy

# Build CLI
go build -o gid

# Run CLI directly
go run ./gid generate ./mykeys/id.pem
```

---

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome!  
Feel free to open issues, pull requests, or suggest new features.
