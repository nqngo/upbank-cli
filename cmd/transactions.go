package cmd

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"upbank-cli/pkg/api"
	"upbank-cli/pkg/models"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// parseDateTime parses a date string that can be either a date (YYYY-MM-DD) or datetime (RFC3339)
// For date-only inputs, it sets the time to 00:00:00
func parseDateTime(input string) (time.Time, error) {
	// Try parsing as date first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", input); err == nil {
		// Set time to start of day
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), nil
	}

	// Try parsing as RFC3339
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format. Use YYYY-MM-DD or RFC3339 format (e.g. 2020-01-01T01:02:03+10:00)")
}

var (
	transactionsCmd = &cobra.Command{
		Use:   "transactions",
		Short: "List all transactions",
		Long:  `List all transactions with their details. Supports filtering by status, date range, category, and tag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			// Get flag values
			rawMode, _ := cmd.Flags().GetBool("raw")
			status, _ := cmd.Flags().GetString("status")
			since, _ := cmd.Flags().GetString("since")
			until, _ := cmd.Flags().GetString("until")
			category, _ := cmd.Flags().GetString("category")
			tag, _ := cmd.Flags().GetString("tag")

			// Build query parameters
			params := make(map[string]string)
			if status != "" {
				params["filter[status]"] = status
			}
			if since != "" {
				// Parse and format the since date
				sinceTime, err := parseDateTime(since)
				if err != nil {
					return fmt.Errorf("invalid since date: %v", err)
				}
				params["filter[since]"] = sinceTime.Format(time.RFC3339)
			}
			if until != "" {
				// Parse and format the until date
				untilTime, err := parseDateTime(until)
				if err != nil {
					return fmt.Errorf("invalid until date: %v", err)
				}
				params["filter[until]"] = untilTime.Format(time.RFC3339)
			}
			if category != "" {
				params["filter[category]"] = category
			}
			if tag != "" {
				params["filter[tag]"] = tag
			}

			// URL encode the parameters
			encodedParams := make(map[string]string)
			for k, v := range params {
				encodedParams[k] = url.QueryEscape(v)
			}

			transactions, err := client.GetTransactions(encodedParams)
			if err != nil {
				return err
			}

			// Sort transactions by date (newest first)
			sort.Sort(models.ByDate(transactions))

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())

			// Set header based on mode
			if rawMode {
				t.AppendHeader(table.Row{"ID", "Date", "Description", "Message", "Amount", "Currency", "Status", "Category", "Tags"})
			} else {
				t.AppendHeader(table.Row{"Date", "Description", "Message", "Amount", "Currency", "Status", "Category", "Tags"})
			}

			// Use built-in dark style
			if !rawMode {
				t.SetStyle(table.StyleColoredRedWhiteOnBlack)
			}

			// Create a new printer for number formatting
			p := message.NewPrinter(language.English)

			var totalDebit, totalCredit float64
			for _, tx := range transactions {
				amount, err := strconv.ParseFloat(tx.Attributes.Amount.Value, 64)
				if err != nil {
					return fmt.Errorf("error parsing amount: %v", err)
				}

				// Track debit and credit totals
				if amount < 0 {
					totalDebit += amount
				} else {
					totalCredit += amount
				}

				// Format amount with thousand separator unless raw mode
				formattedAmount := tx.Attributes.Amount.Value
				if !rawMode {
					// Convert to base units for proper formatting
					baseUnits := tx.Attributes.Amount.ValueInBaseUnits
					// Format with 2 decimal places
					formattedAmount = p.Sprintf("%.2f", float64(baseUnits)/100.0)
				}

				// Format date
				date := tx.Attributes.CreatedAt.Format(time.RFC3339)
				if !rawMode {
					date = tx.Attributes.CreatedAt.Format("Jan 02, 2006 15:04")
				}

				// Get category
				categoryName := ""
				if tx.Relations.Category.Data != nil {
					categoryName = tx.Relations.Category.Data.ID
				}

				// Get tags
				var tags []string
				for _, tag := range tx.Relations.Tags.Data {
					tags = append(tags, tag.ID)
				}
				tagsStr := strings.Join(tags, ", ")

				// Create row based on mode
				if rawMode {
					t.AppendRow(table.Row{
						tx.ID,
						date,
						tx.Attributes.Description,
						tx.Attributes.Message,
						formattedAmount,
						tx.Attributes.Amount.CurrencyCode,
						tx.Attributes.Status,
						categoryName,
						tagsStr,
					})
				} else {
					t.AppendRow(table.Row{
						date,
						tx.Attributes.Description,
						tx.Attributes.Message,
						formattedAmount,
						tx.Attributes.Amount.CurrencyCode,
						tx.Attributes.Status,
						categoryName,
						tagsStr,
					})
				}
			}

			t.AppendSeparator()
			// Format totals with thousand separator unless raw mode
			if !rawMode {
				formattedDebit := p.Sprintf("%.2f", totalDebit)
				formattedCredit := p.Sprintf("%.2f", totalCredit)
				t.AppendFooter(table.Row{
					"", "", "Debits 💸", formattedDebit, "AUD", "", "", "",
				})
				t.AppendFooter(table.Row{
					"", "", "Credits 💰", formattedCredit, "AUD", "", "", "",
				})
				t.AppendFooter(table.Row{
					"", "", "Net 🏦", p.Sprintf("%.2f", totalDebit+totalCredit), "AUD", "", "", "",
				})
			}

			t.Render()
			return nil
		},
	}
)

func init() {
	transactionsCmd.Flags().Bool("raw", false, "Display raw numbers without pretty formatting")
	transactionsCmd.Flags().String("status", "", "Filter transactions by status (HELD, SETTLED)")
	transactionsCmd.Flags().String("since", "", "Filter transactions from this date/time (format: YYYY-MM-DD or RFC3339 e.g. 2020-01-01T01:02:03+10:00). For date-only input, time will be set to 00:00:00")
	transactionsCmd.Flags().String("until", "", "Filter transactions until this date/time (format: YYYY-MM-DD or RFC3339 e.g. 2020-01-01T01:02:03+10:00). For date-only input, time will be set to 00:00:00")
	transactionsCmd.Flags().String("category", "", "Filter transactions by category ID")
	transactionsCmd.Flags().String("tag", "", "Filter transactions by tag ID")
	rootCmd.AddCommand(transactionsCmd)
}
