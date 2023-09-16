package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var host string

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Admin for anti-bruteforce",
	Long:  `Command-Line interface for administering the anti-bruteforce service`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(
		&host,
		"host",
		"localhost:9000",
		"gRPC server host",
	)
}
