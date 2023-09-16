package cmd //nolint:dupl

import (
	"github.com/spf13/cobra"

	"anti_bruteforce/internal/delivery/grpc/pb"
)

var blacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Actions on the black list",
	Long:  `Actions on the black list`,
}

var addToBlackListCmd = &cobra.Command{
	Use:   "add <subnet>",
	Short: "Add ip/mask to black list",
	Long:  `Add ip/mask to black list`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subnet := args[0]
		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		in := &pb.SubnetAddress{SubnetAddress: subnet}
		_, err = client.AddToBlackList(cmd.Context(), in)
		return err
	},
}

var removeFromBlackListCmd = &cobra.Command{
	Use:   "remove <subnet>",
	Short: "Remove ip/mask from black list",
	Long:  `Remove ip/mask from black list`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subnet := args[0]
		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		in := &pb.SubnetAddress{SubnetAddress: subnet}
		_, err = client.RemoveFromBlackList(cmd.Context(), in)
		return err
	},
}

func init() {
	rootCmd.AddCommand(blacklistCmd)
	blacklistCmd.AddCommand(addToBlackListCmd)
	blacklistCmd.AddCommand(removeFromBlackListCmd)
}
