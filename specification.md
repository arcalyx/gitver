# Gitver CLI Specification

## ✨ Purpose

Gitver is a command-line tool for managing semantic versioning according to SemVer 2.0.0. It integrates tightly with Git and various package managers, automating version bumps, tagging, changelog creation, and consistency validation.

---

## 📋 Command Overview

| Command      | Description                                    |
| ------------ | ---------------------------------------------- |
| `version`    | Display the current version or version history |
| `bump`       | Bump the version and perform related actions   |
| `config`     | Set or get configuration values                |
| `completion` | Install shell completion support               |
| `hooks`      | Manage Git hooks                               |

---

## 🔢 Command: `version`

### Description

Displays the current version of the project or the version history.

### Syntax

```bash
gitver version [--history] [--count N] [--format FORMAT]
```

### Options

| Option      | Type    | Default | Description                         |
| ----------- | ------- | ------- | ----------------------------------- |
| `--history` | Flag    | false   | Show full version history           |
| `--count`   | Integer | all     | Limit the number of history entries |
| `--format`  | Choice  | full    | Format: `semver`, `json`, or `full` |

### Examples

```bash
gitver version
gitver version --format semver
gitver version --history --count 5
```

---

## ➕ Command: `bump`

### Description

Bumps the current version, optionally appending metadata, generating changelogs, tagging, or syncing package files. This also updates subproject versions as defined in the configuration.

### Syntax

```bash
gitver bump [--major|--minor|--patch|--auto|--commit]
            [--prerelease <label|bump>]
            [--build-meta TEXT]
            [--tag] [--changelog] [--sync] [--dry-run]
```

### Options

| Option              | Type   | Description                                                  |
|---------------------|--------|--------------------------------------------------------------|
| `--major`           | Flag   | Increase MAJOR version (X.0.0)                               |
| `--minor`           | Flag   | Increase MINOR version (0.X.0)                               |
| `--patch`           | Flag   | Increase PATCH version (0.0.X)                               |
| `--auto`            | Flag   | Determine bump level from commit messages                    |
| `--infer`           | Flag   | Detect bump type using keyword-matching from config.yaml     |
| `--prerelease`      | String | Add or bump a pre-release identifier (e.g., alpha, rc, bump) |
| `--build-meta`      | String | Add build metadata (e.g., +build.123)                        |
| `--tag`             | Flag   | Create a Git tag for the version                             |
| `--changelog`       | Flag   | Generates a changelog entry based on the latest Git commits. Only commits following the Conventional Commits specification (semantic commit messages) will be included. The generated changelog can be saved to a file as defined in the configuration (config.yaml), and may be formatted in Markdown or plain text depending on settings.                    |
| `--update-packages` | Flag   | Update version in package manager files and subprojects      |
| `--dry-run`         | Flag   | Show the result without applying any changes                 |
| `--commit`          | Flag   | After performing the version bump, creates a Git commit that includes the modified version files. The commit message is taken from the configuration (config.yaml), allowing full customization of the commit content. This ensures version changes are properly tracked in Git history. |
| `--push`            | Flag   | After completing the bump and creating a Git commit (if --commit is used), this flag pushes all relevant changes—including commits and tags—to the remote repository.                |

### Validation Rules

* Only one of `--major`, `--minor`, `--patch`, `--auto`, or `--commit` may be used at once
* `--dry-run` disables `--tag`, `--sync`, and file changes

### Examples

```bash
gitver bump --patch

# With pre-release
gitver bump --minor --prerelease beta

# Auto bump with changelog and dry-run
gitver bump --auto --changelog --dry-run

# Bump based on configured commit keywords
gitver bump --infer
```

---

## ⚙️ Command: `config`

### Description

Sets, gets, or lists configuration values. Also used to initialize the configuration directory.

### Syntax

```bash
gitver config <key> <value>
gitver config --get <key>
gitver config --list
gitver config init [--detect-projects]
```

### Subcommand: `init`

Creates a `.gitver` directory in the project root and writes a `config.yaml` with default values.

#### Option:

* `--detect-projects`: Automatically detect subprojects in the Git repository and add them to the configuration file.

### Examples

```bash
gitver config init --detect-projects
```

### Configuration File: `.gitver/config.yaml`

Defines behavior and preferences for Gitver. Example structure:

```yaml
versioning:
  commit_keywords:
    major: ["#breaking", "BREAKING CHANGE"]
    minor: ["feat", "feature"]
    patch: ["fix", "bugfix"]
changelog:
  format: markdown
  output_file: CHANGELOG.md
subprojects:
  - path: libs/core/pom.xml
    type: maven
  - path: frontend/package.json
    type: npm
  - path: cli/pyproject.toml
    type: poetry
```

**Explanation:**

* `versioning.commit_keywords`: Defines commit message keywords that determine bump type when using `--commit`
* `changelog.format`: Output format, e.g., `markdown` or `plaintext`
* `changelog.output_file`: Target file for generated changelog
* `subprojects`: List of subprojects whose versions should also be managed

  * `path`: Relative path to the version file
  * `type`: Project type (e.g., `maven`, `npm`, `poetry`, etc.)

---

## ⌘ Command: `completion`

### Description

Installs shell completion for supported shells.

### Syntax

```bash
gitver completion <shell>
```

### Supported Shells

* bash
* zsh
* fish
* powershell

### Examples

```bash
gitver completion bash
gitver completion zsh
```

---

## ⚔️ Command: `hooks`

### Description

Installs or manages Git hooks for versioning workflows.

### Subcommands

| Subcommand  | Description                     |
| ----------- | ------------------------------- |
| `install`   | Install Git hooks               |
| `uninstall` | Remove all Gitver-managed hooks |
| `list`      | Show active Gitver hooks        |

### Examples

```bash
gitver hooks install
gitver hooks uninstall
gitver hooks list
```

---

## ✅ Code Quality Requirements

All source code for Gitver must:

* Be fully implemented in **English** (identifiers, comments, and documentation)
* Include **comprehensive unit and integration tests**
* Include **CLI integration tests** simulating real command-line usage scenarios
* Be thoroughly **commented** and **documented** using standard conventions
* Be maintained and versioned in a **GitHub repository**

### GitHub CI/CD Workflow

* **Main branch (`master`)** must automatically trigger a **release build**
* A **prerelease build** must be triggered when a commit contains the keyword `#prerelease`
* Builds should validate code, run tests, and publish release artifacts if applicable
* Final builds must be created and tested for **Windows**, **Linux**, and **macOS** platforms
* The logic to build and test should be in shell scripts in the scripts directory

---

## 🚀 Summary

This specification defines a simplified and consistent CLI interface for Gitver. It merges related functionalities under common commands and introduces flexible options to support diverse workflows while keeping the user experience intuitive. It also introduces a standard configuration mechanism via `.gitver/config.yaml` to enable project-level customization, including subproject support and version persistence. The implementation is required to follow strict quality standards, include integration testing for CLI commands, and leverage automated GitHub-based CI/CD workflows that deliver platform-independent builds.
