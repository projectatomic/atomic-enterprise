package origin

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/projectatomic/appinfra-next/pkg/assets"
	configapi "github.com/projectatomic/appinfra-next/pkg/cmd/server/api"
	cmdutil "github.com/projectatomic/appinfra-next/pkg/cmd/util"
	"github.com/projectatomic/appinfra-next/pkg/version"
)

// InstallAPI adds handlers for serving static assets into the provided mux,
// then returns an array of strings indicating what endpoints were started
// (these are format strings that will expect to be sent a single string value).
func (c *AssetConfig) InstallAPI(container *restful.Container) []string {
	if !c.OpenshiftEnabled {
		return []string{fmt.Sprintf("Openshift UI will not be started.")}
	}

	assetHandler, err := c.buildHandler()
	if err != nil {
		glog.Fatal(err)
	}

	publicURL, err := url.Parse(c.Options.PublicURL)
	if err != nil {
		glog.Fatal(err)
	}

	mux := container.ServeMux
	mux.Handle(publicURL.Path, http.StripPrefix(publicURL.Path, assetHandler))

	return []string{fmt.Sprintf("Started OpenShift UI %%s%s", publicURL.Path)}
}

// Run starts an http server for the static assets listening on the configured
// bind address
func (c *AssetConfig) Run() {
	if !c.OpenshiftEnabled {
		return
	}

	assetHandler, err := c.buildHandler()
	if err != nil {
		glog.Fatal(err)
	}

	publicURL, err := url.Parse(c.Options.PublicURL)
	if err != nil {
		glog.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(publicURL.Path, http.StripPrefix(publicURL.Path, assetHandler))

	if publicURL.Path != "/" {
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, publicURL.Path, http.StatusFound)
		})
	}

	timeout := c.Options.ServingInfo.RequestTimeoutSeconds
	if timeout == -1 {
		timeout = 0
	}

	server := &http.Server{
		Addr:           c.Options.ServingInfo.BindAddress,
		Handler:        mux,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	isTLS := configapi.UseTLS(c.Options.ServingInfo.ServingInfo)

	go util.Forever(func() {
		if isTLS {
			server.TLSConfig = &tls.Config{
				// Change default from SSLv3 to TLSv1.0 (because of POODLE vulnerability)
				MinVersion: tls.VersionTLS10,
			}
			glog.Infof("OpenShift UI listening at https://%s", c.Options.ServingInfo.BindAddress)
			glog.Fatal(server.ListenAndServeTLS(c.Options.ServingInfo.ServerCert.CertFile, c.Options.ServingInfo.ServerCert.KeyFile))
		} else {
			glog.Infof("OpenShift UI listening at http://%s", c.Options.ServingInfo.BindAddress)
			glog.Fatal(server.ListenAndServe())
		}
	}, 0)

	// Attempt to verify the server came up for 20 seconds (100 tries * 100ms, 100ms timeout per try)
	cmdutil.WaitForSuccessfulDial(isTLS, "tcp", c.Options.ServingInfo.BindAddress, 100*time.Millisecond, 100*time.Millisecond, 100)

	glog.Infof("OpenShift UI available at %s", c.Options.PublicURL)
}

func (c *AssetConfig) buildHandler() (http.Handler, error) {
	assets.RegisterMimeTypes()

	masterURL, err := url.Parse(c.Options.MasterPublicURL)
	if err != nil {
		return nil, err
	}

	publicURL, err := url.Parse(c.Options.PublicURL)
	if err != nil {
		glog.Fatal(err)
	}

	config := assets.WebConsoleConfig{
		MasterAddr:        masterURL.Host,
		MasterPrefix:      LegacyOpenShiftAPIPrefix, // TODO: change when the UI changes from v1beta3 to v1
		KubernetesAddr:    masterURL.Host,
		KubernetesPrefix:  KubernetesAPIPrefix,
		OAuthAuthorizeURI: OpenShiftOAuthAuthorizeURL(masterURL.String()),
		OAuthRedirectBase: c.Options.PublicURL,
		OAuthClientID:     OpenShiftWebConsoleClientID,
		LogoutURI:         c.Options.LogoutURL,
	}

	handler := http.FileServer(
		&assetfs.AssetFS{
			assets.Asset,
			assets.AssetDir,
			"",
		},
	)

	// Map of context roots (no leading or trailing slash) to the asset path to serve for requests to a missing asset
	subcontextMap := map[string]string{
		"":     "index.html",
		"java": "java/index.html",
	}

	handler, err = assets.HTML5ModeHandler(publicURL.Path, subcontextMap, handler)
	if err != nil {
		return nil, err
	}

	// Cache control should happen after all Vary headers are added, but before
	// any asset related routing (HTML5ModeHandler and FileServer)
	handler = assets.CacheControlHandler(version.Get().GitCommit, handler)

	// Generated config.js can not be cached since it changes depending on startup options
	handler = assets.GeneratedConfigHandler(config, handler)

	// Gzip first so that inner handlers can react to the addition of the Vary header
	handler = assets.GzipHandler(handler)

	return handler, nil
}
