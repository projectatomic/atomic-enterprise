package openshift

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/cli"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/cli/cmd"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/experimental/buildchain"
	exipfailover "github.com/projectatomic/atomic-enterprise/pkg/cmd/experimental/ipfailover"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/experimental/tokens"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/flagtypes"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/infra/builder"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/infra/deployer"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/infra/gitserver"
	irouter "github.com/projectatomic/atomic-enterprise/pkg/cmd/infra/router"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/server/start"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/server/start/kubernetes"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/templates"
	cmdutil "github.com/projectatomic/atomic-enterprise/pkg/cmd/util"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/util/clientcmd"
	"github.com/projectatomic/atomic-enterprise/pkg/version"
)

const openshiftLong = `OpenShift Application Platform.

OpenShift helps you build, deploy, and manage your applications. To start an all-in-one server, run:

  $ openshift start &

OpenShift is built around Docker and the Kubernetes cluster container manager.  You must have
Docker installed on this machine to start your server.`

// CommandFor returns the appropriate command for this base name,
// or the global OpenShift command
func CommandFor(basename string) *cobra.Command {
	var cmd *cobra.Command

	out := os.Stdout

	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	switch basename {
	case "openshift-router":
		cmd = irouter.NewCommandTemplateRouter(basename)
	case "openshift-deploy":
		cmd = deployer.NewCommandDeployer(basename)
	case "openshift-sti-build":
		cmd = builder.NewCommandSTIBuilder(basename)
	case "openshift-docker-build":
		cmd = builder.NewCommandDockerBuilder(basename)
	case "openshift-gitserver":
		cmd = gitserver.NewCommandGitServer(basename)
	case "oc", "osc":
		cmd = cli.NewCommandCLI(basename, basename)
	case "oadm", "osadm":
		cmd = admin.NewCommandAdmin(basename, basename, out)
	case "kubectl":
		cmd = cli.NewCmdKubectl(basename, out)
	case "kube-apiserver":
		cmd = kubernetes.NewAPIServerCommand(basename, basename, out)
	case "kube-controller-manager":
		cmd = kubernetes.NewControllersCommand(basename, basename, out)
	case "kubelet":
		cmd = kubernetes.NewKubeletCommand(basename, basename, out)
	case "kube-proxy":
		cmd = kubernetes.NewProxyCommand(basename, basename, out)
	case "kube-scheduler":
		cmd = kubernetes.NewSchedulerCommand(basename, basename, out)
	case "kubernetes":
		cmd = kubernetes.NewCommand(basename, basename, out)
	case "origin":
		cmd = NewCommandOpenShift("origin")
	default:
		cmd = NewCommandOpenShift("openshift")
	}

	if cmd.UsageFunc() == nil {
		templates.ActsAsRootCommand(cmd)
	}
	flagtypes.GLog(cmd.PersistentFlags())

	return cmd
}

// NewCommandOpenShift creates the standard OpenShift command
func NewCommandOpenShift(name string) *cobra.Command {
	out := os.Stdout

	root := &cobra.Command{
		Use:   name,
		Short: "OpenShift helps you build, deploy, and manage your cloud applications",
		Long:  openshiftLong,
		Run:   cmdutil.DefaultSubCommandRun(out),
	}

	startAllInOne, _ := start.NewCommandStartAllInOne(name, out)
	root.AddCommand(startAllInOne)
	root.AddCommand(admin.NewCommandAdmin("admin", name+" admin", out))
	root.AddCommand(cli.NewCommandCLI("cli", name+" cli"))
	root.AddCommand(cli.NewCmdKubectl("kube", out))
	root.AddCommand(newExperimentalCommand("ex", name+" ex"))
	root.AddCommand(version.NewVersionCommand(name))

	// infra commands are those that are bundled with the binary but not displayed to end users
	// directly
	infra := &cobra.Command{
		Use: "infra", // Because this command exposes no description, it will not be shown in help
	}

	infra.AddCommand(
		irouter.NewCommandTemplateRouter("router"),
		deployer.NewCommandDeployer("deploy"),
		builder.NewCommandSTIBuilder("sti-build"),
		builder.NewCommandDockerBuilder("docker-build"),
		gitserver.NewCommandGitServer("git-server"),
	)
	root.AddCommand(infra)

	root.AddCommand(cmd.NewCmdOptions(out))

	// TODO: add groups
	templates.ActsAsRootCommand(root)

	return root
}

func newExperimentalCommand(name, fullName string) *cobra.Command {
	out := os.Stdout

	experimental := &cobra.Command{
		Use:   name,
		Short: "Experimental commands under active development",
		Long:  "The commands grouped here are under development and may change without notice.",
		Run: func(c *cobra.Command, args []string) {
			c.SetOutput(out)
			c.Help()
		},
	}

	f := clientcmd.New(experimental.PersistentFlags())

	experimental.AddCommand(tokens.NewCmdTokens(tokens.TokenRecommendedCommandName, fullName+" "+tokens.TokenRecommendedCommandName, f, out))
	experimental.AddCommand(exipfailover.NewCmdIPFailoverConfig(f, fullName, "ipfailover", out))
	experimental.AddCommand(buildchain.NewCmdBuildChain(f, fullName, "build-chain"))
	experimental.AddCommand(cmd.NewCmdOptions(out))
	return experimental
}
