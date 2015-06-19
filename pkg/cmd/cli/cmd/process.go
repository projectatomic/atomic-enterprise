package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/meta"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl"
	kcmdutil "github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd/util"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/resource"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/spf13/cobra"

	"github.com/projectatomic/appinfra-next/pkg/cmd/cli/describe"
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/clientcmd"
	"github.com/projectatomic/appinfra-next/pkg/template"
	"github.com/projectatomic/appinfra-next/pkg/template/api"
)

const (
	processLong = `Process template into a list of resources specified in filename or stdin

JSON and YAML formats are accepted.`

	processExample = `  // Convert template.json file into resource list
  $ %[1]s process -f template.json

  // Process template while passing a user-defined label
  $ %[1]s process -f template.json -l name=mytemplate

  // Convert stored template into resource list
  $ %[1]s process foo

  // Convert template.json into resource list
  $ cat template.json | %[1]s process -f -

  // Combine multiple templates into single resource list
  $ cat template.json second_template.json | %[1]s process -f -`
)

// NewCmdProcess implements the OpenShift cli process command
func NewCmdProcess(fullName string, f *clientcmd.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "process (TEMPLATE | -f FILENAME) [-v=KEY=VALUE]",
		Short:   "Process a template into list of resources",
		Long:    processLong,
		Example: fmt.Sprintf(processExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunProcess(f, out, cmd, args)
			kcmdutil.CheckErr(err)
		},
	}

	cmd.Flags().StringP("filename", "f", "", "Filename or URL to file to read a template")
	cmd.Flags().StringP("value", "v", "", "Specify a list of key-value pairs (eg. -v FOO=BAR,BAR=FOO) to set/override parameter values")
	cmd.Flags().BoolP("parameters", "", false, "Do not process but only print available parameters")
	cmd.Flags().StringP("labels", "l", "", "Label to set in all resources for this template")

	cmd.Flags().StringP("output", "o", "json", "Output format. One of: describe|json|yaml|template|templatefile.")
	cmd.Flags().Bool("raw", false, "If true output the processed template instead of the template's objects. Implied by -o describe")
	cmd.Flags().String("output-version", "", "Output the formatted object with the given version (default api-version).")
	cmd.Flags().StringP("template", "t", "", "Template string or path to template file to use when -o=template or -o=templatefile.  The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview]")
	return cmd
}

