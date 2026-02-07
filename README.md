# town

**T**eam **Own**ership — A CLI tool for exploring GitHub organizations, teams, and repository ownership via CODEOWNERS files.

## Features

- **List teams** in a GitHub organization
- **Find repositories** owned by a specific team (via CODEOWNERS)
- **Find repositories** without CODEOWNERS files
- **Clone repositories** in bulk
- **Smart caching** to minimize API calls
- **Shell autocompletion** for team names
- **Secure token storage** via system keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- **Embeddable** as a subcommand in other CLI tools

## Installation

```bash
go install github.com/lordzsolt/town@latest
```

Or build from source:

```bash
git clone https://github.com/lordzsolt/town.git
cd town
go build -o town .
```

## Quick Start

```bash
# List all teams in an organization
town teams --org myorg

# Find repos where a team is mentioned in CODEOWNERS
town repos --org myorg --team platform

# Find repos without CODEOWNERS
town repos --org myorg --no-owner

# Clone all repos owned by a team
town repos --org myorg --team platform --clone --clone-dir ~/projects
```

On first run, you'll be prompted for a GitHub token. The token is securely stored in your system keyring.

## GitHub Token

Create a [fine-grained personal access token](https://github.com/settings/personal-access-tokens/new) with:

- **Resource owner**: Your organization
- **Repository access**: All repositories
- **Permissions**:
  - Repository: Contents (read-only)
  - Organization: Members (read-only)

## Commands

### `town teams`

List all teams in an organization.

```bash
town teams --org myorg
```

Teams are cached locally to enable shell autocompletion for the `--team` flag.

### `town repos`

Find repositories based on CODEOWNERS.

```bash
# Find repos owned by a team
town repos --org myorg --team platform

# Find repos without CODEOWNERS
town repos --org myorg --no-owner

# Clone matching repos
town repos --org myorg --team platform --clone
town repos --org myorg --team platform --clone --clone-dir ~/work
```

Results are cached for 1 hour to avoid unnecessary API calls.

### `town completion`

Generate shell completion scripts.

```bash
# Bash
source <(town completion bash)

# Zsh (add to ~/.zshrc)
source <(town completion zsh)

# Fish
town completion fish | source
```

## Configuration

Town stores configuration in `~/.town/config.json` (or `$XDG_CONFIG_HOME/town/config.json`).

On first use with `--org`, a config file is automatically created with your default organization.

```json
{
  "default_org": "myorg",
  "default_team": "platform"
}
```

With a config file, you can omit flags:

```bash
# Instead of: town teams --org myorg
town teams

# Instead of: town repos --org myorg --team platform
town repos
```

## Caching

Town caches data to minimize GitHub API calls:

| Data | Location | TTL |
|------|----------|-----|
| Teams | `~/.town/cache/<org>/teams` | Until refreshed |
| Repos search | `~/.town/cache/<org>/repos-last.json` | 1 hour |

Delete the cache files to force a refresh.

## Embedding in Other CLIs

Town can be embedded as a subcommand in other Cobra-based CLI tools:

```go
import towncmd "github.com/lordzsolt/town/cmd"

func main() {
    rootCmd := &cobra.Command{Use: "mycli"}
    rootCmd.AddCommand(towncmd.Command())
    rootCmd.Execute()
}
```

Users can then run:

```bash
mycli town teams --org myorg
mycli town repos --team platform
```

## License

MIT
