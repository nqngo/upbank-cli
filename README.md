# UpBank CLI

A command-line interface for Up Bank, allowing you to manage your finances from the terminal.

## Features

- List transactions with filtering options:
  - Filter by status (HELD, SETTLED)
  - Filter by date range (supports both YYYY-MM-DD and RFC3339 formats)
  - Filter by category
  - Filter by tag
  - Filter by foreign currency
  - Filter by description text (case-insensitive)
  - Display transaction totals (debits, credits, and net balance)
  - Multiple display modes (default, detail, raw)
- List accounts and their balances
- Raw mode output for scripting and automation

## Installation

### Quick Install (Linux/macOS)

You can install the CLI with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/nqngo/upbank-cli/refs/heads/main/install.sh | bash
```

### From Source

1. Clone the repository:
```bash
git clone https://github.com/nqngo/upbank-cli.git
cd upbank-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build
```

## Usage

Ensure you have a valid Up Bank API token. See [Up Bank API Documentation](https://developer.up.com.au/) for more information. To start:

```bash
export UPBANK_API_TOKEN=your_api_token_here
```

### List Transactions
```bash
# List all transactions
./upbank-cli transactions

# List transactions with filtering options
./upbank-cli transactions --status SETTLED --since 2024-01-01 --until 2024-01-31

# List transactions in raw mode (without pretty formatting)
./upbank-cli transactions --raw

# List transactions in detail mode (showing message, foreign amounts, and tags)
./upbank-cli transactions --detail

# Filter transactions by foreign currency
./upbank-cli transactions --currency JPY

# Filter transactions by description text
./upbank-cli transactions --description "Osaka"

# Combine multiple filters
./upbank-cli transactions --currency JPY --description "Osaka" --detail
```

#### Display Modes
The CLI supports three display modes for transactions:

1. Default Mode (no flags):
   - Shows essential information: date, description, amount, currency, category
   - Clean and concise output
   - Includes summary totals

2. Detail Mode (`--detail`):
   - Shows additional information: message, foreign amounts, tags
   - Useful for reviewing transaction details
   - Includes summary totals

3. Raw Mode (`--raw`):
   - Shows all available information
   - No pretty formatting or colors
   - Includes transaction IDs
   - Uses RFC3339 timestamp format
   - Raw number values without thousand separators
   - Suitable for scripting and automation

#### Transaction Filtering Options
- `--status`: Filter by transaction status (HELD, SETTLED)
- `--since`: Filter transactions from this date/time
  - Supports both date-only (YYYY-MM-DD) and full datetime (RFC3339) formats
  - Example: `--since 2024-01-01` or `--since "2024-01-01T00:00:00+10:00"`
  - For date-only input, time is automatically set to 00:00:00
- `--until`: Filter transactions until this date/time
  - Same format options as `--since`
- `--category`: Filter by category ID
- `--tag`: Filter by tag ID
- `--currency`: Filter by foreign currency code (e.g., JPY)
  - Client-side filter
  - Case-insensitive matching
  - Only shows transactions with foreign amounts in the specified currency
- `--description`: Filter by description text
  - Client-side filter
  - Case-insensitive partial matching
  - Matches any part of the description
  - Example: `--description "Osaka"` will match "Kids Plaza Osaka", "KAITEIROU SHINSAIBASH,OOSAKAFU", etc.

#### Transaction Display
Transactions are displayed in a table format with the following information:
- Date and time
- Description
- Message (in detail mode)
- Amount (with proper formatting)
- Currency
- Foreign amount and currency (in detail mode)
- Status (in raw mode)
- Category
- Tags (in detail mode)

The display includes summary totals at the bottom:
- 💸 Debits: Total of all negative transactions (money spent)
- 💰 Credits: Total of all positive transactions (money received)
- 🏦 Net: Overall balance (debits + credits)

### List Accounts
```bash
# List all accounts with pretty formatting
./upbank-cli accounts

# List accounts in raw mode (without pretty formatting)
./upbank-cli accounts --raw
```

#### Raw Mode
The `--raw` flag outputs accounts in a format suitable for scripting and automation:
- No pretty formatting or colors
- Includes account IDs
- Raw number values without thousand separators
- Uses RFC3339 timestamp format for dates

Example usage in a script:
```bash
# Get account balances and process with awk
./upbank-cli accounts --raw | awk -F '|' '{print $5, $6}' > balances.txt
```

#### Account Display
Accounts are displayed in a table format with the following information:
- Account name
- Account type
- Balance
- Currency
- Created date

## API Reference

This CLI uses the Up Bank API. For more information about the API endpoints and features, visit:
https://developer.up.com.au/ 
