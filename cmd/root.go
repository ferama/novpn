package cmd

import (
	"github.com/spf13/cobra"
)

// Version is the actual rospo version. This value
// is set during the build process using -ldflags="-X 'github.com/ferama/rospo/cmd.Version=
var Version = "development"

func init() {
}

var rootCmd = &cobra.Command{
	Use:     "vipien",
	Long:    "The vpn tool",
	Version: Version,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
