# UpBank CLI

A command-line interface for Up Bank, allowing you to manage your finances from the terminal.

## Features

- List transactions with filtering options:
  - Filter by status (HELD, SETTLED)
  - Filter by date range (supports both YYYY-MM-DD and RFC3339 formats)
  - Filter by category
  - Filter by tag
  - Display transaction totals (debits, credits, and net balance)
- List accounts and their balances
- Raw mode output for scripting and automation

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/upbank-cli.git
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
export UP_API_TOKEN=your_api_token_here
```

### List Transactions
```bash
# List all transactions
./upbank-cli transactions

# List transactions with filtering options
./upbank-cli transactions --status SETTLED --since 2024-01-01 --until 2024-01-31

# List transactions in raw mode (without pretty formatting)
./upbank-cli transactions --raw
```

#### Raw Mode
The `--raw` flag outputs transactions in a format suitable for scripting and automation:
- No pretty formatting or colors
- No summary totals
- Includes transaction IDs
- Uses RFC3339 timestamp format
- Raw number values without thousand separators

Example usage in a script:
```bash
# Get transactions and process with awk
./upbank-cli transactions --raw | awk -F '|' '{print $3, $5}' > transactions.txt
```

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

#### Transaction Display
Transactions are displayed in a table format with the following information:
- Date and time
- Description
- Message
- Amount (with proper formatting)
- Currency
- Status
- Category
- Tags

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