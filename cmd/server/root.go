package main

import "github.com/spf13/cobra"

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "server",
		Short: "go-lang-template — a 12-factor Go service template",
	}
	root.AddCommand(newServeCmd())
	root.AddCommand(newMigrateCmd())
	return root
}
