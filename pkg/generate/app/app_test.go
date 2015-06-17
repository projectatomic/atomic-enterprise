package app

import (
	"log"
	"net/url"
	"reflect"
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	"github.com/projectatomic/appinfra-next/pkg/api/latest"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	imageapi "github.com/projectatomic/appinfra-next/pkg/image/api"
)

func testImageInfo() *imageapi.DockerImage {
	return &imageapi.DockerImage{
		Config: imageapi.DockerConfig{},
	}
}

func TestWithType(t *testing.T) {
	out := &Generated{
		Items: []runtime.Object{
			&buildapi.BuildConfig{
				ObjectMeta: kapi.ObjectMeta{
					Name: "foo",
				},
			},
			&kapi.Service{
				ObjectMeta: kapi.ObjectMeta{
					Name: "foo",
				},
			},
		},
	}

	builds := []buildapi.BuildConfig{}
	if !out.WithType(&builds) {
		t.Errorf("expected true")
	}
	if len(builds) != 1 {
		t.Errorf("unexpected slice: %#v", builds)
	}

	buildPtrs := []*buildapi.BuildConfig{}
	if out.WithType(&buildPtrs) {
		t.Errorf("expected false")
	}
	if len(buildPtrs) != 0 {
		t.Errorf("unexpected slice: %#v", buildPtrs)
	}
}

func TestBuildConfigNoOutput(t *testing.T) {
	url, err := url.Parse("https://github.com/projectatomic/appinfra-next.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	source := &SourceRef{URL: url}
	build := &BuildRef{Source: source}
	config, err := build.BuildConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Name != "appinfra-next" {
		t.Errorf("unexpected name: %#v", config)
	}
	if !reflect.DeepEqual(config.Parameters.Output, buildapi.BuildOutput{}) {
		t.Errorf("unexpected build output: %#v", config.Parameters.Output)
	}
}

func TestBuildConfigOutput(t *testing.T) {
	url, err := url.Parse("https://github.com/projectatomic/appinfra-next.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := &ImageRef{
		DockerImageReference: imageapi.DockerImageReference{
			Registry:  "myregistry",
			Namespace: "projectatomic",
			Name:      "appinfra-next",
		},
	}
	base := &ImageRef{
		DockerImageReference: imageapi.DockerImageReference{
			Namespace: "projectatomic",
			Name:      "ruby",
		},
		Info:          testImageInfo(),
		AsImageStream: true,
	}
	tests := []struct {
		asImageStream bool
		expectedKind  string
	}{
		{true, "ImageStreamTag"},
		{false, "DockerImage"},
	}
	for i, test := range tests {
		output.AsImageStream = test.asImageStream
		source := &SourceRef{URL: url}
		strategy := &BuildStrategyRef{IsDockerBuild: false, Base: base}
		build := &BuildRef{Source: source, Output: output, Strategy: strategy}
		config, err := build.BuildConfig()
		if err != nil {
			t.Fatalf("(%d) unexpected error: %v", i, err)
		}
		if config.Name != "appinfra-next" {
			t.Errorf("(%d) unexpected name: %s", i, config.Name)
		}
		if config.Parameters.Output.To.Name != "appinfra-next:latest" || config.Parameters.Output.To.Kind != test.expectedKind {
			t.Errorf("(%d) unexpected output image: %s/%s", i, config.Parameters.Output.To.Kind, config.Parameters.Output.To.Name)
		}
		if len(config.Triggers) != 3 {
			t.Errorf("(%d) unexpected number of triggers %d: %#v\n", i, len(config.Triggers), config.Triggers)
		}
		imageChangeTrigger := false
		for _, trigger := range config.Triggers {
			if trigger.Type == buildapi.ImageChangeBuildTriggerType {
				imageChangeTrigger = true
				if trigger.ImageChange == nil {
					t.Errorf("(%d) invalid image change trigger found: %#v", i, trigger)
				}
			}
		}
		if !imageChangeTrigger {
			t.Errorf("expecting image change trigger in build config")
		}
	}
}

func TestSimpleDeploymentConfig(t *testing.T) {
	image := &ImageRef{
		DockerImageReference: imageapi.DockerImageReference{
			Registry:  "myregistry",
			Namespace: "projectatomic",
			Name:      "appinfra-next",
		},
		Info:          testImageInfo(),
		AsImageStream: true,
	}
	deploy := &DeploymentConfigRef{Images: []*ImageRef{image}}
	config, err := deploy.DeploymentConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Name != "appinfra-next" || len(config.Triggers) != 2 || config.Template.ControllerTemplate.Template.Spec.Containers[0].Image != image.String() {
		t.Errorf("unexpected value: %#v", config)
	}
	for _, trigger := range config.Triggers {
		if trigger.Type == deployapi.DeploymentTriggerOnImageChange {
			from := trigger.ImageChangeParams.From
			if from.Kind != "ImageStream" {
				t.Errorf("unexpected from kind in image change trigger")
			}
			if from.Name != "appinfra-next" && from.Namespace != "projectatomic" {
				t.Errorf("unexpected  from name and namespace in image change trigger: %s, %s", from.Name, from.Namespace)
			}
		}
	}
}

func TestImageRefDeployableContainerPorts(t *testing.T) {
	tests := []struct {
		name          string
		inputPorts    map[string]struct{}
		expectedPorts map[int]string
		expectError   bool
	}{
		{
			name: "tcp implied, individual ports",
			inputPorts: map[string]struct{}{
				"123": {},
				"456": {},
			},
			expectedPorts: map[int]string{
				123: "TCP",
				456: "TCP",
			},
			expectError: false,
		},
		{
			name: "tcp implied, multiple ports",
			inputPorts: map[string]struct{}{
				"123 456":  {},
				"678 1123": {},
			},
			expectedPorts: map[int]string{
				123:  "TCP",
				678:  "TCP",
				456:  "TCP",
				1123: "TCP",
			},
			expectError: false,
		},
		{
			name: "tcp and udp, individual ports",
			inputPorts: map[string]struct{}{
				"123/tcp": {},
				"456/udp": {},
			},
			expectedPorts: map[int]string{
				123: "TCP",
				456: "UDP",
			},
			expectError: false,
		},
		{
			name: "tcp implied, multiple ports",
			inputPorts: map[string]struct{}{
				"123/tcp 456/udp":  {},
				"678/udp 1123/tcp": {},
			},
			expectedPorts: map[int]string{
				123:  "TCP",
				456:  "UDP",
				678:  "UDP",
				1123: "TCP",
			},
			expectError: false,
		},
		{
			name: "invalid format",
			inputPorts: map[string]struct{}{
				"123/tcp abc": {},
			},
			expectedPorts: map[int]string{},
			expectError:   true,
		},
	}
	for _, test := range tests {
		imageRef := &ImageRef{
			DockerImageReference: imageapi.DockerImageReference{
				Namespace: "test",
				Name:      "image",
				Tag:       imageapi.DefaultImageTag,
			},
			Info: &imageapi.DockerImage{
				Config: imageapi.DockerConfig{
					ExposedPorts: test.inputPorts,
				},
			},
		}
		container, _, err := imageRef.DeployableContainer()
		if err != nil && !test.expectError {
			t.Errorf("%s: unexpected error: %v", test.name, err)
			continue
		}
		if err == nil && test.expectError {
			t.Errorf("%s: got no error and expected an error", test.name)
			continue
		}
		if test.expectError {
			continue
		}
		remaining := test.expectedPorts
		for _, port := range container.Ports {
			proto, ok := remaining[port.ContainerPort]
			if !ok {
				t.Errorf("%s: got unexpected port: %v", test.name, port)
				continue
			}
			if kapi.Protocol(proto) != port.Protocol {
				t.Errorf("%s: got unexpected protocol %s for port %v", test.name, port.Protocol, port)
			}
			delete(remaining, port.ContainerPort)
		}
		if len(remaining) > 0 {
			t.Errorf("%s: did not find expected ports: %#v", test.name, remaining)
		}
	}
}

func TestSourceRefBuildSourceURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL without hash",
			input:    "https://github.com/projectatomic/ruby-hello-world.git",
			expected: "https://github.com/projectatomic/ruby-hello-world.git",
		},
		{
			name:     "URL with hash",
			input:    "https://github.com/projectatomic/ruby-hello-world.git#testref",
			expected: "https://github.com/projectatomic/ruby-hello-world.git",
		},
	}
	for _, tst := range tests {
		u, _ := url.Parse(tst.input)
		s := SourceRef{
			URL: u,
		}
		buildSource, _ := s.BuildSource()
		if buildSource.Git.URI != tst.expected {
			t.Errorf("%s: unexpected build source URI: %s. Expected: %s", tst.name, buildSource.Git.URI, tst.expected)
		}
	}
}

