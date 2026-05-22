---
name: pp-doorloop
description: "Every DoorLoop resource at your fingertips, plus offline delinquency reports, vacancy maps, and cash flow snapshots... Trigger phrases: `check DoorLoop delinquency`, `which tenants owe rent`, `DoorLoop vacancy report`, `which units are vacant`, `lease expiring soon`, `cash flow for my properties`, `use DoorLoop`, `run DoorLoop`."
author: "Craig Collier"
license: "Apache-2.0"
argument-hint: "<command> [args] | install cli|mcp"
allowed-tools: "Read Bash"
metadata:
  openclaw:
    requires:
      bins:
        - doorloop-pp-cli
---

# DoorLoop — Printing Press CLI

## Prerequisites: Install the CLI

This skill drives the `doorloop-pp-cli` binary. **You must verify the CLI is installed before invoking any command from this skill.** If it is missing, install it first:

1. Install via the Printing Press installer:
   ```bash
   npx -y @mvanhorn/printing-press install doorloop --cli-only
   ```
2. Verify: `doorloop-pp-cli --version`
3. Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is on `$PATH`.

If the `npx` install fails before this CLI has a public-library category, install Node or use the category-specific Go fallback after publish.

If `--version` reports "command not found" after install, the install step did not put the binary on `$PATH`. Do not proceed with skill commands until verification succeeds.

This CLI syncs your entire DoorLoop portfolio to a local SQLite database, then unlocks compound queries no API call can answer: who owes rent and how to reach them, which units are vacant today, which leases expire this month, and net cash flow per property — all in under a second, offline, scriptable, and agent-ready.

## When to Use This CLI

Use doorloop-pp-cli when you need to query across multiple DoorLoop resources at once (delinquency + contacts, vacancy + lease history, cash flow across all properties), when you want scriptable or automated reporting without the web UI, or when an agent needs structured property management data to make decisions.

## Unique Capabilities

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

## Command Reference

**accounts** — Chart-of-accounts entries

- `doorloop-pp-cli accounts get` — 
- `doorloop-pp-cli accounts list` — 

**communication_logs** — Logged communications with tenants and owners

- `doorloop-pp-cli communication_logs create` — 
- `doorloop-pp-cli communication_logs delete` — 
- `doorloop-pp-cli communication_logs get` — 
- `doorloop-pp-cli communication_logs list` — 
- `doorloop-pp-cli communication_logs update` — 

**expenses** — Property operating expenses

- `doorloop-pp-cli expenses create` — 
- `doorloop-pp-cli expenses delete` — 
- `doorloop-pp-cli expenses get` — 
- `doorloop-pp-cli expenses list` — 
- `doorloop-pp-cli expenses update` — 

**files** — Documents and file attachments

- `doorloop-pp-cli files create` — 
- `doorloop-pp-cli files delete` — 
- `doorloop-pp-cli files get` — 
- `doorloop-pp-cli files list` — 
- `doorloop-pp-cli files update` — 

**lease_charges** — Charges applied to leases (rent, fees, etc.)

- `doorloop-pp-cli lease_charges create` — 
- `doorloop-pp-cli lease_charges delete` — 
- `doorloop-pp-cli lease_charges get` — 
- `doorloop-pp-cli lease_charges list` — 
- `doorloop-pp-cli lease_charges update` — 

**lease_credits** — Credits applied to leases

- `doorloop-pp-cli lease_credits create` — 
- `doorloop-pp-cli lease_credits delete` — 
- `doorloop-pp-cli lease_credits get` — 
- `doorloop-pp-cli lease_credits list` — 
- `doorloop-pp-cli lease_credits update` — 

**lease_payments** — Rent payments recorded against leases

- `doorloop-pp-cli lease_payments create` — 
- `doorloop-pp-cli lease_payments delete` — 
- `doorloop-pp-cli lease_payments get` — 
- `doorloop-pp-cli lease_payments list` — 
- `doorloop-pp-cli lease_payments update` — 

**lease_returned_payments** — Returned/bounced payments on leases

- `doorloop-pp-cli lease_returned_payments create` — 
- `doorloop-pp-cli lease_returned_payments delete` — 
- `doorloop-pp-cli lease_returned_payments get` — 
- `doorloop-pp-cli lease_returned_payments list` — 
- `doorloop-pp-cli lease_returned_payments update` — 

**leases** — Lease agreements linking tenants to units

- `doorloop-pp-cli leases get` — 
- `doorloop-pp-cli leases list` — 
- `doorloop-pp-cli leases move-in` — 
- `doorloop-pp-cli leases move-out` — 
- `doorloop-pp-cli leases tenants` — 

**notes** — Notes attached to records

