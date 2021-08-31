package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ferama/vipien/pkg/client"
	"github.com/ferama/vipien/pkg/iface"
	"github.com/ferama/vipien/pkg/tun"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client [ws|wss]://server:port",
	Short: "client",
	Long:  "client",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tun := tun.CreateTun()
		iface := iface.New(tun)
		iface.Setup("172.16.0.2/24")
		iface.AddRoute("172.21.0.0/16", "172.16.0.1")
		iface.SetupDns([]string{
			"8.8.8.8",
			"172.21.0.10",
		})

		client := client.New(args[0], iface)
		go client.Start()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		log.Println("exiting...")
		iface.SetupDns([]string{})
	},
}

func main() {
	clientCmd.Execute()
}
