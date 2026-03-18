# Cursor Usage Analyzer

Interactive Go CLI to analyze Cursor team usage CSV exports with clear daily totals, model breakdowns, token share percentages, and cost summaries.

## Features

- Interactive CSV path prompt in terminal UI.
- Daily usage summary (sorted by date ascending).
- Model breakdown (sorted by total tokens descending, then model name ascending).
- Token share percentage per model.
- Overall tokens and overall cost.
- Friendly parsing for cost formats:
  - numeric cost values (for example `0.27`)
  - `$`-prefixed values (for example `$0.27`)
  - `-` (treated as zero)
  - `Free` (treated as zero)
- Path input normalization:
  - quoted paths
  - `~` home expansion
  - accidental line breaks in pasted paths

## Requirements

- Go `1.24.2+`

## Quick Start

Run directly with Go:

```bash
go run ./cmd/cursor-usage-analyzer
```

Build binary:

```bash
go build -o cursor-usage-analyzer ./cmd/cursor-usage-analyzer
./cursor-usage-analyzer
```

Install globally:

```bash
go install github.com/rennanbadaro/cursor-usage-analyzer/cmd/cursor-usage-analyzer@latest
```

Show help:

```bash
cursor-usage-analyzer --help
```

Show version:

```bash
cursor-usage-analyzer --version
```

## Pre-built Binaries

Download the appropriate binary for your platform from the `bin/` directory:

| Platform            | Binary                                |
|---------------------|---------------------------------------|
| macOS (Apple Silicon)| `bin/cursor-usage-analyzer-darwin-arm64` |
| macOS (Intel)       | `bin/cursor-usage-analyzer-darwin-amd64` |
| Linux (x86_64)      | `bin/cursor-usage-analyzer-linux-amd64`  |
| Linux (ARM64)       | `bin/cursor-usage-analyzer-linux-arm64`  |
| Windows (x86_64)    | `bin/cursor-usage-analyzer-windows-amd64.exe` |

After downloading, make it executable (macOS/Linux):

```bash
chmod +x bin/cursor-usage-analyzer-*
./bin/cursor-usage-analyzer-darwin-arm64
```

## Example Output

```text
Daily Usage
-----------
2026-02-18: 9,761,554 tokens | $9.30
2026-02-19: 25,238,533 tokens | $35.22

Model Breakdown
---------------
claude-4.6-opus-high-thinking: 26,893,134 tokens | 76.84% | $36.97
gpt-5.3-codex: 8,106,953 tokens | 23.16% | $7.55

Overall Tokens: 35,000,087
Overall Cost: $44.52
```

## Input Behavior

- Press `Enter` to submit the path.
- Press `Esc` or `Ctrl+C` to cancel.
- Returns a non-zero exit code on cancel, validation, parsing, or file/open errors.

## Testing

```bash
go test ./...
```
