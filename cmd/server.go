package cmd

import (
	"github.com/ferama/vipien/pkg/server"
	"github.com/ferama/vipien/pkg/tun"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  "server",
	Run: func(cmd *cobra.Command, args []string) {
		iface := tun.CreateTun("172.16.0.1/24")

		server := server.New(iface)
		server.Run("0.0.0.0:8800")
	},
}
