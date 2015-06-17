package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

// FakeBuilds implements BuildInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeBuilds struct {
	Fake      *Fake
	Namespace string
}

func (c *FakeBuilds) List(label labels.Selector, field fields.Selector) (*buildapi.BuildList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-builds"}, &buildapi.BuildList{})
	return obj.(*buildapi.BuildList), err
}

func (c *FakeBuilds) Get(name string) (*buildapi.Build, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-build"}, &buildapi.Build{})
	return obj.(*buildapi.Build), err
}

func (c *FakeBuilds) Create(build *buildapi.Build) (*buildapi.Build, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-build", Value: build}, &buildapi.Build{})
	return obj.(*buildapi.Build), err
}

func (c *FakeBuilds) Update(build *buildapi.Build) (*buildapi.Build, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "update-build"}, &buildapi.Build{})
	return obj.(*buildapi.Build), err
}

func (c *FakeBuilds) Delete(name string) error {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "delete-build", Value: name})
	return nil
}

func (c *FakeBuilds) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "watch-builds"})
	return nil, nil
}

func (c *FakeBuilds) Clone(request *buildapi.BuildRequest) (result *buildapi.Build, err error) {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "clone-build"})
	return nil, nil
}
