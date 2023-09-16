package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type resetFlags struct {
	login    string
	password string
	ip       string
}

var resetF = resetFlags{}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete stats for given login, password or ip",
	Long:  `Delete stats for given login, password or ip`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resetF == (resetFlags{}) {
			return fmt.Errorf("at least one of the flags in the group [login, password, ip] is required")
		}
		cmd.Println(resetF)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().StringVarP(&resetF.login, "login", "l", "", "Reset by login")
	resetCmd.Flags().StringVarP(&resetF.password, "password", "p", "", "Reset by password")
	resetCmd.Flags().StringVarP(&resetF.ip, "ip", "i", "", "Reset by ip")
	resetCmd.MarkFlagsRequiredTogether()
}
