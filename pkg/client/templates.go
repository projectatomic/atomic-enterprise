package client

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
)

// TemplatesNamespacer has methods to work with Template resources in a namespace
type TemplatesNamespacer interface {
	Templates(namespace string) TemplateInterface
}

// TemplateInterface exposes methods on Template resources.
type TemplateInterface interface {
	List(label labels.Selector, field fields.Selector) (*templateapi.TemplateList, error)
	Get(name string) (*templateapi.Template, error)
	Create(template *templateapi.Template) (*templateapi.Template, error)
	Update(template *templateapi.Template) (*templateapi.Template, error)
	Delete(name string) error
	Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

// templates implements TemplatesNamespacer interface
type templates struct {
	r  *Client
	ns string
}

// newTemplates returns a templates
func newTemplates(c *Client, namespace string) *templates {
	return &templates{
		r:  c,
		ns: namespace,
	}
}

// List returns a list of templates that match the label and field selectors.
func (c *templates) List(label labels.Selector, field fields.Selector) (result *templateapi.TemplateList, err error) {
	result = &templateapi.TemplateList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("templates").
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Do().
		Into(result)
	return
}

// Get returns information about a particular template and error if one occurs.
func (c *templates) Get(name string) (result *templateapi.Template, err error) {
	result = &templateapi.Template{}
	err = c.r.Get().Namespace(c.ns).Resource("templates").Name(name).Do().Into(result)
	return
}

// Create creates new template. Returns the server's representation of the template and error if one occurs.
func (c *templates) Create(template *templateapi.Template) (result *templateapi.Template, err error) {
	result = &templateapi.Template{}
	err = c.r.Post().Namespace(c.ns).Resource("templates").Body(template).Do().Into(result)
	return
}

// Update updates the template on server. Returns the server's representation of the template and error if one occurs.
func (c *templates) Update(template *templateapi.Template) (result *templateapi.Template, err error) {
	result = &templateapi.Template{}
	err = c.r.Put().Namespace(c.ns).Resource("templates").Name(template.Name).Body(template).Do().Into(result)
	return
}

// Delete deletes a template, returns error if one occurs.
func (c *templates) Delete(name string) (err error) {
	err = c.r.Delete().Namespace(c.ns).Resource("templates").Name(name).Do().Error()
	return
}

// Watch returns a watch.Interface that watches the requested templates
func (c *templates) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("templates").
		Param("resourceVersion", resourceVersion).
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Watch()
}
