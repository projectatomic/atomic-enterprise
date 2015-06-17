package kubernetes

import (
	"net/url"

	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"

	"github.com/projectatomic/appinfra-next/pkg/util/httpproxy"
)

type ProxyConfig struct {
	ClientConfig *kclient.Config
}

func (c *ProxyConfig) InstallAPI(container *restful.Container) []string {
	kubeAddr, err := url.Parse(c.ClientConfig.Host)
	if err != nil {
		glog.Fatal(err)
	}

	proxy, err := httpproxy.NewUpgradeAwareSingleHostReverseProxy(c.ClientConfig, kubeAddr)
	if err != nil {
		glog.Fatalf("Unable to initialize the Kubernetes proxy: %v", err)
	}

	container.Handle("/api/", proxy)

	return []string{
		"Started Kubernetes proxy at %s/api/",
	}
}
