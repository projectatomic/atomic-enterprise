package autobuild

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	kclientcmd "github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/errors"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/generate/git"
	"github.com/projectatomic/appinfra-next/pkg/gitserver"
)

type AutoLinkBuilds struct {
	Namespaces []string
	Builders   []kapi.ObjectReference
	Client     client.BuildConfigsNamespacer

	CurrentNamespace string

	PostReceiveHook string

	LinkFn func(name string) *url.URL
}

var ErrNotEnabled = fmt.Errorf("not enabled")

func NewAutoLinkBuildsFromEnvironment() (*AutoLinkBuilds, error) {
	config := &AutoLinkBuilds{}

	file := os.Getenv("AUTOLINK_CONFIG")
	if len(file) == 0 {
		return nil, ErrNotEnabled
	}
	clientConfig, namespace, err := clientFromConfig(file)
	if err != nil {
		return nil, err
	}
	client, err := client.New(clientConfig)
	if err != nil {
		return nil, err
	}
	config.Client = client

	if value := os.Getenv("AUTOLINK_NAMESPACE"); len(value) > 0 {
		namespace = value
	}
	if len(namespace) == 0 {
		return nil, ErrNotEnabled
	}

	if value := os.Getenv("AUTOLINK_HOOK"); len(value) > 0 {
		abs, err := filepath.Abs(value)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(abs); err != nil {
			return nil, err
		}
		config.PostReceiveHook = abs
	}

	config.Namespaces = []string{namespace}
	config.CurrentNamespace = namespace
	return config, nil
}

func clientFromConfig(path string) (*kclient.Config, string, error) {
	rules := &kclientcmd.ClientConfigLoadingRules{ExplicitPath: path}
	credentials, err := rules.Load()
	if err != nil {
		return nil, "", fmt.Errorf("the provided credentials %q could not be loaded: %v", path, err)
	}
	cfg := kclientcmd.NewDefaultClientConfig(*credentials, &kclientcmd.ConfigOverrides{})
	config, err := cfg.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("the provided credentials %q could not be used: %v", path, err)
	}
	namespace, _ := cfg.Namespace()
	return config, namespace, nil
}

func (a *AutoLinkBuilds) Link() (map[string]gitserver.Clone, error) {
	errs := []error{}
	builders := []*buildapi.BuildConfig{}
	for _, namespace := range a.Namespaces {
		list, err := a.Client.BuildConfigs(namespace).List(labels.Everything(), fields.Everything())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for i := range list.Items {
			builders = append(builders, &list.Items[i])
		}
	}
	for _, b := range a.Builders {
		if hasItem(builders, b) {
			continue
		}
		config, err := a.Client.BuildConfigs(b.Namespace).Get(b.Name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		builders = append(builders, config)
	}

	hooks := make(map[string]string)
	if len(a.PostReceiveHook) > 0 {
		hooks["post-receive"] = a.PostReceiveHook
	}

	clones := make(map[string]gitserver.Clone)
	for _, builder := range builders {
		source := builder.Parameters.Source.Git
		if source == nil {
			continue
		}
		if builder.Annotations == nil {
			builder.Annotations = make(map[string]string)
		}

		// calculate the origin URL
		uri := source.URI
		if value, ok := builder.Annotations["git.openshift.io/origin-url"]; ok {
			uri = value
		}
		if len(uri) == 0 {
			continue
		}
		origin, err := git.ParseRepository(uri)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// calculate the local repository name and self URL
		name := builder.Name
		if a.CurrentNamespace != builder.Namespace {
			name = fmt.Sprintf("%s.%s", builder.Namespace, name)
		}
		name = fmt.Sprintf("%s.git", name)
		self := a.LinkFn(name)
		if self == nil {
			errs = append(errs, fmt.Errorf("no self URL available, can't update %s", name))
			continue
		}

		// we can't clone from ourself
		if self.Host == origin.Host {
			continue
		}

		// update the existing builder
		changed := false
		if builder.Annotations["git.openshift.io/origin-url"] != origin.String() {
			builder.Annotations["git.openshift.io/origin-url"] = origin.String()
			changed = true
		}
		if source.URI != self.String() {
			source.URI = self.String()
			changed = true
		}
		if changed {
			if _, err := a.Client.BuildConfigs(builder.Namespace).Update(builder); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		clones[name] = gitserver.Clone{
			URL:   *origin,
			Hooks: hooks,
		}
	}
	return clones, errors.NewAggregate(errs)
}

func hasItem(items []*buildapi.BuildConfig, item kapi.ObjectReference) bool {
	for _, c := range items {
		if c.Namespace == item.Namespace && c.Name == item.Name {
			return true
		}
	}
	return false
}
