# doorloop-pp-cli Usage Guide

## Build the CLI and MCP binaries
make build-all

## Install the binaries
make install      # Installs doorloop-pp-cli
make install-mcp  # Installs doorloop-pp-mcp

## Setup
Set your API token (generate in DoorLoop Settings → API):
```bash
export DOORLOOP_TOKEN="your_token_here"
```

## First Run
Sync your portfolio to a local SQLite database for instant, offline queries:
```bash
doorloop-pp-cli sync --full --db doorloop.db
```

## Quick Reference

| Task | Command |
| :--- | :--- |
| **Health Check** | `doorloop-pp-cli doctor` |
| **Delinquency** | `doorloop-pp-cli delinquency --min-balance 100` |
| **Expiring Leases** | `doorloop-pp-cli expiring --days 30` |
| **Vacant Units** | `doorloop-pp-cli vacancy` |
| **Cash Flow** | `doorloop-pp-cli cashflow --from 2026-05-01 --to 2026-05-31` |
| **Lease Ledger** | `doorloop-pp-cli ledger <lease_id>` |
| **Today's Tasks** | `doorloop-pp-cli tasks agenda` |
| **Portfolio Health** | `doorloop-pp-cli portfolio-health` |

## Formatting & Output
- **JSON:** Append `--json` for machine-readable output.
- **Select Fields:** Use `--select id,name` to limit output.
- **Agent Mode:** Use `--agent` for compact JSON (ideal for LLMs).
- **Dry Run:** Use `--dry-run` to see the API request without executing it.

## Troubleshooting
If data feels stale or commands return empty results, run:
```bash
doorloop-pp-cli sync
```
Use the `--live` flag on any command to bypass the local cache and query the API directly.
