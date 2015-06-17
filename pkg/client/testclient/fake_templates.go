package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
)

// FakeTemplates implements TemplateInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeTemplates struct {
	Fake      *Fake
	Namespace string
}

func (c *FakeTemplates) List(label labels.Selector, field fields.Selector) (*templateapi.TemplateList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-templates"}, &templateapi.TemplateList{})
	return obj.(*templateapi.TemplateList), err
}

func (c *FakeTemplates) Get(name string) (*templateapi.Template, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-template"}, &templateapi.Template{})
	return obj.(*templateapi.Template), err
}

func (c *FakeTemplates) Create(template *templateapi.Template) (*templateapi.Template, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-template", Value: template}, &templateapi.Template{})
	return obj.(*templateapi.Template), err
}

func (c *FakeTemplates) Update(template *templateapi.Template) (*templateapi.Template, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "update-template"}, &templateapi.Template{})
	return obj.(*templateapi.Template), err
}

func (c *FakeTemplates) Delete(name string) error {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "delete-template", Value: name})
	return nil
}

func (c *FakeTemplates) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "watch-templates"})
	return nil, nil
}
