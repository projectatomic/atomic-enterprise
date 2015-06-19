package client

import (
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"

	api "github.com/projectatomic/appinfra-next/pkg/build/api"
)

// BuildLogsNamespacer has methods to work with BuildLogs resources in a namespace
type BuildLogsNamespacer interface {
	BuildLogs(namespace string) BuildLogsInterface
}

// BuildLogsInterface exposes methods on BuildLogs resources.
type BuildLogsInterface interface {
	Get(name string, opts api.BuildLogOptions) *kclient.Request
}

// buildLogs implements BuildLogsNamespacer interface
type buildLogs struct {
	r  *Client
	ns string
}

// newBuildLogs returns a buildLogs
func newBuildLogs(c *Client, namespace string) *buildLogs {
	return &buildLogs{
		r:  c,
		ns: namespace,
	}
}

// Get builds and returns a buildLog request
func (c *buildLogs) Get(name string, opt api.BuildLogOptions) *kclient.Request {
	req := c.r.Get().Namespace(c.ns).Resource("builds").Name(name).SubResource("log")
	if opt.NoWait {
		req.Param("nowait", "true")
	}
	if opt.Follow {
		req.Param("follow", "true")
	}
	return req
}
