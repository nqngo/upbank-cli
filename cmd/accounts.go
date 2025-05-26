package cmd

import (
	"fmt"
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

var (
	accountsCmd = &cobra.Command{
		Use:   "accounts",
		Short: "List all accounts",
		Long:  `List all accounts with their detail with optional filters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			// Get flag values
			rawMode, _ := cmd.Flags().GetBool("raw")
			accountType, _ := cmd.Flags().GetString("type")
			ownershipType, _ := cmd.Flags().GetString("ownership")

			// Build query parameters
			params := make(map[string]string)
			if accountType != "" {
				params["filter[accountType]"] = accountType
			}
			if ownershipType != "" {
				params["filter[ownershipType]"] = ownershipType
			}

			accounts, err := client.GetAccounts(params)
			if err != nil {
				return err
			}

			// Sort accounts by type and name
			sort.Sort(models.ByTypeAndName(accounts))

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())

			// Set header based on mode
			if rawMode {
				t.AppendHeader(table.Row{"ID", "Type", "Ownership", "Name", "Balance", "Currency", "Created At"})
			} else {
				t.AppendHeader(table.Row{"Type", "Ownership", "Name", "Balance", "Currency", "Created At"})
			}

			// Use built-in dark style
			if !rawMode {
				t.SetStyle(table.StyleColoredRedWhiteOnBlack)
			}

			// Create a new printer for number formatting
			p := message.NewPrinter(language.English)

			var totalBalance float64
			for _, account := range accounts {
				balance, err := strconv.ParseFloat(account.Attributes.Balance.Value, 64)
				if err != nil {
					return fmt.Errorf("error parsing balance: %v", err)
				}
				totalBalance += balance

				// Format balance with thousand separator unless raw mode
				formattedBalance := account.Attributes.Balance.Value
				if !rawMode {
					// Convert to base units for proper formatting
					baseUnits := account.Attributes.Balance.ValueInBaseUnits
					// Format with 2 decimal places
					formattedBalance = p.Sprintf("%.2f", float64(baseUnits)/100.0)
				}

				// Format creation date
				createdAt := account.Attributes.CreatedAt
				if !rawMode {
					if t, err := time.Parse(time.RFC3339, account.Attributes.CreatedAt); err == nil {
						createdAt = t.Format("Jan 02, 2006 15:04")
					}
				}

				// Create row based on mode
				if rawMode {
					t.AppendRow(table.Row{
						account.ID,
						account.Attributes.AccountType,
						account.Attributes.OwnershipType,
						account.Attributes.DisplayName,
						formattedBalance,
						account.Attributes.Balance.CurrencyCode,
						createdAt,
					})
				} else {
					t.AppendRow(table.Row{
						account.Attributes.AccountType,
						account.Attributes.OwnershipType,
						account.Attributes.DisplayName,
						formattedBalance,
						account.Attributes.Balance.CurrencyCode,
						createdAt,
					})
				}
			}

			t.AppendSeparator()
			// Format total with thousand separator unless raw mode
			if !rawMode {
				formattedTotal := p.Sprintf("%.2f", totalBalance)
				t.AppendFooter(table.Row{"", "", "Total", formattedTotal, "AUD", ""})
			}

			t.Render()
			return nil
		},
	}
)

func init() {
	accountsCmd.Flags().Bool("raw", false, "Display raw numbers without pretty formatting")
	accountsCmd.Flags().String("type", "", "Filter accounts by type (e.g., SAVER)")
	accountsCmd.Flags().String("ownership", "", "Filter accounts by ownership type (e.g., INDIVIDUAL)")
	rootCmd.AddCommand(accountsCmd)
}
