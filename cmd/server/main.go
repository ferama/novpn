package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ferama/vipien/pkg/iface"
	"github.com/ferama/vipien/pkg/server"
	"github.com/ferama/vipien/pkg/tun"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  "server",
	Run: func(cmd *cobra.Command, args []string) {
		tun := tun.CreateTun()
		iface := iface.New(tun)
		iface.Setup("172.16.0.1/24")
		iface.Masquerade()

		server := server.New(iface)
		go server.Run("0.0.0.0:8800")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		log.Println("exiting...")
	},
}

func main() {
	serverCmd.Execute()
}
