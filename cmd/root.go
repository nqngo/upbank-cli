package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "upbank-cli",
	Short: "A CLI tool to interact with Upbank API",
	Long: `A CLI tool that allows you to interact with Upbank API to:
- List and retrieve account information
- Find and list transactions with total amounts`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
