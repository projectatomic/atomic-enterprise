package client

import (
	"fmt"

	"github.com/projectatomic/appinfra-next/pkg/image/api"
)

// ImageStreamTagsNamespacer has methods to work with ImageStreamTag resources in a namespace
type ImageStreamTagsNamespacer interface {
	ImageStreamTags(namespace string) ImageStreamTagInterface
}

// ImageStreamTagInterface exposes methods on ImageStreamTag resources.
type ImageStreamTagInterface interface {
	Get(name, tag string) (*api.ImageStreamTag, error)
	Delete(name, tag string) error
}

// imageStreamTags implements ImageStreamTagsNamespacer interface
type imageStreamTags struct {
	r  *Client
	ns string
}

// newImageStreamTags returns an imageStreamTags
func newImageStreamTags(c *Client, namespace string) *imageStreamTags {
	return &imageStreamTags{
		r:  c,
		ns: namespace,
	}
}

// Get finds the specified image by name of an image stream and tag.
func (c *imageStreamTags) Get(name, tag string) (result *api.ImageStreamTag, err error) {
	result = &api.ImageStreamTag{}
	err = c.r.Get().Namespace(c.ns).Resource("imageStreamTags").Name(fmt.Sprintf("%s:%s", name, tag)).Do().Into(result)
	return
}

// Delete deletes the specified tag from the image stream.
func (c *imageStreamTags) Delete(name, tag string) error {
	return c.r.Delete().Namespace(c.ns).Resource("imageStreamTags").Name(fmt.Sprintf("%s:%s", name, tag)).Do().Error()
}
