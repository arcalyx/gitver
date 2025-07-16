# Gitver - Git Semantic Versioning Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/arcalyx/gitver)](https://goreportcard.com/report/github.com/arcalyx/gitver)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Gitver is a command-line tool that helps you manage semantic versioning in your Git repositories. It provides commands for bumping versions, creating releases, managing pre-release versions, and generating changelogs.

## Features

- **Semantic Versioning**: Follows [SemVer 2.0.0](https://semver.org/) specification
- **Version Bumping**: Easily bump major, minor, or patch versions
- **Auto Bump**: Automatically determine version bump based on commit messages
- **Pre-release Support**: Manage alpha, beta, and release candidate versions
- **Git Integration**: Create tags, commits, and push changes
- **Changelog Generation**: Generate changelogs based on commit messages
- **Version History**: View version history and details
- **Package Manager Integration**: Automatically update versions in package manager files (Maven, NPM, Gradle, Kotlin Gradle, Go, Python)

## Installation

### Using Go

```bash
go install github.com/arcalyx/gitver@latest
```

### From Source

```bash
git clone https://github.com/arcalyx/gitver.git
cd gitver
./install.sh
```

## Quick Start

1. Initialize a new project:

```bash
gitver init
```

2. Bump the version:

```bash
gitver bump --minor
```

3. Create a release:

```bash
gitver release
```

## Commands

### Init

Initialize a new project with version 0.1.0:

```bash
gitver init
```

### Bump

Bump the version number:

```bash
gitver bump --major     # Bump major version (X.0.0)
gitver bump --minor     # Bump minor version (0.X.0)
gitver bump --patch     # Bump patch version (0.0.X)
gitver bump --auto      # Automatically determine version bump based on commit messages
```

Options:
- `--commit`: Create a commit with the version change
- `--tag`: Create a tag for the new version
- `--push`: Push changes to remote repository
- `--update-packages`: Update version in package manager files (default: true)

### Prerelease

Manage pre-release versions:

```bash
gitver prerelease --alpha     # Set/bump alpha version (1.0.0-alpha.1)
gitver prerelease --beta      # Set/bump beta version (1.0.0-beta.1)
gitver prerelease --rc        # Set/bump release candidate version (1.0.0-rc.1)
gitver prerelease --bump      # Bump the current pre-release version
gitver prerelease --promote   # Promote to the next pre-release level (alpha → beta → rc)
gitver prerelease --clear     # Clear pre-release version (1.0.0-beta.1 → 1.0.0)
```

Options:
- `--commit`: Create a commit with the version change
- `--tag`: Create a tag for the new version
- `--push`: Push changes to remote repository

### Release

Create a release tag:

```bash
gitver release
```

Options:
- `--push`: Push the tag to remote repository

### Version

Display version information:

```bash
gitver version              # Show current version
gitver version --history    # Show version history
gitver version --history --count 5  # Show last 5 versions
```

### Changelog

Generate a changelog from commit messages:

```bash
gitver changelog                      # Generate changelog from latest tag to HEAD
gitver changelog --from v1.0.0        # Generate changelog from v1.0.0 to HEAD
gitver changelog --from v1.0.0 --to v2.0.0  # Generate changelog between two tags
gitver changelog --output CHANGELOG.md      # Save changelog to a file
gitver changelog --markdown                 # Format output as markdown
```

## Package Manager Integration

Gitver automatically updates version numbers in popular package manager files when bumping versions. This ensures that your project version is consistent across all files.

Supported package managers:

- **Maven**: Updates version in `pom.xml`
- **NPM**: Updates version in `package.json`
- **Gradle**: Updates version in `build.gradle`
- **Kotlin Gradle**: Updates version in `build.gradle.kts`
- **Go**: Updates version comment in `go.mod` and updates module path for major versions
- **Python**: Updates version in `setup.py` and `pyproject.toml` (supports both standard project and Poetry formats)

You can disable this feature using the `--update-packages=false` flag:

```bash
gitver bump --minor --update-packages=false
```

## Commit Message Format

Gitver follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for parsing commit messages when using the `--auto` flag:

- `feat: ...` or `feat(scope): ...`: A new feature (minor version bump)
- `fix: ...` or `fix(scope): ...`: A bug fix (patch version bump)
- `BREAKING CHANGE: ...`: A breaking API change (major version bump)

## Configuration

Gitver stores its configuration in a `.gitver` directory in your project root.

## License

MIT License
