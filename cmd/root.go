package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd root cobra command.
var RootCmd = &cobra.Command{
	Use:   "shorthand <subcommand>",
	Short: "shorthand implements more readable k8s manifests",
}

var printCmd = &cobra.Command{
	Use:   "print [optional: v1 type name, e.g. Pod]",
	Short: "print all fields (recursively) of a k8s v1 type",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			loadAndPrint(args[0])
		} else {
			loadAndPrint("")
		}
	},
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "do some random stuff (development purposes only)",
	Run: func(cmd *cobra.Command, args []string) {
		genThings()
	},
}

func init() {
	RootCmd.AddCommand(printCmd, playCmd)
}
