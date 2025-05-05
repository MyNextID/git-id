# gid

**gid** — Tiny Git Identity toolkit for Ed25519 keys ⚡

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 🔐 Generate new Ed25519 key pairs (private + public keys)
- 📂 Load existing identities from disk
- 🌎 Fetch public keys remotely from GitHub
- ⚡ Lightweight, no external crypto dependencies

## Installation

```bash
go install github.com/mynextid/git-id/cli@latest
```

Or clone manually:

```bash
git clone https://github.com/mynextid/git-id.git
cd cli
go build
```

## Usage

```bash
git-id [command] [flags]
```

Available commands:

| Command                  | Description                                          |
| :----------------------- | :--------------------------------------------------- |
| `generate [path]`        | Generate a new Ed25519 identity and save to `[path]` |
| `load [path]`            | Load an existing identity and display the public key |
| `fetch [GitHub handler]` | Fetch a public key from a GitHub repository          |

### Generate an Identity

```bash
git-id generate ./keys/identity.pem
```

This will:

- Create a new Ed25519 private key at `./keys/identity.pem`
- Save the corresponding public key at `./keys/git-id.pem`

### Load an Identity

```bash
git-id load ./keys/identity.pem
```

Outputs the **public key** from your saved identity.

### Fetch a Public Key from GitHub

```bash
git-id fetch user
```

Fetches and displays the public key from a GitHub repo, branch, and file path.

Example:

```bash
git-id fetch mynextid
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

## Development

```bash
# Install dependencies
go mod tidy

# Build CLI
go build

# Run CLI directly
go run ./git-id generate ./mykeys/id.pem
```

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome!  
Feel free to open issues, pull requests, or suggest new features.
