package cmd

import (
	"github.com/ferama/vipien/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client [ws|wss]://server:port",
	Short: "client",
	Long:  "client",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := client.New(args[0])
		client.Start()
	},
}
