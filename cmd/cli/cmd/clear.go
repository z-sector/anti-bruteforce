package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear white and black list",
	Long:  `Clear white and black list`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGRPCClient(host)
		if err != nil {
			return err
		}

		_, err = client.ClearLists(cmd.Context(), &emptypb.Empty{})
		return err
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
