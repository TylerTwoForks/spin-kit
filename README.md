# SPIN Kit
**S**andbox **P**rovisioning and **IN**spection Kit

A terminal UI for managing Salesforce sandbox provisioning, built with [Bubbletea](https://github.com/charmbracelet/bubbletea).

## Features

1. Save and manage frequently used sandbox names per production org.
2. Refresh sandboxes with a single keypress, or type in a custom sandbox name.
3. List all local org connections.
4. Reconnect to a sandbox after a refresh (re-authenticate via browser).
5. Check real-time sandbox refresh status.
6. Support for multiple production orgs with easy switching.
7. All configuration stored locally in SQLite (`~/.spin-kit/spin-kit.db`).

## Prerequisites

- [Salesforce CLI (`sf`)](https://developer.salesforce.com/tools/salesforcecli) installed and available in your PATH.
- Go 1.22+ (only needed if building from source).

## Install

### From source

```bash
go install github.com/TylerTwoForks/spin-kit/cmd/spin-kit@latest
```

### Build locally

```bash
git clone https://github.com/TylerTwoForks/spin-kit.git
cd spin-kit
go build -o spin-kit ./cmd/spin-kit
```

Move the binary somewhere in your `$PATH`:

```bash
mv spin-kit ~/scripts/    # or /usr/local/bin/, etc.
```

## Usage

```bash
spin-kit
```

On first launch, you will be prompted to add a production org alias. This is the alias (or username) of your production org that has sandbox refresh rights.

### Main Menu

From the main menu you can:

| Action | Description |
|---|---|
| Select a sandbox | Refresh that sandbox |
| List Org Connections | Run `sf org list` and display results |
| Connect to Prod Org | Login to a production org via browser |
| Refresh Custom Sandbox | Type a sandbox name and refresh it |
| Reconnect to Sandbox | Re-authenticate to a sandbox via browser |
| Sandbox Refresh Status | Query in-flight sandbox refresh progress |
| Manage Sandboxes | Add or remove saved sandbox names |
| Settings | Add, remove, or switch production orgs |

### Key Bindings

| Key | Action |
|---|---|
| `up`/`down`, `j`/`k` | Navigate |
| `enter` | Select / confirm |
| `s` | Jump to settings |
| `m` | Jump to settings |
| `a` | Add item (in management screens) |
| `d` | Delete item (in management screens) |
| `esc` | Go back |
| `q`, `ctrl+c` | Quit |

## Data Storage

All configuration is stored in a local SQLite database at `~/.spin-kit/spin-kit.db`. This includes:

- **Production orgs** -- aliases of your production orgs, with one marked as active.
- **Sandboxes** -- saved sandbox names, scoped to each production org.

No credentials are stored. Authentication is handled by the Salesforce CLI.

## Legacy Shell Script

The original `spin-kit.sh` shell script is still included in the repository for reference. The Go binary is the primary tool going forward.