// RunProject contains all the necessary functionality for the OpenShift cli process command
func RunProcess(f *clientcmd.Factory, out io.Writer, cmd *cobra.Command, args []string) error {
	storedTemplate := ""
	if len(args) > 0 {
		storedTemplate = args[0]
	}

	filename := kcmdutil.GetFlagString(cmd, "filename")
	if len(storedTemplate) == 0 && len(filename) == 0 {
		return kcmdutil.UsageError(cmd, "Must pass a filename or name of stored template")
	}

	namespace, err := f.DefaultNamespace()
	if err != nil {
		return err
	}

	mapper, typer := f.Object()

	client, _, err := f.Clients()
	if err != nil {
		return err
	}

	var (
		objects []runtime.Object
		infos   []*resource.Info
		mapping *meta.RESTMapping
	)

	version, kind, err := mapper.VersionAndKindForResource("template")
	if mapping, err = mapper.RESTMapping(kind, version); err != nil {
		return err
	}

	// When storedTemplate is not empty, then we fetch the template from the
	// server, otherwise we require to set the `-f` parameter.
	if len(storedTemplate) > 0 {
		templateObj, err := client.Templates(namespace).Get(storedTemplate)
		if err != nil {
			if errors.IsNotFound(err) {
				return fmt.Errorf("template %q could not be found", storedTemplate)
			}
			return err
		}
		templateObj.CreationTimestamp = util.Now()
		infos = append(infos, &resource.Info{Object: templateObj})
	} else {
		infos, err = resource.NewBuilder(mapper, typer, f.ClientMapperForCommand()).
			NamespaceParam(namespace).RequireNamespace().
			FilenameParam(filename).
			Do().
			Infos()
		if err != nil {
			return err
		}
	}

	outputFormat := kcmdutil.GetFlagString(cmd, "output")

	for i := range infos {
		obj, ok := infos[i].Object.(*api.Template)
		if !ok {
			sourceName := filename
			if len(storedTemplate) > 0 {
				sourceName = namespace + "/" + storedTemplate
			}
			fmt.Fprintf(cmd.Out(), "unable to parse %q, not a valid Template but %s\n", sourceName, reflect.TypeOf(infos[i].Object))
			continue
		}

		// If 'parameters' flag is set it does not do processing but only print
		// the template parameters to console for inspection.
		// If multiple templates are passed, this will print combined output for all
		// templates.
		if kcmdutil.GetFlagBool(cmd, "parameters") {
			if len(infos) > 1 {
				fmt.Fprintf(out, "\n%s:\n", obj.Name)
			}
			if err := describe.PrintTemplateParameters(obj.Parameters, out); err != nil {
				fmt.Fprintf(cmd.Out(), "error printing parameters for %q: %v\n", obj.Name, err)
			}
			continue
		}

		if label := kcmdutil.GetFlagString(cmd, "labels"); len(label) > 0 {
			lbl, err := kubectl.ParseLabels(label)
			if err != nil {
				fmt.Fprintf(cmd.Out(), "error parsing labels: %v\n", err)
				continue
			}
			if obj.ObjectLabels == nil {
				obj.ObjectLabels = make(map[string]string)
			}
			for key, value := range lbl {
				obj.ObjectLabels[key] = value
			}
		}

		// Override the values for the current template parameters
		// when user specify the --value
		if cmd.Flag("value").Changed {
			injectUserVars(cmd, obj)
		}

		resultObj, err := client.TemplateConfigs(namespace).Create(obj)
		if err != nil {
			fmt.Fprintf(cmd.Out(), "error processing the template %q: %v\n", obj.Name, err)
			continue
		}

		if outputFormat == "describe" {
			if s, err := (&describe.TemplateDescriber{
				MetadataAccessor: meta.NewAccessor(),
				ObjectTyper:      kapi.Scheme,
				ObjectDescriber:  nil,
			}).DescribeTemplate(resultObj); err != nil {
				fmt.Fprintf(cmd.Out(), "error describing %q: %v\n", obj.Name, err)
			} else {
				fmt.Fprintf(out, s)
			}
			continue
		}
		objects = append(objects, resultObj.Objects...)
	}

	// Do not print the processed templates when asked to only show parameters or
	// describe.
	if kcmdutil.GetFlagBool(cmd, "parameters") || outputFormat == "describe" {
		return nil
	}

	p, _, err := kubectl.GetPrinter(outputFormat, "")
	if err != nil {
		return err
	}
	p = kubectl.NewVersionedPrinter(p, kapi.Scheme, kcmdutil.OutputVersion(cmd, mapping.APIVersion))

	// use generic output
	if kcmdutil.GetFlagBool(cmd, "raw") {
		for i := range objects {
			p.PrintObj(objects[i], out)
		}
		return nil
	}

	return p.PrintObj(&kapi.List{
		ListMeta: kapi.ListMeta{},
		Items:    objects,
	}, out)
}

// injectUserVars injects user specified variables into the Template
func injectUserVars(cmd *cobra.Command, t *api.Template) {
	values := util.StringList{}
	values.Set(kcmdutil.GetFlagString(cmd, "value"))
	for _, keypair := range values {
		p := strings.SplitN(keypair, "=", 2)
		if len(p) != 2 {
			fmt.Fprintf(cmd.Out(), "invalid parameter assignment in %q: %q\n", t.Name, keypair)
			continue
		}
		if v := template.GetParameterByName(t, p[0]); v != nil {
			v.Value = p[1]
			v.Generate = ""
			template.AddParameter(t, *v)
		} else {
			fmt.Fprintf(cmd.Out(), "unknown parameter name %q\n", p[0])
		}
	}
}
