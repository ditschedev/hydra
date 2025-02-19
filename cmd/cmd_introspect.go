package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
)

func NewIntrospectCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "introspect",
		Short: "Introspect resources",
	}
	cmdx.RegisterHTTPClientFlags(cmd.PersistentFlags())
	cmdx.RegisterFormatFlags(cmd.PersistentFlags())
	return cmd
}