- `doorloop-pp-cli notes create` — 
- `doorloop-pp-cli notes delete` — 
- `doorloop-pp-cli notes get` — 
- `doorloop-pp-cli notes list` — 
- `doorloop-pp-cli notes update` — 

**owners** — Property owners

- `doorloop-pp-cli owners create` — 
- `doorloop-pp-cli owners delete` — 
- `doorloop-pp-cli owners get` — 
- `doorloop-pp-cli owners list` — 
- `doorloop-pp-cli owners update` — 

**portfolios** — Portfolio groupings of properties

- `doorloop-pp-cli portfolios create` — 
- `doorloop-pp-cli portfolios delete` — 
- `doorloop-pp-cli portfolios get` — 
- `doorloop-pp-cli portfolios list` — 
- `doorloop-pp-cli portfolios update` — 

**properties** — Physical properties — residential or commercial

- `doorloop-pp-cli properties get` — 
- `doorloop-pp-cli properties list` — 

**reports** — Generated financial and operational reports

- `doorloop-pp-cli reports` — 

**tasks** — Maintenance and operational tasks

- `doorloop-pp-cli tasks create` — 
- `doorloop-pp-cli tasks delete` — 
- `doorloop-pp-cli tasks get` — 
- `doorloop-pp-cli tasks list` — 
- `doorloop-pp-cli tasks post-update` — 
- `doorloop-pp-cli tasks update` — 

**tenants** — Tenants — active lease tenants and prospects

- `doorloop-pp-cli tenants create` — 
- `doorloop-pp-cli tenants delete` — 
- `doorloop-pp-cli tenants get` — 
- `doorloop-pp-cli tenants list` — 
- `doorloop-pp-cli tenants update` — 

**units** — Individual rental units within properties

- `doorloop-pp-cli units get` — 
- `doorloop-pp-cli units list` — 

**users** — DoorLoop users (staff)

- `doorloop-pp-cli users create` — 
- `doorloop-pp-cli users delete` — 
- `doorloop-pp-cli users get` — 
- `doorloop-pp-cli users list` — 
- `doorloop-pp-cli users update` — 

**vendor_bills** — Bills from vendors against properties

- `doorloop-pp-cli vendor_bills create` — 
- `doorloop-pp-cli vendor_bills delete` — 
- `doorloop-pp-cli vendor_bills get` — 
- `doorloop-pp-cli vendor_bills list` — 
- `doorloop-pp-cli vendor_bills update` — 

**vendor_credits** — Credits from vendors

- `doorloop-pp-cli vendor_credits create` — 
- `doorloop-pp-cli vendor_credits delete` — 
- `doorloop-pp-cli vendor_credits get` — 
- `doorloop-pp-cli vendor_credits list` — 
- `doorloop-pp-cli vendor_credits update` — 

**vendors** — Vendors and contractors

- `doorloop-pp-cli vendors create` — 
- `doorloop-pp-cli vendors delete` — 
- `doorloop-pp-cli vendors get` — 
- `doorloop-pp-cli vendors list` — 
- `doorloop-pp-cli vendors update` — 


### Finding the right command

When you know what you want to do but not which command does it, ask the CLI directly:

```bash
doorloop-pp-cli which "<capability in your own words>"
```

`which` resolves a natural-language capability query to the best matching command from this CLI's curated feature index. Exit code `0` means at least one match; exit code `2` means no confident match — fall back to `--help` or use a narrower query.

## Recipes


### Monday delinquency call list

```bash
doorloop-pp-cli delinquency --min-balance 0 --json --agent --select tenant_name,phone,email,unit,property,outstanding_balance,days_past_due
```

Outputs just the fields your collections coordinator needs, ready to pipe into a CSV or email template.

### Weekly vacancy and expiry pipeline

```bash
doorloop-pp-cli vacancy --json && doorloop-pp-cli expiring --days 60 --json
```

Two commands cover Thursday's leasing pipeline check: vacant units and expiring leases in one terminal session.

### Per-lease tenant dispute ledger

```bash
doorloop-pp-cli ledger lease_456 --json --agent --select date,type,amount,running_balance
```

Shows a bank-statement-style charge+payment history with running balance for any lease — structured for agent use.

### Portfolio health for Monday standup

```bash
doorloop-pp-cli portfolio-health --json --agent --select portfolio_name,total_units,occupancy_rate,outstanding_balance,expiring_soon
```

One row per portfolio with occupancy and financial health — pipe to jq to flag portfolios below 90% occupancy.

### Friday cash flow snapshot

```bash
doorloop-pp-cli cashflow --from 2026-05-01 --to 2026-05-31 --json --agent --select property_name,income,expenses,net
```