func ExampleGenerateSimpleDockerApp() {
	// TODO: determine if the repo is secured prior to fetching
	// TODO: determine whether we want to clone this repo, or use it directly. Using it directly would require setting hooks
	// if we have source, assume we are going to go into a build flow.
	// TODO: get info about git url: does this need STI?
	url, _ := url.Parse("https://github.com/projectatomic/appinfra-next.git")
	source := &SourceRef{URL: url}
	// generate a local name for the repo
	name, _ := source.SuggestName()
	// BUG: an image repo (if we want to create one) needs to tell other objects its pullspec, but we don't know what that will be
	// until the object is placed into a namespace and we lookup what registry (registries?) serve the object.
	// QUESTION: Is it ok for generation to require a namespace?  Do we want to be able to create apps with builds, image repos, and
	// deployment configs in templates (hint: yes).
	// SOLUTION? Make deployment config accept unqualified image repo names (foo) and then prior to creating the RC resolve those.
	output := &ImageRef{
		DockerImageReference: imageapi.DockerImageReference{
			Name: name,
		},
		AsImageStream: true,
	}
	// create our build based on source and input
	// TODO: we might need to pick a base image if this is STI
	build := &BuildRef{Source: source, Output: output}
	// take the output image and wire it into a deployment config
	deploy := &DeploymentConfigRef{Images: []*ImageRef{output}}

	outputRepo, _ := output.ImageStream()
	buildConfig, _ := build.BuildConfig()
	deployConfig, _ := deploy.DeploymentConfig()
	items := []runtime.Object{
		outputRepo,
		buildConfig,
		deployConfig,
	}
	out := &kapi.List{
		Items: items,
	}

	data, err := latest.Codec.Encode(out)
	if err != nil {
		log.Fatalf("Unable to generate output: %v", err)
	}
	log.Print(string(data))
	// output:
}
