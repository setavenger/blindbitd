package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	socketPath string
	RootCmd    = &cobra.Command{
		Use:   "blindbit-cli",
		Short: "A cli application to interact with the blindbit daemon",
		Long: `blindbit-cli is a CLI application to interact with the blindbit daemon:

It connects to the unix socket of blindbitd and users can manage their wallet through the cli.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// defines where to attach to the ipc socket
	RootCmd.PersistentFlags().StringVarP(&socketPath, "socket", "s", DefaultSocketPath, "Set the socket path. This is set to blindbitd default value")

	// required flags
	err := cobra.MarkFlagRequired(RootCmd.PersistentFlags(), "socket")
	if err != nil {
		log.Fatalln(err)
	}
}
