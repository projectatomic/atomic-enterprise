package policy

import (
	"io"

	"github.com/spf13/cobra"

	adminpolicy "github.com/projectatomic/appinfra-next/pkg/cmd/admin/policy"
	cmdutil "github.com/projectatomic/appinfra-next/pkg/cmd/util"
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/clientcmd"
)

const PolicyRecommendedName = "policy"

func NewCmdPolicy(name, fullName string, f *clientcmd.Factory, out io.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   name,
		Short: "Manage authorization policy",
		Long:  `Manage authorization policy`,
		Run:   cmdutil.DefaultSubCommandRun(out),
	}

	cmds.AddCommand(adminpolicy.NewCmdWhoCan(adminpolicy.WhoCanRecommendedName, fullName+" "+adminpolicy.WhoCanRecommendedName, f, out))

	cmds.AddCommand(adminpolicy.NewCmdAddRoleToUser(adminpolicy.AddRoleToUserRecommendedName, fullName+" "+adminpolicy.AddRoleToUserRecommendedName, f, out))
	cmds.AddCommand(adminpolicy.NewCmdRemoveRoleFromUser(adminpolicy.RemoveRoleFromUserRecommendedName, fullName+" "+adminpolicy.RemoveRoleFromUserRecommendedName, f, out))
	cmds.AddCommand(adminpolicy.NewCmdRemoveUserFromProject(adminpolicy.RemoveUserRecommendedName, fullName+" "+adminpolicy.RemoveUserRecommendedName, f, out))
	cmds.AddCommand(adminpolicy.NewCmdAddRoleToGroup(adminpolicy.AddRoleToGroupRecommendedName, fullName+" "+adminpolicy.AddRoleToGroupRecommendedName, f, out))
	cmds.AddCommand(adminpolicy.NewCmdRemoveRoleFromGroup(adminpolicy.RemoveRoleFromGroupRecommendedName, fullName+" "+adminpolicy.RemoveRoleFromGroupRecommendedName, f, out))
	cmds.AddCommand(adminpolicy.NewCmdRemoveGroupFromProject(adminpolicy.RemoveGroupRecommendedName, fullName+" "+adminpolicy.RemoveGroupRecommendedName, f, out))

	return cmds
}
