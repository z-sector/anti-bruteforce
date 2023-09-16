package cmd //nolint:dupl

import (
	"github.com/spf13/cobra"

	"anti_bruteforce/internal/delivery/grpc/pb"
)

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Actions on the white list",
	Long:  `Actions on the white list`,
}

var addToWhiteListCmd = &cobra.Command{
	Use:   "add <subnet>",
	Short: "Add ip/mask to white list",
	Long:  `Add ip/mask to white list`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subnet := args[0]
		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		in := &pb.SubnetAddress{SubnetAddress: subnet}
		_, err = client.AddToWhiteList(cmd.Context(), in)
		return err
	},
}

var removeFromWhiteListCmd = &cobra.Command{
	Use:   "remove <subnet>",
	Short: "Remove ip/mask from white list",
	Long:  `Remove ip/mask from white list`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subnet := args[0]
		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		in := &pb.SubnetAddress{SubnetAddress: subnet}
		_, err = client.RemoveFromWhiteList(cmd.Context(), in)
		return err
	},
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
	whitelistCmd.AddCommand(addToWhiteListCmd)
	whitelistCmd.AddCommand(removeFromWhiteListCmd)
}
