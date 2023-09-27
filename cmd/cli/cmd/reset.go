package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"anti_bruteforce/internal/delivery/grpc/pb"
)

type resetFlags struct {
	login string
	ip    string
}

var resetF = resetFlags{}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete stats for given login or ip",
	Long:  `Delete stats for given login or ip`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resetF == (resetFlags{}) {
			return fmt.Errorf("at least one of the flags in the group [login, ip] is required")
		}

		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		in := &pb.ResetBucketRequest{
			Login: resetF.login,
			Ip:    resetF.ip,
		}

		_, err = client.ResetBucket(cmd.Context(), in)
		return err
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().StringVarP(&resetF.login, "login", "l", "", "Reset by login")
	resetCmd.Flags().StringVarP(&resetF.ip, "ip", "i", "", "Reset by ip")
}
