# Gitver - Git Semantic Versioning Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/arcalyx/gitver)](https://goreportcard.com/report/github.com/arcalyx/gitver)
[![License: BSL 1.1](https://img.shields.io/badge/License-BSL%201.1-blue.svg)](https://github.com/arcalyx/gitver/blob/main/LICENSE)

Gitver is a command-line tool that helps you manage semantic versioning in your Git repositories. It provides commands for bumping versions, managing Git hooks, and generating shell completions.

## Features

- **Semantic Versioning**: Follows [SemVer 2.0.0](https://semver.org/) specification
- **Version Bumping**: Easily bump major, minor, or patch versions
- **Auto Bump**: Automatically determine version bump based on commit messages
- **Pre-release Support**: Manage alpha, beta, and release candidate versions
- **Build Metadata Support**: Include build metadata in version strings (e.g., 1.2.3+build.123)
- **Git Integration**: Create tags, commits, and push changes
- **Package Manager Integration**: Automatically update versions in package manager files (Go, Python)
- **Git Hooks**: Install, list, and uninstall Git hooks for version management
- **Shell Completion**: Generate shell completion scripts for bash, zsh, and fish

## Installation

### Using Go

```bash
go install github.com/arcalyx/gitver@latest
```

### From Source

```bash
git clone https://github.com/arcalyx/gitver.git
cd gitver
./build.sh --install
```

## Quick Start

1. Initialize a new project:

```bash
gitver init
```

2. Bump the version:

```bash
gitver bump minor
```

3. View the current version:

```bash
gitver version
```

## Commands

### Basic Commands

- `gitver init`: Initialize gitver in your project
- `gitver version`: Display version information

### Bump Command

- `gitver bump`: Bump the version of the project
  - `major`: Bump the major version
  - `minor`: Bump the minor version
  - `patch`: Bump the patch version
  - `--auto`: Automatically determine version bump based on commit messages
  - `--prerelease <type>`: Set pre-release version (alpha, beta, rc)
  - `--build-meta <metadata>`: Set build metadata (e.g., build.123, sha.abc123)

### Config Command

- `gitver config`: Manage configuration
  - `init`: Initialize configuration
  - `list`: List configuration values
  - `set <key> <value>`: Set configuration value
  - `get <key>`: Get configuration value

### Hooks Command

- `gitver hooks`: Manage Git hooks
  - `install`: Install Git hooks
  - `uninstall`: Uninstall Git hooks
  - `list`: List installed Git hooks

### Completion Command

- `gitver completion`: Generate shell completion scripts
  - `bash`: Generate bash completion script
  - `zsh`: Generate zsh completion script
  - `fish`: Generate fish completion script

## Examples

### Init

Initialize a new project with version 0.1.0:

```bash
gitver init
```

Initialize with a specific version:

```bash
gitver init --version 1.0.0
```

### Bump

Bump the version number:

```bash
gitver bump major     # Bump major version (X.0.0)
gitver bump minor     # Bump minor version (0.X.0)
gitver bump patch     # Bump patch version (0.0.X)
```

Bump with pre-release or build metadata:

```bash
gitver bump minor --prerelease beta   # Bump to 0.X.0-beta.1
gitver bump patch --build-meta build.123   # Bump to 0.0.X+build.123
```

Automatically determine bump type based on commit messages:

```bash
gitver bump --auto
```

### Config

Set configuration values:

```bash
gitver config set hooks.enabled true
```

List all configuration values:

```bash
gitver config list
```

### Hooks

Install Git hooks:

```bash
gitver hooks install
```

List installed hooks:

```bash
gitver hooks list
```

## License

This project is licensed under the Business Source License 1.1 - see the [LICENSE](LICENSE) file for details.
