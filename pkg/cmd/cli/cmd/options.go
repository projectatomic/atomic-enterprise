package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/projectatomic/atomic-enterprise/pkg/cmd/templates"
)

// NewCmdOptions implements the OpenShift cli options command
func NewCmdOptions(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "options",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	templates.UseOptionsTemplates(cmd)

	return cmd
}
