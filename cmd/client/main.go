package main

import (
	"github.com/ferama/vipien/pkg/client"
	"github.com/ferama/vipien/pkg/tun"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client [ws|wss]://server:port",
	Short: "client",
	Long:  "client",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		iface := tun.CreateTun("172.16.0.2/24")

		client := client.New(args[0], iface)
		client.Start()
	},
}

func main() {
	clientCmd.Execute()
}
