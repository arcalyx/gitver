# Gitver - Git Semantic Versioning Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/arcalyx/gitver)](https://goreportcard.com/report/github.com/arcalyx/gitver)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Gitver is a command-line tool that helps you manage semantic versioning in your Git repositories. It provides commands for bumping versions, creating releases, managing pre-release versions, and generating changelogs.

## Features

- **Semantic Versioning**: Follows [SemVer 2.0.0](https://semver.org/) specification
- **Version Bumping**: Easily bump major, minor, or patch versions
- **Auto Bump**: Automatically determine version bump based on commit messages
- **Pre-release Support**: Manage alpha, beta, and release candidate versions
- **Build Metadata Support**: Include build metadata in version strings (e.g., 1.2.3+build.123)
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

### Basic Commands

- `gitver init`: Initialize gitver in your project
- `gitver bump`: Bump the version of the project
  - `--major`: Bump the major version
  - `--minor`: Bump the minor version
  - `--patch`: Bump the patch version
  - `--auto`: Automatically determine version bump based on commit messages
- `gitver version`: Display version information

### Pre-release Commands

- `gitver prerelease set <type> <version>`: Set pre-release version (alpha, beta, rc)
- `gitver prerelease bump`: Bump pre-release version
- `gitver prerelease promote`: Promote pre-release to next level (alpha → beta → rc → release)
- `gitver prerelease clear`: Clear pre-release version

### Build Metadata Commands

- `gitver build-meta set <metadata>`: Set build metadata (e.g., build.123, sha.abc123)
- `gitver build-meta clear`: Clear build metadata
- `gitver build-meta show`: Display current build metadata

### Release Commands

- `gitver release`: Create a release tag for the current version

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

### Build Metadata

Gitver supports build metadata in version strings as per the SemVer 2.0.0 specification. Build metadata is specified by appending a plus sign and a series of dot-separated identifiers immediately following the patch or pre-release version.

Examples of valid version strings with build metadata:
```
1.2.3+build.123
1.2.3-alpha.1+build.456
1.2.3-beta.5+sha.abc123
```

Build metadata is supported in all version operations and is preserved when reading/writing versions. Note that build metadata does not affect version precedence - versions with the same version numbers but different build metadata are considered equal.

To set or modify build metadata, use the `version` command with the full version string:

```bash
gitver version 1.2.3+build.123
```

Alternatively, you can use the `build-meta` commands to set, clear, or display build metadata:

```bash
gitver build-meta set build.123
gitver build-meta clear
gitver build-meta show
```

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

Gitver can automatically update version numbers in various package manager files:

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

## Build Metadata in SemVer 2.0.0

According to the [SemVer 2.0.0 specification](https://semver.org/), build metadata is specified by appending a plus sign and a series of dot-separated identifiers immediately following the patch or pre-release version:

```
1.0.0+build.1
1.2.3-alpha.1+build.123
```

Build metadata is used to provide additional information about a build but does not affect version precedence. This means that `1.0.0+build.1` and `1.0.0+build.2` are considered equal in terms of precedence.

Gitver fully supports build metadata in accordance with the SemVer 2.0.0 specification:

1. Build metadata can be added to any version (regular or pre-release)
2. Build metadata is preserved during version operations but ignored in version comparisons
3. Build metadata identifiers can only contain alphanumeric characters and hyphens [0-9A-Za-z-]
4. Build metadata identifiers cannot be empty
5. Build metadata is dot-separated (e.g., `build.123`, `sha.abc123.date.20250716`)

## Git Hooks

Gitver can install Git hooks to automate version management tasks in your workflow:

```bash
gitver hooks --install-all                # Install all supported hooks
gitver hooks --pre-commit                 # Install only the pre-commit hook
gitver hooks --pre-commit --keyword=bump  # Install hook that triggers on "bump" keyword
gitver hooks --remove                     # Remove all gitver hooks
```

### Available Hooks

- **pre-commit**: Validates version consistency before commits
- **pre-push**: Ensures version is properly updated before pushing
- **post-merge**: Updates package versions after merging

### Keyword Triggers

You can specify a keyword that must be present in commit messages to trigger version-related actions:

```bash
gitver hooks --install-all --keyword=version
```

With this configuration:
- Only commits containing the word "version" will trigger version validation
- Only pushes with commit messages containing "version" will check for needed version bumps
- Only merges with commit messages containing "version" will automatically sync package versions

This allows you to control when version management is performed based on your commit messages.

### Supporting Commands

The hooks use these commands which you can also run manually:

```bash
gitver validate       # Check version consistency across package files
gitver sync-packages  # Update all package files to match the current version
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
