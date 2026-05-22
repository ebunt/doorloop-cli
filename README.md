# DoorLoop CLI

**Every DoorLoop resource at your fingertips, plus offline delinquency reports, vacancy maps, and cash flow snapshots no browser can match.**

doorloop-pp-cli syncs your entire DoorLoop portfolio to a local SQLite database, then unlocks compound queries no API call can answer: who owes rent and how to reach them, which units are vacant today, which leases expire this month, and net cash flow per property — all in under a second, offline, scriptable, and agent-ready.


## Install

The recommended path installs both the `doorloop-pp-cli` binary and the `pp-doorloop` agent skill (Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, and other agents supported by the upstream [`skills`](https://github.com/vercel-labs/skills) CLI) in one shot:

```bash
npx -y @mvanhorn/printing-press install doorloop
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press install doorloop --cli-only
```

For skill only — installs the skill into the same agents as the default command above, but skips the CLI binary (use this to update or reinstall just the skill):

```bash
npx -y @mvanhorn/printing-press install doorloop --skill-only
```

To constrain the skill install to one or more specific agents (repeatable — agent names match the [`skills`](https://github.com/vercel-labs/skills) CLI):

```bash
npx -y @mvanhorn/printing-press install doorloop --agent claude-code
npx -y @mvanhorn/printing-press install doorloop --agent claude-code --agent codex
```

### Without Node

The generated install path is category-agnostic until this CLI is published. If `npx` is not available before publish, install Node or use the category-specific Go fallback from the public-library entry after publish.

### Pre-built binary

Download a pre-built binary for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/doorloop-current). On macOS, clear the Gatekeeper quarantine: `xattr -d com.apple.quarantine <binary>`. On Unix, mark it executable: `chmod +x <binary>`.

<!-- pp-hermes-install-anchor -->
## Install for Hermes

From the Hermes CLI:

```bash
hermes skills install mvanhorn/printing-press-library/cli-skills/pp-doorloop --force
```

Inside a Hermes chat session:

```bash
/skills install mvanhorn/printing-press-library/cli-skills/pp-doorloop --force
```

## Install for OpenClaw

Tell your OpenClaw agent (copy this):

```
Install the pp-doorloop skill from https://github.com/mvanhorn/printing-press-library/tree/main/cli-skills/pp-doorloop. The skill defines how its required CLI can be installed.
```

## Use with Claude Desktop

This CLI ships an [MCPB](https://github.com/modelcontextprotocol/mcpb) bundle — Claude Desktop's standard format for one-click MCP extension installs (no JSON config required).

To install:

1. Download the `.mcpb` for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/doorloop-current).
2. Double-click the `.mcpb` file. Claude Desktop opens and walks you through the install.
3. Fill in `DOORLOOP_TOKEN` when Claude Desktop prompts you.

Requires Claude Desktop 1.0.0 or later. Pre-built bundles ship for macOS Apple Silicon (`darwin-arm64`) and Windows (`amd64`, `arm64`); for other platforms, use the manual config below.

<details>
<summary>Manual JSON config (advanced)</summary>

If you can't use the MCPB bundle (older Claude Desktop, unsupported platform), install the MCP binary and configure it manually.


Install the MCP binary from this CLI's published public-library entry or pre-built release.

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "doorloop": {
      "command": "doorloop-pp-mcp",
      "env": {
        "DOORLOOP_TOKEN": "<your-key>"
      }
    }
  }
}
```

</details>

## Authentication

Generate an API token in DoorLoop under Settings → API, then set DOORLOOP_TOKEN in your shell. All commands read this variable automatically.

## Quick Start

```bash
# Pull all properties, units, leases, tenants, payments, expenses, and tasks into local SQLite
doorloop-pp-cli sync --full


# See every tenant with an outstanding balance over $100 with contact info
doorloop-pp-cli delinquency --min-balance 100


# Find leases expiring in the next 30 days for renewal outreach
doorloop-pp-cli expiring --days 30


# Find units that have been vacant for over 30 days
doorloop-pp-cli vacancy --json | jq '.[] | select(.days_vacant > 30)'


# Get net cash flow per property for May
doorloop-pp-cli cashflow --from 2026-05-01 --to 2026-05-31

```

## Unique Features

These capabilities aren't available in any other tool for this API.

### Local state that compounds
- **`sync`** — Pull all DoorLoop entities into a local SQLite database for instant offline queries.

  _Every other novel command depends on synced local data. Run sync first, then query instantly._

  ```bash
  doorloop-pp-cli sync --full
  ```
- **`delinquency`** — See every tenant with an outstanding balance, with their contact info and days past due — one command instead of three clicks per tenant.

  _Use this to build a collections call queue or flag tenants for outreach without opening the web UI._

  ```bash
  doorloop-pp-cli delinquency --min-balance 100 --json --agent
  ```
- **`cashflow`** — Get income, expenses, and net cash flow per property for any date range — faster than navigating DoorLoop's PDF reports.

  _Pipe to jq to extract the net figure, or run for all properties to compare portfolio-wide cash position._

  ```bash
  doorloop-pp-cli cashflow --property prop_123 --from 2026-05-01 --to 2026-05-31 --json
  ```
- **`ledger`** — View a chronological bank-statement-style ledger for any lease: every charge, payment, credit, and returned payment with running balance.

  _Use when a tenant disputes their balance — instantly show the full charge and payment history with running balance._

  ```bash
  doorloop-pp-cli ledger lease_456 --json --agent
  ```
- **`vacancy`** — Find every vacant unit across your portfolio — units with no active lease — with days vacant and last lease end date.

  _Use for leasing pipeline planning — see which units need immediate outreach and how long they have been empty._

  ```bash
  doorloop-pp-cli vacancy --property prop_123 --json --agent
  ```

### Agent-native plumbing
- **`expiring`** — List all leases expiring within N days with tenant contact info and whether a prospect already exists for the unit.

  _Run weekly to build the renewal outreach list before leases lapse into month-to-month or vacancy._

  ```bash
  doorloop-pp-cli expiring --days 60 --json --agent
  ```
- **`tasks agenda`** — See today's tasks bucketed as OVERDUE, DUE TODAY, and DUE THIS WEEK — with property and assignee — ready to act on without opening a browser.

  _Use as a morning standup tool for maintenance teams, or pipe --json output to build a daily briefing._

  ```bash
  doorloop-pp-cli tasks agenda --property prop_123 --assignee user_789 --json
  ```
- **`portfolio-health`** — One-row-per-portfolio summary: total units, occupancy rate, total outstanding balance, and leases expiring within 30 days.

  _Use for the Monday portfolio review — spot which portfolios have rising delinquency or upcoming vacancy before problems compound._

  ```bash
  doorloop-pp-cli portfolio-health --json --agent
  ```

## Usage

Run `doorloop-pp-cli --help` for the full command reference and flag list.

## Commands

### accounts

Chart-of-accounts entries

- **`doorloop-pp-cli accounts get`** - 
- **`doorloop-pp-cli accounts list`** - 

### communication_logs

Logged communications with tenants and owners

- **`doorloop-pp-cli communication_logs create`** - 
- **`doorloop-pp-cli communication_logs delete`** - 
- **`doorloop-pp-cli communication_logs get`** - 
- **`doorloop-pp-cli communication_logs list`** - 
- **`doorloop-pp-cli communication_logs update`** - 

### expenses

Property operating expenses

- **`doorloop-pp-cli expenses create`** - 
- **`doorloop-pp-cli expenses delete`** - 
- **`doorloop-pp-cli expenses get`** - 
- **`doorloop-pp-cli expenses list`** - 
- **`doorloop-pp-cli expenses update`** - 

### files

Documents and file attachments

- **`doorloop-pp-cli files create`** - 
- **`doorloop-pp-cli files delete`** - 
- **`doorloop-pp-cli files get`** - 
- **`doorloop-pp-cli files list`** - 
- **`doorloop-pp-cli files update`** - 

### lease_charges

Charges applied to leases (rent, fees, etc.)

- **`doorloop-pp-cli lease_charges create`** - 
- **`doorloop-pp-cli lease_charges delete`** - 
- **`doorloop-pp-cli lease_charges get`** - 
- **`doorloop-pp-cli lease_charges list`** - 
- **`doorloop-pp-cli lease_charges update`** - 

### lease_credits

Credits applied to leases

- **`doorloop-pp-cli lease_credits create`** - 
- **`doorloop-pp-cli lease_credits delete`** - 
- **`doorloop-pp-cli lease_credits get`** - 
- **`doorloop-pp-cli lease_credits list`** - 
- **`doorloop-pp-cli lease_credits update`** - 

### lease_payments

Rent payments recorded against leases

- **`doorloop-pp-cli lease_payments create`** - 
- **`doorloop-pp-cli lease_payments delete`** - 
- **`doorloop-pp-cli lease_payments get`** - 
- **`doorloop-pp-cli lease_payments list`** - 
- **`doorloop-pp-cli lease_payments update`** - 

### lease_returned_payments

Returned/bounced payments on leases

- **`doorloop-pp-cli lease_returned_payments create`** - 
- **`doorloop-pp-cli lease_returned_payments delete`** - 
- **`doorloop-pp-cli lease_returned_payments get`** - 
- **`doorloop-pp-cli lease_returned_payments list`** - 
- **`doorloop-pp-cli lease_returned_payments update`** - 

### leases

Lease agreements linking tenants to units

- **`doorloop-pp-cli leases get`** - 
- **`doorloop-pp-cli leases list`** - 
- **`doorloop-pp-cli leases move-in`** - 
- **`doorloop-pp-cli leases move-out`** - 
- **`doorloop-pp-cli leases tenants`** - 

### notes

Notes attached to records

- **`doorloop-pp-cli notes create`** - 
- **`doorloop-pp-cli notes delete`** - 
- **`doorloop-pp-cli notes get`** - 
- **`doorloop-pp-cli notes list`** - 
- **`doorloop-pp-cli notes update`** - 

### owners

Property owners

- **`doorloop-pp-cli owners create`** - 
- **`doorloop-pp-cli owners delete`** - 
- **`doorloop-pp-cli owners get`** - 
- **`doorloop-pp-cli owners list`** - 
- **`doorloop-pp-cli owners update`** - 

### portfolios

Portfolio groupings of properties

- **`doorloop-pp-cli portfolios create`** - 
- **`doorloop-pp-cli portfolios delete`** - 
- **`doorloop-pp-cli portfolios get`** - 
- **`doorloop-pp-cli portfolios list`** - 
- **`doorloop-pp-cli portfolios update`** - 

### properties

Physical properties — residential or commercial

- **`doorloop-pp-cli properties get`** - 
- **`doorloop-pp-cli properties list`** - 

### reports

Generated financial and operational reports

- **`doorloop-pp-cli reports`** - 

### tasks

Maintenance and operational tasks

- **`doorloop-pp-cli tasks create`** - 
- **`doorloop-pp-cli tasks delete`** - 
- **`doorloop-pp-cli tasks get`** - 
- **`doorloop-pp-cli tasks list`** - 
- **`doorloop-pp-cli tasks post-update`** - 
- **`doorloop-pp-cli tasks update`** - 

### tenants

Tenants — active lease tenants and prospects

- **`doorloop-pp-cli tenants create`** - 
- **`doorloop-pp-cli tenants delete`** - 
- **`doorloop-pp-cli tenants get`** - 
- **`doorloop-pp-cli tenants list`** - 
- **`doorloop-pp-cli tenants update`** - 

### units

Individual rental units within properties

- **`doorloop-pp-cli units get`** - 
- **`doorloop-pp-cli units list`** - 

### users

DoorLoop users (staff)

- **`doorloop-pp-cli users create`** - 
- **`doorloop-pp-cli users delete`** - 
- **`doorloop-pp-cli users get`** - 
- **`doorloop-pp-cli users list`** - 
- **`doorloop-pp-cli users update`** - 

### vendor_bills

Bills from vendors against properties

- **`doorloop-pp-cli vendor_bills create`** - 
- **`doorloop-pp-cli vendor_bills delete`** - 
- **`doorloop-pp-cli vendor_bills get`** - 
- **`doorloop-pp-cli vendor_bills list`** - 
- **`doorloop-pp-cli vendor_bills update`** - 

### vendor_credits

Credits from vendors

- **`doorloop-pp-cli vendor_credits create`** - 
- **`doorloop-pp-cli vendor_credits delete`** - 
- **`doorloop-pp-cli vendor_credits get`** - 
- **`doorloop-pp-cli vendor_credits list`** - 
- **`doorloop-pp-cli vendor_credits update`** - 

### vendors

Vendors and contractors

- **`doorloop-pp-cli vendors create`** - 
- **`doorloop-pp-cli vendors delete`** - 
- **`doorloop-pp-cli vendors get`** - 
- **`doorloop-pp-cli vendors list`** - 
- **`doorloop-pp-cli vendors update`** - 


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
doorloop-pp-cli accounts list

# JSON for scripting and agents
doorloop-pp-cli accounts list --json

# Filter to specific fields
doorloop-pp-cli accounts list --json --select id,name,status

# Dry run — show the request without sending
doorloop-pp-cli accounts list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
doorloop-pp-cli accounts list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Explicit retries** - add `--idempotent` to create retries and `--ignore-missing` to delete retries when a no-op success is acceptable
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - write commands can accept structured input when their help lists `--stdin`
- **Offline-friendly** - sync/search commands can use the local SQLite store when available
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
doorloop-pp-cli doctor
```

Verifies configuration, credentials, and connectivity to the API.

## Configuration

Config file: ``

Static request headers can be configured under `headers`; per-command header overrides take precedence.

Environment variables:

| Name | Kind | Required | Description |
| --- | --- | --- | --- |
| `DOORLOOP_TOKEN` | per_call | Yes | Set to your API credential. |

## Troubleshooting
**Authentication errors (exit code 4)**
- Run `doorloop-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $DOORLOOP_TOKEN`
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

### API-specific

- **401 Unauthorized** — Run: export DOORLOOP_TOKEN=<your-token> — generate one at DoorLoop Settings → API
- **Empty results from delinquency, vacancy, or cashflow** — Run doorloop-pp-cli sync --full first to populate the local database
- **Stale data in offline commands** — Run doorloop-pp-cli sync to refresh; use --live flag to query the API directly

---

## Sources & Inspiration

This CLI was built by studying these projects and resources:

- [**mezbahalam/doorloop**](https://github.com/mezbahalam/doorloop) — Ruby (2 stars)
- [**cjmcfaul/doorloop**](https://github.com/cjmcfaul/doorloop) — Python (1 stars)

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
