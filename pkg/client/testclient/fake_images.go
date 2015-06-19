package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	"github.com/projectatomic/appinfra-next/pkg/client"
	imageapi "github.com/projectatomic/appinfra-next/pkg/image/api"
)

// FakeImages implements ImageInterface. Meant to be embedded into a struct to
// get a default implementation. This makes faking out just the methods you
// want to test easier.
type FakeImages struct {
	Fake *Fake
}

var _ client.ImageInterface = &FakeImages{}

func (c *FakeImages) List(label labels.Selector, field fields.Selector) (*imageapi.ImageList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-images"}, &imageapi.ImageList{})
	return obj.(*imageapi.ImageList), err
}

func (c *FakeImages) Get(name string) (*imageapi.Image, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-image", Value: name}, &imageapi.Image{})
	return obj.(*imageapi.Image), err
}

func (c *FakeImages) Create(image *imageapi.Image) (*imageapi.Image, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-image"}, &imageapi.Image{})
	return obj.(*imageapi.Image), err
}

func (c *FakeImages) Delete(name string) error {
	_, err := c.Fake.Invokes(FakeAction{Action: "delete-image", Value: name}, nil)
	return err
}