Income minus expenses per property for the month — use --select to narrow the deeply nested response to just the numbers you need.

## Auth Setup

Generate an API token in DoorLoop under Settings → API, then set DOORLOOP_TOKEN in your shell. All commands read this variable automatically.

Run `doorloop-pp-cli doctor` to verify setup.

## Agent Mode

Add `--agent` to any command. Expands to: `--json --compact --no-input --no-color --yes`.

- **Pipeable** — JSON on stdout, errors on stderr
- **Filterable** — `--select` keeps a subset of fields. Dotted paths descend into nested structures; arrays traverse element-wise. Critical for keeping context small on verbose APIs:

  ```bash
  doorloop-pp-cli accounts list --agent --select id,name,status
  ```
- **Previewable** — `--dry-run` shows the request without sending
- **Offline-friendly** — sync/search commands can use the local SQLite store when available
- **Non-interactive** — never prompts, every input is a flag
- **Explicit retries** — use `--idempotent` only when an already-existing create should count as success, and `--ignore-missing` only when a missing delete target should count as success

### Response envelope

Commands that read from the local store or the API wrap output in a provenance envelope:

```json
{
  "meta": {"source": "live" | "local", "synced_at": "...", "reason": "..."},
  "results": <data>
}
```

Parse `.results` for data and `.meta.source` to know whether it's live or local. A human-readable `N results (live)` summary is printed to stderr only when stdout is a terminal AND no machine-format flag (`--json`, `--csv`, `--compact`, `--quiet`, `--plain`, `--select`) is set — piped/agent consumers and explicit-format runs get pure JSON on stdout.

## Agent Feedback

When you (or the agent) notice something off about this CLI, record it:

```
doorloop-pp-cli feedback "the --since flag is inclusive but docs say exclusive"
doorloop-pp-cli feedback --stdin < notes.txt
doorloop-pp-cli feedback list --json --limit 10
```

Entries are stored locally at `~/.doorloop-pp-cli/feedback.jsonl`. They are never POSTed unless `DOORLOOP_FEEDBACK_ENDPOINT` is set AND either `--send` is passed or `DOORLOOP_FEEDBACK_AUTO_SEND=true`. Default behavior is local-only.

Write what *surprised* you, not a bug report. Short, specific, one line: that is the part that compounds.

## Output Delivery

Every command accepts `--deliver <sink>`. The output goes to the named sink in addition to (or instead of) stdout, so agents can route command results without hand-piping. Three sinks are supported:

| Sink | Effect |
|------|--------|
| `stdout` | Default; write to stdout only |
| `file:<path>` | Atomically write output to `<path>` (tmp + rename) |
| `webhook:<url>` | POST the output body to the URL (`application/json` or `application/x-ndjson` when `--compact`) |

Unknown schemes are refused with a structured error naming the supported set. Webhook failures return non-zero and log the URL + HTTP status on stderr.

## Named Profiles

A profile is a saved set of flag values, reused across invocations. Use it when a scheduled agent calls the same command every run with the same configuration - HeyGen's "Beacon" pattern.

```
doorloop-pp-cli profile save briefing --json
doorloop-pp-cli --profile briefing accounts list
doorloop-pp-cli profile list --json
doorloop-pp-cli profile show briefing
doorloop-pp-cli profile delete briefing --yes
```

Explicit flags always win over profile values; profile values win over defaults. `agent-context` lists all available profiles under `available_profiles` so introspecting agents discover them at runtime.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage error (wrong arguments) |
| 3 | Resource not found |
| 4 | Authentication required |
| 5 | API error (upstream issue) |
| 7 | Rate limited (wait and retry) |
| 10 | Config error |

## Argument Parsing

Parse `$ARGUMENTS`:

1. **Empty, `help`, or `--help`** → show `doorloop-pp-cli --help` output
2. **Starts with `install`** → ends with `mcp` → MCP installation; otherwise → see Prerequisites above
3. **Anything else** → Direct Use (execute as CLI command with `--agent`)

## MCP Server Installation

Install the MCP binary from this CLI's published public-library entry or pre-built release, then register it:

```bash
claude mcp add doorloop-pp-mcp -- doorloop-pp-mcp
```

Verify: `claude mcp list`

## Direct Use

1. Check if installed: `which doorloop-pp-cli`
   If not found, offer to install (see Prerequisites at the top of this skill).
2. Match the user query to the best command from the Unique Capabilities and Command Reference above.
3. Execute with the `--agent` flag:
   ```bash
   doorloop-pp-cli <command> [subcommand] [args] --agent
   ```
4. If ambiguous, drill into subcommand help: `doorloop-pp-cli <command> --help`.
