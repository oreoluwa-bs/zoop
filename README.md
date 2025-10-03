# Zoop

A fast, minimal CLI for storing and retrieving anything—API keys, passwords, tokens, notes, secrets.

## Features

- **Simple**: Store and retrieve data with key-value pairs
- **Secure**: Optional encryption using age (X25519)
- **Fast**: Minimal dependencies, quick operations
- **Cross-platform**: Works on Linux, macOS, Windows

## Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/oreoluwa-bs/zoop@latest
```

### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/oreoluwa-bs/zoop/releases).

Make the binary executable and move to your PATH:

```bash
chmod +x zoop
sudo mv zoop /usr/local/bin/
```

### Option 3: Build from Source

```bash
git clone https://github.com/oreoluwa-bs/zoop.git
cd zoop
go build -o zoop main.go
```

## Quick Start

1. Initialize Zoop:

```bash
zoop init
```

This creates `~/.zoop/` directory with config and keys.

2. Store a secret:

```bash
zoop set my-api-key sk-1234567890abcdef
```

3. Retrieve it:

```bash
zoop get my-api-key
```

## Usage

### Commands

- `zoop init [--encrypt] [--key-file PATH] [--force]`: Initialize Zoop. Generates encryption keys if encryption enabled.
- `zoop set KEY VALUE`: Store a value with the given key.
- `zoop get KEY`: Retrieve the value for the given key.
- `zoop delete KEY`: Delete the value for the given key.
- `zoop config set KEY=VALUE`: Set a configuration option.
- `zoop config get KEY`: Get a configuration value.
- `zoop config list`: List all configuration values.
- `zoop migrate [decrypt|encrypt]`: Migrate data between encrypted and unencrypted storage.
- `zoop version`: Show version information.

### Configuration

Zoop uses a config file at `~/.zoop/config.yaml`. Default settings:

- `data_file`: `~/.zoop/store.json`
- `key_file`: `~/.zoop/key.txt`
- `encryption`: `false`

You can override with environment variables prefixed with `ZOOP_`, e.g., `ZOOP_ENCRYPTION=true`.

### Encryption

By default, encryption is disabled. To enable:

```bash
zoop init --encrypt
```

Or after init:

```bash
zoop config set encryption=true
```

Then migrate existing data:

```bash
zoop migrate encrypt
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Cross-platform Builds

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o zoop-linux-amd64 main.go

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o zoop-darwin-amd64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o zoop-windows-amd64.exe main.go
```

## License

MIT