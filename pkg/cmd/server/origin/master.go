package origin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/ww/goautoneg"

	etcdclient "github.com/coreos/go-etcd/etcd"
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/golang/glog"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kapierror "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	klatest "github.com/GoogleCloudPlatform/kubernetes/pkg/api/latest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	kmaster "github.com/GoogleCloudPlatform/kubernetes/pkg/master"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/service/allocator"
	etcdallocator "github.com/GoogleCloudPlatform/kubernetes/pkg/registry/service/allocator/etcd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/serviceaccount"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	serviceaccountadmission "github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/admission/serviceaccount"

	"github.com/openshift/origin/pkg/api/latest"
	"github.com/openshift/origin/pkg/api/v1"
	"github.com/openshift/origin/pkg/api/v1beta3"
	buildclient "github.com/openshift/origin/pkg/build/client"
	buildcontrollerfactory "github.com/openshift/origin/pkg/build/controller/factory"
	buildstrategy "github.com/openshift/origin/pkg/build/controller/strategy"
	buildgenerator "github.com/openshift/origin/pkg/build/generator"
	buildregistry "github.com/openshift/origin/pkg/build/registry/build"
	buildetcd "github.com/openshift/origin/pkg/build/registry/build/etcd"
	buildconfigregistry "github.com/openshift/origin/pkg/build/registry/buildconfig"
	buildconfigetcd "github.com/openshift/origin/pkg/build/registry/buildconfig/etcd"
	buildlogregistry "github.com/openshift/origin/pkg/build/registry/buildlog"
	"github.com/openshift/origin/pkg/build/webhook"
	"github.com/openshift/origin/pkg/build/webhook/generic"
	"github.com/openshift/origin/pkg/build/webhook/github"
	"github.com/openshift/origin/pkg/cmd/admin/policy"
	cmdutil "github.com/openshift/origin/pkg/cmd/util"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	configchangecontroller "github.com/openshift/origin/pkg/deploy/controller/configchange"
	deployerpodcontroller "github.com/openshift/origin/pkg/deploy/controller/deployerpod"
	deploycontroller "github.com/openshift/origin/pkg/deploy/controller/deployment"
	deployconfigcontroller "github.com/openshift/origin/pkg/deploy/controller/deploymentconfig"
	imagechangecontroller "github.com/openshift/origin/pkg/deploy/controller/imagechange"
	deployconfiggenerator "github.com/openshift/origin/pkg/deploy/generator"
	deployconfigregistry "github.com/openshift/origin/pkg/deploy/registry/deployconfig"
	deployconfigetcd "github.com/openshift/origin/pkg/deploy/registry/deployconfig/etcd"
	deployrollback "github.com/openshift/origin/pkg/deploy/registry/rollback"
	"github.com/openshift/origin/pkg/dns"
	imagecontroller "github.com/openshift/origin/pkg/image/controller"
	"github.com/openshift/origin/pkg/image/registry/image"
	imageetcd "github.com/openshift/origin/pkg/image/registry/image/etcd"
	"github.com/openshift/origin/pkg/image/registry/imagestream"
	imagestreametcd "github.com/openshift/origin/pkg/image/registry/imagestream/etcd"
	"github.com/openshift/origin/pkg/image/registry/imagestreamimage"
	"github.com/openshift/origin/pkg/image/registry/imagestreammapping"
	"github.com/openshift/origin/pkg/image/registry/imagestreamtag"
	accesstokenetcd "github.com/openshift/origin/pkg/oauth/registry/oauthaccesstoken/etcd"
	authorizetokenetcd "github.com/openshift/origin/pkg/oauth/registry/oauthauthorizetoken/etcd"
	clientetcd "github.com/openshift/origin/pkg/oauth/registry/oauthclient/etcd"
	clientauthetcd "github.com/openshift/origin/pkg/oauth/registry/oauthclientauthorization/etcd"
	projectcache "github.com/openshift/origin/pkg/project/cache"
	projectcontroller "github.com/openshift/origin/pkg/project/controller"
	projectproxy "github.com/openshift/origin/pkg/project/registry/project/proxy"
	projectrequeststorage "github.com/openshift/origin/pkg/project/registry/projectrequest/delegated"
	routeallocationcontroller "github.com/openshift/origin/pkg/route/controller/allocation"
	routeetcd "github.com/openshift/origin/pkg/route/registry/etcd"
	routeregistry "github.com/openshift/origin/pkg/route/registry/route"
	clusternetworketcd "github.com/openshift/origin/pkg/sdn/registry/clusternetwork/etcd"
	hostsubnetetcd "github.com/openshift/origin/pkg/sdn/registry/hostsubnet/etcd"
	securitycontroller "github.com/openshift/origin/pkg/security/controller"
	"github.com/openshift/origin/pkg/security/mcs"
	"github.com/openshift/origin/pkg/security/uid"
	"github.com/openshift/origin/pkg/security/uidallocator"
	"github.com/openshift/origin/pkg/service"
	templateregistry "github.com/openshift/origin/pkg/template/registry"
	templateetcd "github.com/openshift/origin/pkg/template/registry/etcd"
	identityregistry "github.com/openshift/origin/pkg/user/registry/identity"
	identityetcd "github.com/openshift/origin/pkg/user/registry/identity/etcd"
	userregistry "github.com/openshift/origin/pkg/user/registry/user"
	useretcd "github.com/openshift/origin/pkg/user/registry/user/etcd"
	"github.com/openshift/origin/pkg/user/registry/useridentitymapping"

	buildclonestorage "github.com/openshift/origin/pkg/build/registry/clone/generator"
	buildinstantiatestorage "github.com/openshift/origin/pkg/build/registry/instantiate/generator"

	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	clusterpolicyregistry "github.com/openshift/origin/pkg/authorization/registry/clusterpolicy"
	clusterpolicystorage "github.com/openshift/origin/pkg/authorization/registry/clusterpolicy/etcd"
	clusterpolicybindingregistry "github.com/openshift/origin/pkg/authorization/registry/clusterpolicybinding"
	clusterpolicybindingstorage "github.com/openshift/origin/pkg/authorization/registry/clusterpolicybinding/etcd"
	clusterrolestorage "github.com/openshift/origin/pkg/authorization/registry/clusterrole/proxy"
	clusterrolebindingstorage "github.com/openshift/origin/pkg/authorization/registry/clusterrolebinding/proxy"
	policyregistry "github.com/openshift/origin/pkg/authorization/registry/policy"
	policyetcd "github.com/openshift/origin/pkg/authorization/registry/policy/etcd"
	policybindingregistry "github.com/openshift/origin/pkg/authorization/registry/policybinding"
	policybindingetcd "github.com/openshift/origin/pkg/authorization/registry/policybinding/etcd"
	resourceaccessreviewregistry "github.com/openshift/origin/pkg/authorization/registry/resourceaccessreview"
	rolestorage "github.com/openshift/origin/pkg/authorization/registry/role/policybased"
	rolebindingstorage "github.com/openshift/origin/pkg/authorization/registry/rolebinding/policybased"
	"github.com/openshift/origin/pkg/authorization/registry/subjectaccessreview"
	"github.com/openshift/origin/pkg/cmd/server/admin"
	configapi "github.com/openshift/origin/pkg/cmd/server/api"
	"github.com/openshift/origin/pkg/cmd/server/bootstrappolicy"
	serviceaccountcontrollers "github.com/openshift/origin/pkg/serviceaccounts/controllers"
	"github.com/openshift/origin/plugins/osdn"
	routeplugin "github.com/openshift/origin/plugins/route/allocation/simple"
)

const (
	LegacyOpenShiftAPIPrefix  = "/osapi" // TODO: make configurable
	OpenShiftAPIPrefix        = "/oapi"  // TODO: make configurable
	KubernetesAPIPrefix       = "/api"   // TODO: make configurable
	OpenShiftAPIV1Beta3       = "v1beta3"
	OpenShiftAPIV1            = "v1"
	OpenShiftAPIPrefixV1Beta3 = LegacyOpenShiftAPIPrefix + "/" + OpenShiftAPIV1Beta3
	OpenShiftAPIPrefixV1      = OpenShiftAPIPrefix + "/" + OpenShiftAPIV1
	swaggerAPIPrefix          = "/swaggerapi/"
)

var (
	excludedV1Beta3Types = util.NewStringSet()
	excludedV1Types      = excludedV1Beta3Types

	// TODO: correctly solve identifying requests by type
	longRunningRE = regexp.MustCompile("watch|proxy|logs|exec|portforward")
)

// APIInstaller installs additional API components into this server
type APIInstaller interface {
	// InstallAPI returns an array of strings describing what was installed
	InstallAPI(*restful.Container) []string
}

// APIInstallFunc is a function for installing APIs
type APIInstallFunc func(*restful.Container) []string

// InstallAPI implements APIInstaller
func (fn APIInstallFunc) InstallAPI(container *restful.Container) []string {
	return fn(container)
}

func (c *MasterConfig) GetRestStorage() map[string]rest.Storage {
	defaultRegistry := env("OPENSHIFT_DEFAULT_REGISTRY", "${DOCKER_REGISTRY_SERVICE_HOST}:${DOCKER_REGISTRY_SERVICE_PORT}")
	svcCache := service.NewServiceResolverCache(c.KubeClient().Services(api.NamespaceDefault).Get)
	defaultRegistryFunc, err := svcCache.Defer(defaultRegistry)
	if err != nil {
		glog.Fatalf("OPENSHIFT_DEFAULT_REGISTRY variable is invalid %q: %v", defaultRegistry, err)
	}

	kubeletClient, err := kclient.NewKubeletClient(c.KubeletClientConfig)
	if err != nil {
		glog.Fatalf("Unable to configure Kubelet client: %v", err)
	}

	buildStorage := buildetcd.NewStorage(c.EtcdHelper)
	buildRegistry := buildregistry.NewRegistry(buildStorage)

	buildConfigStorage := buildconfigetcd.NewStorage(c.EtcdHelper)
	buildConfigRegistry := buildconfigregistry.NewRegistry(buildConfigStorage)

	deployConfigStorage := deployconfigetcd.NewStorage(c.EtcdHelper)
	deployConfigRegistry := deployconfigregistry.NewRegistry(deployConfigStorage)

	routeEtcd := routeetcd.New(c.EtcdHelper)
	hostSubnetStorage := hostsubnetetcd.NewREST(c.EtcdHelper)
	clusterNetworkStorage := clusternetworketcd.NewREST(c.EtcdHelper)

	userStorage := useretcd.NewREST(c.EtcdHelper)
	userRegistry := userregistry.NewRegistry(userStorage)
	identityStorage := identityetcd.NewREST(c.EtcdHelper)
	identityRegistry := identityregistry.NewRegistry(identityStorage)
	userIdentityMappingStorage := useridentitymapping.NewREST(userRegistry, identityRegistry)

	policyStorage := policyetcd.NewStorage(c.EtcdHelper)
	policyRegistry := policyregistry.NewRegistry(policyStorage)
	policyBindingStorage := policybindingetcd.NewStorage(c.EtcdHelper)
	policyBindingRegistry := policybindingregistry.NewRegistry(policyBindingStorage)

	clusterPolicyStorage := clusterpolicystorage.NewStorage(c.EtcdHelper)
	clusterPolicyRegistry := clusterpolicyregistry.NewRegistry(clusterPolicyStorage)
	clusterPolicyBindingStorage := clusterpolicybindingstorage.NewStorage(c.EtcdHelper)
	clusterPolicyBindingRegistry := clusterpolicybindingregistry.NewRegistry(clusterPolicyBindingStorage)

	roleStorage := rolestorage.NewVirtualStorage(policyRegistry)
	roleBindingStorage := rolebindingstorage.NewVirtualStorage(policyRegistry, policyBindingRegistry, clusterPolicyRegistry, clusterPolicyBindingRegistry)
	clusterRoleStorage := clusterrolestorage.NewClusterRoleStorage(clusterPolicyRegistry)
	clusterRoleBindingStorage := clusterrolebindingstorage.NewClusterRoleBindingStorage(clusterPolicyRegistry, clusterPolicyBindingRegistry)

	subjectAccessReviewStorage := subjectaccessreview.NewREST(c.Authorizer)
	subjectAccessReviewRegistry := subjectaccessreview.NewRegistry(subjectAccessReviewStorage)

	imageStorage := imageetcd.NewREST(c.EtcdHelper)
	imageRegistry := image.NewRegistry(imageStorage)
	imageStreamStorage, imageStreamStatusStorage := imagestreametcd.NewREST(c.EtcdHelper, imagestream.DefaultRegistryFunc(defaultRegistryFunc), subjectAccessReviewRegistry)
	imageStreamRegistry := imagestream.NewRegistry(imageStreamStorage, imageStreamStatusStorage)
	imageStreamMappingStorage := imagestreammapping.NewREST(imageRegistry, imageStreamRegistry)
	imageStreamTagStorage := imagestreamtag.NewREST(imageRegistry, imageStreamRegistry)
	imageStreamTagRegistry := imagestreamtag.NewRegistry(imageStreamTagStorage)
	imageStreamImageStorage := imagestreamimage.NewREST(imageRegistry, imageStreamRegistry)
	imageStreamImageRegistry := imagestreamimage.NewRegistry(imageStreamImageStorage)

	routeAllocator := c.RouteAllocator()

	buildGenerator := &buildgenerator.BuildGenerator{
		Client: buildgenerator.Client{
			GetBuildConfigFunc:      buildConfigRegistry.GetBuildConfig,
			UpdateBuildConfigFunc:   buildConfigRegistry.UpdateBuildConfig,
			GetBuildFunc:            buildRegistry.GetBuild,
			CreateBuildFunc:         buildRegistry.CreateBuild,
			GetImageStreamFunc:      imageStreamRegistry.GetImageStream,
			GetImageStreamImageFunc: imageStreamImageRegistry.GetImageStreamImage,
			GetImageStreamTagFunc:   imageStreamTagRegistry.GetImageStreamTag,
		},
		ServiceAccounts: c.KubeClient(),
		Secrets:         c.KubeClient(),
	}

	// TODO: with sharding, this needs to be changed
	deployConfigGenerator := &deployconfiggenerator.DeploymentConfigGenerator{
		Client: deployconfiggenerator.Client{
			DCFn:   deployConfigRegistry.GetDeploymentConfig,
			ISFn:   imageStreamRegistry.GetImageStream,
			LISFn2: imageStreamRegistry.ListImageStreams,
		},
	}
	_, kclient := c.DeploymentConfigControllerClients()
	deployRollback := &deployrollback.RollbackGenerator{}
	deployRollbackClient := deployrollback.Client{
		DCFn: deployConfigRegistry.GetDeploymentConfig,
		RCFn: clientDeploymentInterface{kclient}.GetDeployment,
		GRFn: deployRollback.GenerateRollback,
	}

	projectStorage := projectproxy.NewREST(kclient.Namespaces(), c.ProjectAuthorizationCache)

	namespace, templateName, err := configapi.ParseNamespaceAndName(c.Options.ProjectConfig.ProjectRequestTemplate)
	if err != nil {
		glog.Errorf("Error parsing project request template value: %v", err)
		// we can continue on, the storage that gets created will be valid, it simply won't work properly.  There's no reason to kill the master
	}
	projectRequestStorage := projectrequeststorage.NewREST(c.Options.ProjectConfig.ProjectRequestMessage, namespace, templateName, c.PrivilegedLoopbackOpenShiftClient)

	bcClient := c.BuildConfigWebHookClient()
	buildConfigWebHooks := buildconfigregistry.NewWebHookREST(
		buildConfigRegistry,
		buildclient.NewOSClientBuildConfigInstantiatorClient(bcClient),
		map[string]webhook.Plugin{
			"generic": generic.New(),
			"github":  github.New(),
		},
	)

	storage := map[string]rest.Storage{
		"builds":                   buildStorage,
		"buildConfigs":             buildConfigStorage,
		"buildConfigs/webhooks":    buildConfigWebHooks,
		"builds/clone":             buildclonestorage.NewStorage(buildGenerator),
		"buildConfigs/instantiate": buildinstantiatestorage.NewStorage(buildGenerator),
		"builds/log":               buildlogregistry.NewREST(buildRegistry, c.BuildLogClient(), kubeletClient),

		"images":              imageStorage,
		"imageStreams":        imageStreamStorage,
		"imageStreams/status": imageStreamStatusStorage,
		"imageStreamImages":   imageStreamImageStorage,
		"imageStreamMappings": imageStreamMappingStorage,
		"imageStreamTags":     imageStreamTagStorage,

		"deploymentConfigs":         deployConfigStorage,
		"generateDeploymentConfigs": deployconfiggenerator.NewREST(deployConfigGenerator, c.EtcdHelper.Codec),
		"deploymentConfigRollbacks": deployrollback.NewREST(deployRollbackClient, c.EtcdHelper.Codec),

		"processedTemplates": templateregistry.NewREST(),
		"templates":          templateetcd.NewREST(c.EtcdHelper),

		"routes": routeregistry.NewREST(routeEtcd, routeAllocator),

		"projects":        projectStorage,
		"projectRequests": projectRequestStorage,

		"hostSubnets":     hostSubnetStorage,
		"clusterNetworks": clusterNetworkStorage,

		"users":                userStorage,
		"identities":           identityStorage,
		"userIdentityMappings": userIdentityMappingStorage,

		"oAuthAuthorizeTokens":      authorizetokenetcd.NewREST(c.EtcdHelper),
		"oAuthAccessTokens":         accesstokenetcd.NewREST(c.EtcdHelper),
		"oAuthClients":              clientetcd.NewREST(c.EtcdHelper),
		"oAuthClientAuthorizations": clientauthetcd.NewREST(c.EtcdHelper),

		"resourceAccessReviews": resourceaccessreviewregistry.NewREST(c.Authorizer),
		"subjectAccessReviews":  subjectAccessReviewStorage,

		"policies":       policyStorage,
		"policyBindings": policyBindingStorage,
		"roles":          roleStorage,
		"roleBindings":   roleBindingStorage,

		"clusterPolicies":       clusterPolicyStorage,
		"clusterPolicyBindings": clusterPolicyBindingStorage,
		"clusterRoleBindings":   clusterRoleBindingStorage,
		"clusterRoles":          clusterRoleStorage,
	}

	return storage
}

func (c *MasterConfig) InstallProtectedAPI(container *restful.Container) []string {
	// initialize OpenShift API
	storage := c.GetRestStorage()

	messages := []string{}
	legacyAPIVersions := []string{}
	currentAPIVersions := []string{}

	if configapi.HasOpenShiftAPILevel(c.Options, OpenShiftAPIV1Beta3) {
		if err := c.api_v1beta3(storage).InstallREST(container); err != nil {
			glog.Fatalf("Unable to initialize v1beta3 API: %v", err)
		}
		messages = append(messages, fmt.Sprintf("Started OpenShift API at %%s%s", OpenShiftAPIPrefixV1Beta3))
		legacyAPIVersions = append(legacyAPIVersions, OpenShiftAPIV1Beta3)
	}

	if configapi.HasOpenShiftAPILevel(c.Options, OpenShiftAPIV1) {
		if err := c.api_v1(storage).InstallREST(container); err != nil {
			glog.Fatalf("Unable to initialize v1 API: %v", err)
		}
		messages = append(messages, fmt.Sprintf("Started OpenShift API at %%s%s (experimental)", OpenShiftAPIPrefixV1))
		currentAPIVersions = append(currentAPIVersions, OpenShiftAPIV1)
	}

	var root *restful.WebService
	for _, svc := range container.RegisteredWebServices() {
		switch svc.RootPath() {
		case "/":
			root = svc
		case OpenShiftAPIPrefixV1Beta3:
			svc.Doc("OpenShift REST API, version v1beta3").ApiVersion("v1beta3")
		case OpenShiftAPIPrefixV1:
			svc.Doc("OpenShift REST API, version v1").ApiVersion("v1")
		}
	}

	if root == nil {
		root = new(restful.WebService)
		container.Add(root)
	}

	initControllerRoutes(root, "/controllers", c.Options.Controllers != configapi.ControllersDisabled, c.ControllerPlug)
	initAPIVersionRoute(root, LegacyOpenShiftAPIPrefix, legacyAPIVersions...)
	initAPIVersionRoute(root, OpenShiftAPIPrefix, currentAPIVersions...)
	initReadinessCheckRoute(root, "healthz/ready", c.ProjectAuthorizationCache.ReadyForAccess)

	return messages
}

func (c *MasterConfig) InstallUnprotectedAPI(container *restful.Container) []string {
	return []string{}
}

// initAPIVersionRoute initializes the osapi endpoint to behave similar to the upstream api endpoint
func initAPIVersionRoute(root *restful.WebService, prefix string, versions ...string) {
	if len(versions) == 0 {
		return
	}

	versionHandler := apiserver.APIVersionHandler(versions...)
	root.Route(root.GET(prefix).To(versionHandler).
		Doc("list supported server API versions").
		Produces(restful.MIME_JSON).
		Consumes(restful.MIME_JSON))
}

// initReadinessCheckRoute initializes an HTTP endpoint for readiness checking
func initReadinessCheckRoute(root *restful.WebService, path string, readyFunc func() bool) {
	root.Route(root.GET(path).To(func(req *restful.Request, resp *restful.Response) {
		if readyFunc() {
			resp.ResponseWriter.WriteHeader(http.StatusOK)
			resp.ResponseWriter.Write([]byte("ok"))

		} else {
			resp.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
		}
	}).Doc("return the readiness state of OpenShift").
		Returns(http.StatusOK, "if OpenShift is ready", nil).
		Returns(http.StatusServiceUnavailable, "if OpenShift is not ready", nil).
		Produces(restful.MIME_JSON))
}

// If we know the location of the asset server, redirect to it when / is requested
// and the Accept header supports text/html
func assetServerRedirect(handler http.Handler, assetPublicURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			accepts := goautoneg.ParseAccept(req.Header.Get("Accept"))
			for _, accept := range accepts {
				if accept.Type == "text" && accept.SubType == "html" {
					http.Redirect(w, req, assetPublicURL, http.StatusFound)
					return
				}
			}
		}
		// Dispatch to the next handler
		handler.ServeHTTP(w, req)
	})
}

// assetServerOffNotice returns a notice that the web ui is off if the configuration boolean is
// isn't on
func assetServerOffNotice(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// List of paths to show notice on
		if req.URL.Path == "/login" || req.URL.Path == "/logout" || req.URL.Path == "/console" || strings.HasPrefix(req.URL.Path, "/console/") {
			w.Write([]byte("You need to upgrade to OpenShift in order to take advantage of this feature"))
			return
		}
		// Dispatch to the next handler
		handler.ServeHTTP(w, req)
	})
}

// TODO We would like to use the IndexHandler from k8s but we do not yet have a
// MuxHelper to track all registered paths
func indexAPIPaths(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			// TODO once we have a MuxHelper we will not need to hardcode this list of paths
			object := api.RootPaths{Paths: []string{
				"/api",
				"/api/v1beta3",
				"/api/v1",
				"/controllers",
				"/healthz",
				"/healthz/ping",
				"/logs/",
				"/metrics",
				"/ready",
				"/osapi",
				"/osapi/v1beta3",
				"/oapi",
				"/oapi/v1",
				"/swaggerapi/",
			}}
			// TODO it would be nice if apiserver.writeRawJSON was not private
			output, err := json.MarshalIndent(object, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", restful.MIME_JSON)
			w.WriteHeader(http.StatusOK)
			w.Write(output)
		} else {
			// Dispatch to the next handler
			handler.ServeHTTP(w, req)
		}
	})
}

// Run launches the OpenShift master. It takes optional installers that may install additional endpoints into the server.
// All endpoints get configured CORS behavior
// Protected installers' endpoints are protected by API authentication and authorization.
// Unprotected installers' endpoints do not have any additional protection added.
func (c *MasterConfig) Run(protected []APIInstaller, unprotected []APIInstaller) {
	var extra []string

	safe := kmaster.NewHandlerContainer(http.NewServeMux())
	open := kmaster.NewHandlerContainer(http.NewServeMux())

	// enforce authentication on protected endpoints
	protected = append(protected, APIInstallFunc(c.InstallProtectedAPI))
	for _, i := range protected {
		extra = append(extra, i.InstallAPI(safe)...)
	}
	handler := c.authorizationFilter(safe)
	handler = authenticationHandlerFilter(handler, c.Authenticator, c.getRequestContextMapper())
	handler = namespacingFilter(handler, c.getRequestContextMapper())
	handler = cacheControlFilter(handler, "no-store") // protected endpoints should not be cached

	// unprotected resources
	unprotected = append(unprotected, APIInstallFunc(c.InstallUnprotectedAPI))
	for _, i := range unprotected {
		extra = append(extra, i.InstallAPI(open)...)
	}

	handler = indexAPIPaths(handler)

	open.Handle("/", handler)

	// install swagger
	swaggerConfig := swagger.Config{
		WebServicesUrl:   c.Options.MasterPublicURL,
		WebServices:      append(safe.RegisteredWebServices(), open.RegisteredWebServices()...),
		ApiPath:          swaggerAPIPrefix,
		PostBuildHandler: customizeSwaggerDefinition,
	}
	// log nothing from swagger
	swagger.LogInfo = func(format string, v ...interface{}) {}
	swagger.RegisterSwaggerService(swaggerConfig, open)
	extra = append(extra, fmt.Sprintf("Started Swagger Schema API at %%s%s", swaggerAPIPrefix))

	handler = open

	// add CORS support
	if origins := c.ensureCORSAllowedOrigins(); len(origins) != 0 {
		handler = apiserver.CORS(handler, origins, nil, nil, "true")
	}

	if c.Options.AssetConfig != nil {
		if c.Options.OpenshiftEnabled {
			handler = assetServerRedirect(handler, c.Options.AssetConfig.PublicURL)
		} else {
			handler = assetServerOffNotice(handler)
		}
	}

	// Make the outermost filter the requestContextMapper to ensure all components share the same context
	if contextHandler, err := kapi.NewRequestContextFilter(c.getRequestContextMapper(), handler); err != nil {
		glog.Fatalf("Error setting up request context filter: %v", err)
	} else {
		handler = contextHandler
	}

	// TODO: MaxRequestsInFlight should be subdivided by intent, type of behavior, and speed of
	// execution - updates vs reads, long reads vs short reads, fat reads vs skinny reads.
	if c.Options.ServingInfo.MaxRequestsInFlight > 0 {
		sem := make(chan bool, c.Options.ServingInfo.MaxRequestsInFlight)
		handler = apiserver.MaxInFlightLimit(sem, longRunningRE, handler)
	}

	timeout := c.Options.ServingInfo.RequestTimeoutSeconds
	if timeout == -1 {
		timeout = 0
	}

	server := &http.Server{
		Addr:           c.Options.ServingInfo.BindAddress,
		Handler:        handler,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go util.Forever(func() {
		for _, s := range extra {
			glog.Infof(s, c.Options.ServingInfo.BindAddress)
		}
		if c.TLS {
			server.TLSConfig = &tls.Config{
				// Change default from SSLv3 to TLSv1.0 (because of POODLE vulnerability)
				MinVersion: tls.VersionTLS10,
				// Populate PeerCertificates in requests, but don't reject connections without certificates
				// This allows certificates to be validated by authenticators, while still allowing other auth types
				ClientAuth: tls.RequestClientCert,
				ClientCAs:  c.ClientCAs,
			}
			glog.Fatal(server.ListenAndServeTLS(c.Options.ServingInfo.ServerCert.CertFile, c.Options.ServingInfo.ServerCert.KeyFile))
		} else {
			glog.Fatal(server.ListenAndServe())
		}
	}, 0)

	// Attempt to verify the server came up for 20 seconds (100 tries * 100ms, 100ms timeout per try)
	cmdutil.WaitForSuccessfulDial(c.TLS, "tcp", c.Options.ServingInfo.BindAddress, 100*time.Millisecond, 100*time.Millisecond, 100)

	// Create required policy rules if needed
	c.ensureComponentAuthorizationRules()
	// Ensure the default SCCs are created
	c.ensureDefaultSecurityContextConstraints()
	// Bind default roles for service accounts in the default namespace if needed
	c.ensureDefaultNamespaceServiceAccountRoles()

	// Create the infra namespace
	c.ensureOpenShiftInfraNamespace()

	// Create the shared resource namespace
	c.ensureOpenShiftSharedResourcesNamespace()
}

func (c *MasterConfig) defaultAPIGroupVersion() *apiserver.APIGroupVersion {
	return &apiserver.APIGroupVersion{
		Root: OpenShiftAPIPrefix,

		Mapper: latest.RESTMapper,

		Creater:   kapi.Scheme,
		Typer:     kapi.Scheme,
		Convertor: kapi.Scheme,
		Linker:    latest.SelfLinker,

		Admit:   c.AdmissionControl,
		Context: c.getRequestContextMapper(),
	}
}

// api_v1beta3 returns the resources and codec for API version v1beta3.
func (c *MasterConfig) api_v1beta3(all map[string]rest.Storage) *apiserver.APIGroupVersion {
	storage := make(map[string]rest.Storage)
	for k, v := range all {
		if excludedV1Beta3Types.Has(k) {
			continue
		}
		storage[strings.ToLower(k)] = v
	}
	version := c.defaultAPIGroupVersion()
	version.Root = LegacyOpenShiftAPIPrefix
	version.Storage = storage
	version.Version = OpenShiftAPIV1Beta3
	version.Codec = v1beta3.Codec
	return version
}

// api_v1 returns the resources and codec for API version v1.
func (c *MasterConfig) api_v1(all map[string]rest.Storage) *apiserver.APIGroupVersion {
	storage := make(map[string]rest.Storage)
	for k, v := range all {
		if excludedV1Types.Has(k) {
			continue
		}
		storage[strings.ToLower(k)] = v
	}
	version := c.defaultAPIGroupVersion()
	version.Storage = storage
	version.Version = OpenShiftAPIV1
	version.Codec = v1.Codec
	return version
}

// getRequestContextMapper returns a mapper from requests to contexts, initializing it if needed
func (c *MasterConfig) getRequestContextMapper() kapi.RequestContextMapper {
	if c.RequestContextMapper == nil {
		c.RequestContextMapper = kapi.NewRequestContextMapper()
	}
	return c.RequestContextMapper
}

// ensureOpenShiftSharedResourcesNamespace is called as part of global policy initialization to ensure shared namespace exists
func (c *MasterConfig) ensureOpenShiftSharedResourcesNamespace() {
	if _, err := c.KubeClient().Namespaces().Get(c.Options.PolicyConfig.OpenShiftSharedResourcesNamespace); kapierror.IsNotFound(err) {
		namespace := &kapi.Namespace{
			ObjectMeta: kapi.ObjectMeta{Name: c.Options.PolicyConfig.OpenShiftSharedResourcesNamespace},
		}
		_, err = c.KubeClient().Namespaces().Create(namespace)
		if err != nil {
			glog.Errorf("Error creating namespace: %v due to %v\n", namespace, err)
		}
	}
}

// ensureOpenShiftInfraNamespace is called as part of global policy initialization to ensure infra namespace exists
func (c *MasterConfig) ensureOpenShiftInfraNamespace() {
	ns := c.Options.PolicyConfig.OpenShiftInfrastructureNamespace

	// Ensure namespace exists
	_, err := c.KubeClient().Namespaces().Create(&kapi.Namespace{ObjectMeta: kapi.ObjectMeta{Name: ns}})
	if err != nil && !kapierror.IsAlreadyExists(err) {
		glog.Errorf("Error creating namespace %s: %v", ns, err)
	}

	// Ensure service accounts exist
	serviceAccounts := []string{c.BuildControllerServiceAccount, c.DeploymentControllerServiceAccount, c.ReplicationControllerServiceAccount}
	for _, serviceAccountName := range serviceAccounts {
		_, err := c.KubeClient().ServiceAccounts(ns).Create(&kapi.ServiceAccount{ObjectMeta: kapi.ObjectMeta{Name: serviceAccountName}})
		if err != nil && !kapierror.IsAlreadyExists(err) {
			glog.Errorf("Error creating service account %s/%s: %v", ns, serviceAccountName, err)
		}
	}

	// Ensure service account cluster role bindings exist
	clusterRolesToUsernames := map[string][]string{
		bootstrappolicy.BuildControllerRoleName:       {serviceaccount.MakeUsername(ns, c.BuildControllerServiceAccount)},
		bootstrappolicy.DeploymentControllerRoleName:  {serviceaccount.MakeUsername(ns, c.DeploymentControllerServiceAccount)},
		bootstrappolicy.ReplicationControllerRoleName: {serviceaccount.MakeUsername(ns, c.ReplicationControllerServiceAccount)},
	}
	roleAccessor := policy.NewClusterRoleBindingAccessor(c.ServiceAccountRoleBindingClient())
	for clusterRole, usernames := range clusterRolesToUsernames {
		addRole := &policy.RoleModificationOptions{
			RoleName:            clusterRole,
			RoleBindingAccessor: roleAccessor,
			Users:               usernames,
		}
		if err := addRole.AddRole(); err != nil {
			glog.Errorf("Could not add %v users to the %v cluster role: %v\n", ns, usernames, clusterRole, err)
		} else {
			glog.V(2).Infof("Added %v users to the %v cluster role: %v\n", usernames, clusterRole, err)
		}
	}
}

// ensureComponentAuthorizationRules initializes the cluster policies
func (c *MasterConfig) ensureComponentAuthorizationRules() {
	clusterPolicyRegistry := clusterpolicyregistry.NewRegistry(clusterpolicystorage.NewStorage(c.EtcdHelper))
	ctx := kapi.WithNamespace(kapi.NewContext(), "")

	if _, err := clusterPolicyRegistry.GetClusterPolicy(ctx, authorizationapi.PolicyName); kapierror.IsNotFound(err) {
		glog.Infof("No cluster policy found.  Creating bootstrap policy based on: %v", c.Options.PolicyConfig.BootstrapPolicyFile)

		if err := admin.OverwriteBootstrapPolicy(c.EtcdHelper, c.Options.PolicyConfig.BootstrapPolicyFile, admin.CreateBootstrapPolicyFileFullCommand, true, ioutil.Discard); err != nil {
			glog.Errorf("Error creating bootstrap policy: %v", err)
		}

	} else {
		glog.V(2).Infof("Ignoring bootstrap policy file because cluster policy found")
	}
}

// ensureDefaultNamespaceServiceAccountRoles initializes roles for service accounts in the default namespace
func (c *MasterConfig) ensureDefaultNamespaceServiceAccountRoles() {
	const ServiceAccountRolesInitializedAnnotation = "openshift.io/sa.initialized-roles"

	// Wait for the default namespace
	var defaultNamespace *kapi.Namespace
	for i := 0; i < 30; i++ {
		ns, err := c.KubeClient().Namespaces().Get(kapi.NamespaceDefault)
		if err == nil {
			defaultNamespace = ns
			break
		}
		if kapierror.IsNotFound(err) {
			time.Sleep(time.Second)
			continue
		}
		glog.Errorf("Error adding service account roles to default namespace: %v", err)
		return
	}
	if defaultNamespace == nil {
		glog.Errorf("Default namespace not found, could not initialize default service account roles")
		return
	}

	// Short-circuit if we're already initialized
	if defaultNamespace.Annotations[ServiceAccountRolesInitializedAnnotation] == "true" {
		return
	}

	hasErrors := false
	for _, binding := range bootstrappolicy.GetBootstrapServiceAccountProjectRoleBindings(kapi.NamespaceDefault) {
		addRole := &policy.RoleModificationOptions{
			RoleName:            binding.RoleRef.Name,
			RoleNamespace:       binding.RoleRef.Namespace,
			RoleBindingAccessor: policy.NewLocalRoleBindingAccessor(kapi.NamespaceDefault, c.ServiceAccountRoleBindingClient()),
			Users:               binding.Users.List(),
			Groups:              binding.Groups.List(),
		}
		if err := addRole.AddRole(); err != nil {
			glog.Errorf("Could not add service accounts to the %v role in the %v namespace: %v\n", binding.RoleRef.Name, kapi.NamespaceDefault, err)
			hasErrors = true
		}
	}

	// If we had errors, don't register initialization so we can try again
	if !hasErrors {
		if defaultNamespace.Annotations == nil {
			defaultNamespace.Annotations = map[string]string{}
		}
		defaultNamespace.Annotations[ServiceAccountRolesInitializedAnnotation] = "true"
		if _, err := c.KubeClient().Namespaces().Update(defaultNamespace); err != nil {
			glog.Errorf("Error recording adding service account roles to default namespace: %v", err)
		}
	}
}

func (c *MasterConfig) ensureDefaultSecurityContextConstraints() {
	sccList, err := c.KubeClient().SecurityContextConstraints().List(labels.Everything(), fields.Everything())
	if err != nil {
		glog.Errorf("Unable to initialize security context constraints: %v", err)
	}
	if len(sccList.Items) > 0 {
		return
	}

	glog.Infof("No security context constraints detected, adding defaults")
	ns := c.Options.PolicyConfig.OpenShiftInfrastructureNamespace
	buildControllerUsername := serviceaccount.MakeUsername(ns, c.BuildControllerServiceAccount)
	for _, scc := range bootstrappolicy.GetBootstrapSecurityContextConstraints(buildControllerUsername) {
		_, err = c.KubeClient().SecurityContextConstraints().Create(&scc)
		if err != nil {
			glog.Errorf("Unable to create default security context constraint %s.  Got error: %v", scc.Name, err)
		}
	}
}

func (c *MasterConfig) authorizationFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		attributes, err := c.AuthorizationAttributeBuilder.GetAttributes(req)
		if err != nil {
			forbidden(err.Error(), "", w, req)
			return
		}
		if attributes == nil {
			forbidden("No attributes", "", w, req)
			return
		}

		ctx, exists := c.RequestContextMapper.Get(req)
		if !exists {
			forbidden("context not found", attributes.GetAPIVersion(), w, req)
			return
		}

		allowed, reason, err := c.Authorizer.Authorize(ctx, attributes)
		if err != nil {
			forbidden(err.Error(), attributes.GetAPIVersion(), w, req)
			return
		}
		if !allowed {
			forbidden(reason, attributes.GetAPIVersion(), w, req)
			return
		}

		handler.ServeHTTP(w, req)
	})
}

// cacheControlFilter sets the Cache-Control header to the specified value.
func cacheControlFilter(handler http.Handler, value string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", value)
		handler.ServeHTTP(w, req)
	})
}

// forbidden renders a simple forbidden error
func forbidden(reason, apiVersion string, w http.ResponseWriter, req *http.Request) {
	// the api version can be empty for two basic reasons:
	// 1. malformed API request
	// 2. not an API request at all
	// In these cases, just assume the latest version that will work better than nothing
	if len(apiVersion) == 0 {
		apiVersion = klatest.Version
	}

	// Reason is an opaque string that describes why access is allowed or forbidden (forbidden by the time we reach here).
	// We don't have direct access to kind or name (not that those apply either in the general case)
	// We create a NewForbidden to stay close the API, but then we override the message to get a serialization
	// that makes sense when a human reads it.
	forbiddenError, _ := kapierror.NewForbidden("", "", errors.New("")).(*kapierror.StatusError)
	forbiddenError.ErrStatus.Message = reason

	// Not all API versions in valid API requests will have a matching codec in kubernetes.  If we can't find one,
	// just default to the latest kube codec.
	codec := klatest.Codec
	if requestedCodec, err := klatest.InterfacesFor(apiVersion); err == nil {
		codec = requestedCodec
	}

	formatted := &bytes.Buffer{}
	output, err := codec.Encode(&forbiddenError.ErrStatus)
	if err != nil {
		fmt.Fprintf(formatted, "%s", forbiddenError.Error())
	} else {
		_ = json.Indent(formatted, output, "", "  ")
	}

	w.Header().Set("Content-Type", restful.MIME_JSON)
	w.WriteHeader(http.StatusForbidden)
	w.Write(formatted.Bytes())
}

// RunProjectAuthorizationCache starts the project authorization cache
func (c *MasterConfig) RunProjectAuthorizationCache() {
	// TODO: look at exposing a configuration option in future to control how often we run this loop
	period := 1 * time.Second
	c.ProjectAuthorizationCache.Run(period)
}

// RunOriginNamespaceController starts the controller that takes part in namespace termination of openshift content
func (c *MasterConfig) RunOriginNamespaceController() {
	osclient, kclient := c.OriginNamespaceControllerClients()
	factory := projectcontroller.NamespaceControllerFactory{
		Client:     osclient,
		KubeClient: kclient,
	}
	controller := factory.Create()
	controller.Run()
}

// RunServiceAccountsController starts the service account controller
func (c *MasterConfig) RunServiceAccountsController() {
	if len(c.Options.ServiceAccountConfig.ManagedNames) == 0 {
		glog.Infof("Skipped starting Service Account Manager, no managed names specified")
		return
	}
	options := serviceaccount.DefaultServiceAccountsControllerOptions()
	options.Names = util.NewStringSet(c.Options.ServiceAccountConfig.ManagedNames...)
	serviceaccount.NewServiceAccountsController(c.KubeClient(), options).Run()
	glog.Infof("Started Service Account Manager")
}

// RunServiceAccountTokensController starts the service account token controller
func (c *MasterConfig) RunServiceAccountTokensController() {
	if len(c.Options.ServiceAccountConfig.PrivateKeyFile) == 0 {
		glog.Infof("Skipped starting Service Account Token Manager, no private key specified")
		return
	}

	privateKey, err := serviceaccount.ReadPrivateKey(c.Options.ServiceAccountConfig.PrivateKeyFile)
	if err != nil {
		glog.Fatalf("Error reading signing key for Service Account Token Manager: %v", err)
	}
	options := serviceaccount.DefaultTokenControllerOptions(serviceaccount.JWTTokenGenerator(privateKey))

	serviceaccount.NewTokensController(c.KubeClient(), options).Run()
	glog.Infof("Started Service Account Token Manager")
}

// RunServiceAccountPullSecretsControllers starts the service account pull secret controllers
func (c *MasterConfig) RunServiceAccountPullSecretsControllers() {
	serviceaccountcontrollers.NewDockercfgDeletedController(c.KubeClient(), serviceaccountcontrollers.DockercfgDeletedControllerOptions{}).Run()
	serviceaccountcontrollers.NewDockercfgTokenDeletedController(c.KubeClient(), serviceaccountcontrollers.DockercfgTokenDeletedControllerOptions{}).Run()

	dockercfgController := serviceaccountcontrollers.NewDockercfgController(c.KubeClient(), serviceaccountcontrollers.DockercfgControllerOptions{DefaultDockerURL: serviceaccountcontrollers.DefaultOpenshiftDockerURL})
	dockercfgController.Run()

	dockerRegistryControllerOptions := serviceaccountcontrollers.DockerRegistryServiceControllerOptions{
		RegistryNamespace:   "default",
		RegistryServiceName: "docker-registry",
		DockercfgController: dockercfgController,
		DefaultDockerURL:    serviceaccountcontrollers.DefaultOpenshiftDockerURL,
	}
	serviceaccountcontrollers.NewDockerRegistryServiceController(c.KubeClient(), dockerRegistryControllerOptions).Run()

	glog.Infof("Started Service Account Pull Secret Controllers")
}

// RunPolicyCache starts the policy cache
func (c *MasterConfig) RunPolicyCache() {
	c.PolicyCache.Run()
}

// RunAssetServer starts the asset server for the OpenShift UI.
func (c *MasterConfig) RunAssetServer() {

}

// RunDNSServer starts the DNS server
func (c *MasterConfig) RunDNSServer() {
	config, err := dns.NewServerDefaults()
	if err != nil {
		glog.Fatalf("Could not start DNS: %v", err)
	}
	config.DnsAddr = c.Options.DNSConfig.BindAddress
	config.NoRec = true // do not want to deploy an open resolver

	_, port, err := net.SplitHostPort(c.Options.DNSConfig.BindAddress)
	if err != nil {
		glog.Fatalf("Could not start DNS: %v", err)
	}
	if port != "53" {
		glog.Warningf("Binding DNS on port %v instead of 53 (you may need to run as root and update your config), using %s which will not resolve from all locations", port, c.Options.DNSConfig.BindAddress)
	}

	if ok, err := cmdutil.TryListen(c.Options.DNSConfig.BindAddress); !ok {
		glog.Warningf("Could not start DNS: %v", err)
		return
	}

	go func() {
		err := dns.ListenAndServe(config, c.DNSServerClient(), c.EtcdHelper.Client.(*etcdclient.Client))
		glog.Fatalf("Could not start DNS: %v", err)
	}()

	cmdutil.WaitForSuccessfulDial(false, "tcp", c.Options.DNSConfig.BindAddress, 100*time.Millisecond, 100*time.Millisecond, 100)

	glog.Infof("OpenShift DNS listening at %s", c.Options.DNSConfig.BindAddress)
}

// RunProjectCache populates project cache, used by scheduler and project admission controller.
func (c *MasterConfig) RunProjectCache() {
	glog.Infof("Using default project node label selector: %s", c.Options.ProjectConfig.DefaultNodeSelector)
	projectcache.RunProjectCache(c.PrivilegedLoopbackKubernetesClient, c.Options.ProjectConfig.DefaultNodeSelector)
}

// RunBuildController starts the build sync loop for builds and buildConfig processing.
func (c *MasterConfig) RunBuildController() {
	// initialize build controller
	dockerImage := c.ImageFor("docker-builder")
	stiImage := c.ImageFor("sti-builder")

	storageVersion := c.Options.EtcdStorageConfig.OpenShiftStorageVersion
	interfaces, err := latest.InterfacesFor(storageVersion)
	if err != nil {
		glog.Fatalf("Unable to load storage version %s: %v", storageVersion, err)
	}

	osclient, kclient := c.BuildControllerClients()
	factory := buildcontrollerfactory.BuildControllerFactory{
		OSClient:     osclient,
		KubeClient:   kclient,
		BuildUpdater: buildclient.NewOSClientBuildClient(osclient),
		OpenshiftEnabled: c.OpenshiftEnabled,
		DockerBuildStrategy: &buildstrategy.DockerBuildStrategy{
			Image: dockerImage,
			// TODO: this will be set to --storage-version (the internal schema we use)
			Codec: interfaces.Codec,
		},
		SourceBuildStrategy: &buildstrategy.SourceBuildStrategy{
			Image:                stiImage,
			TempDirectoryCreator: buildstrategy.STITempDirectoryCreator,
			// TODO: this will be set to --storage-version (the internal schema we use)
			Codec: interfaces.Codec,
		},
		CustomBuildStrategy: &buildstrategy.CustomBuildStrategy{
			// TODO: this will be set to --storage-version (the internal schema we use)
			Codec: interfaces.Codec,
		},
	}

	controller := factory.Create()
	controller.Run()
	deleteController := factory.CreateDeleteController()
	deleteController.Run()
}

// RunBuildPodController starts the build/pod status sync loop for build status
func (c *MasterConfig) RunBuildPodController() {
	osclient, kclient := c.BuildPodControllerClients()
	factory := buildcontrollerfactory.BuildPodControllerFactory{
		OSClient:     osclient,
		KubeClient:   kclient,
		BuildUpdater: buildclient.NewOSClientBuildClient(osclient),
	}
	controller := factory.Create()
	controller.Run()
	deletecontroller := factory.CreateDeleteController()
	deletecontroller.Run()
}

// RunBuildImageChangeTriggerController starts the build image change trigger controller process.
func (c *MasterConfig) RunBuildImageChangeTriggerController() {
	bcClient, _ := c.BuildImageChangeTriggerControllerClients()
	bcInstantiator := buildclient.NewOSClientBuildConfigInstantiatorClient(bcClient)
	factory := buildcontrollerfactory.ImageChangeControllerFactory{Client: bcClient, BuildConfigInstantiator: bcInstantiator}
	factory.Create().Run()
}

// RunDeploymentController starts the deployment controller process.
func (c *MasterConfig) RunDeploymentController() {
	_, kclient := c.DeploymentControllerClients()

	_, kclientConfig, err := configapi.GetKubeClient(c.Options.MasterClients.OpenShiftLoopbackKubeConfig)
	if err != nil {
		glog.Fatalf("Unable to initialize deployment controller: %v", err)
	}
	// TODO eliminate these environment variables once service accounts provide a kubeconfig that includes all of this info
	env := clientcmd.EnvVars(
		kclientConfig.Host,
		kclientConfig.CAData,
		kclientConfig.Insecure,
		path.Join(serviceaccountadmission.DefaultAPITokenMountPath, kapi.ServiceAccountTokenKey),
	)

	factory := deploycontroller.DeploymentControllerFactory{
		KubeClient:     kclient,
		Codec:          c.EtcdHelper.Codec,
		Environment:    env,
		DeployerImage:  c.ImageFor("deployer"),
		ServiceAccount: bootstrappolicy.DeployerServiceAccountName,
	}

	controller := factory.Create()
	controller.Run()
}

// RunDeployerPodController starts the deployer pod controller process.
func (c *MasterConfig) RunDeployerPodController() {
	_, kclient := c.DeployerPodControllerClients()
	factory := deployerpodcontroller.DeployerPodControllerFactory{
		KubeClient: kclient,
	}

	controller := factory.Create()
	controller.Run()
}

// RunDeploymentConfigController starts the deployment config controller process.
func (c *MasterConfig) RunDeploymentConfigController() {
	osclient, kclient := c.DeploymentConfigControllerClients()
	factory := deployconfigcontroller.DeploymentConfigControllerFactory{
		Client:     osclient,
		KubeClient: kclient,
		Codec:      c.EtcdHelper.Codec,
	}
	controller := factory.Create()
	controller.Run()
}

// RunDeploymentConfigChangeController starts the deployment config change controller process.
func (c *MasterConfig) RunDeploymentConfigChangeController() {
	osclient, kclient := c.DeploymentConfigChangeControllerClients()
	factory := configchangecontroller.DeploymentConfigChangeControllerFactory{
		Client:     osclient,
		KubeClient: kclient,
		Codec:      c.EtcdHelper.Codec,
	}
	controller := factory.Create()
	controller.Run()
}

// RunDeploymentImageChangeTriggerController starts the image change trigger controller process.
func (c *MasterConfig) RunDeploymentImageChangeTriggerController() {
	osclient := c.DeploymentImageChangeTriggerControllerClient()
	factory := imagechangecontroller.ImageChangeControllerFactory{Client: osclient}
	controller := factory.Create()
	controller.Run()
}

// RunSDNController runs openshift-sdn if the said network plugin is provided
func (c *MasterConfig) RunSDNController() {
	osclient, kclient := c.SDNControllerClients()
	if c.Options.NetworkConfig.NetworkPluginName == osdn.NetworkPluginName() {
		osdn.Master(osclient, kclient, c.Options.NetworkConfig.ClusterNetworkCIDR, c.Options.NetworkConfig.HostSubnetLength)
	}
}

// RouteAllocator returns a route allocation controller.
func (c *MasterConfig) RouteAllocator() *routeallocationcontroller.RouteAllocationController {
	osclient, kclient := c.RouteAllocatorClients()
	factory := routeallocationcontroller.RouteAllocationControllerFactory{
		OSClient:   osclient,
		KubeClient: kclient,
	}

	plugin, err := routeplugin.NewSimpleAllocationPlugin(c.Options.RoutingConfig.Subdomain)
	if err != nil {
		glog.Fatalf("Route plugin initialization failed: %v", err)
	}

	return factory.Create(plugin)
}

// RunImageImportController starts the image import trigger controller process.
func (c *MasterConfig) RunImageImportController() {
	osclient := c.ImageImportControllerClient()
	factory := imagecontroller.ImportControllerFactory{
		Client: osclient,
	}
	controller := factory.Create()
	controller.Run()
}

// RunSecurityAllocationController starts the security allocation controller process.
func (c *MasterConfig) RunSecurityAllocationController() {
	alloc := c.Options.ProjectConfig.SecurityAllocator
	if alloc == nil {
		glog.V(3).Infof("Security allocator is disabled - no UIDs assigned to projects")
		return
	}

	// TODO: move range initialization to run_config
	uidRange, err := uid.ParseRange(alloc.UIDAllocatorRange)

	if err != nil {
		glog.Fatalf("Unable to describe UID range: %v", err)
	}
	var etcdAlloc *etcdallocator.Etcd
	uidAllocator := uidallocator.New(uidRange, func(max int, rangeSpec string) allocator.Interface {
		mem := allocator.NewContiguousAllocationMap(max, rangeSpec)
		etcdAlloc = etcdallocator.NewEtcd(mem, "/ranges/uids", "uidallocation", c.EtcdHelper)
		return etcdAlloc
	})
	mcsRange, err := mcs.ParseRange(alloc.MCSAllocatorRange)
	if err != nil {
		glog.Fatalf("Unable to describe MCS category range: %v", err)
	}

	kclient := c.SecurityAllocationControllerClient()

	repair := securitycontroller.NewRepair(time.Minute, kclient.Namespaces(), uidRange, etcdAlloc)
	if err := repair.RunOnce(); err != nil {
		// TODO: v scary, may need to use direct etcd calls?
		glog.Fatalf("Unable to initialize namespaces: %v", err)
	}

	factory := securitycontroller.AllocationFactory{
		UIDAllocator: uidAllocator,
		MCSAllocator: securitycontroller.DefaultMCSAllocation(uidRange, mcsRange, alloc.MCSLabelsPerProject),
		Client:       kclient.Namespaces(),
		// TODO: reuse namespace cache
	}
	controller := factory.Create()
	controller.Run()
}

// ensureCORSAllowedOrigins takes a string list of origins and attempts to covert them to CORS origin
// regexes, or exits if it cannot.
func (c *MasterConfig) ensureCORSAllowedOrigins() []*regexp.Regexp {
	if len(c.Options.CORSAllowedOrigins) == 0 {
		return []*regexp.Regexp{}
	}
	allowedOriginRegexps, err := util.CompileRegexps(util.StringList(c.Options.CORSAllowedOrigins))
	if err != nil {
		glog.Fatalf("Invalid --cors-allowed-origins: %v", err)
	}
	return allowedOriginRegexps
}

// env returns an environment variable, or the defaultValue if it is not set.
func env(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

type clientDeploymentInterface struct {
	KubeClient kclient.Interface
}

// GetDeployment returns the deployment with the provided context and name
func (c clientDeploymentInterface) GetDeployment(ctx api.Context, name string) (*api.ReplicationController, error) {
	return c.KubeClient.ReplicationControllers(api.NamespaceValue(ctx)).Get(name)
}

// namespacingFilter adds a filter that adds the namespace of the request to the context.  Not all requests will have namespaces,
// but any that do will have the appropriate value added.
func namespacingFilter(handler http.Handler, contextMapper kapi.RequestContextMapper) http.Handler {
	infoResolver := &apiserver.APIRequestInfoResolver{util.NewStringSet("api", "osapi", "oapi"), latest.RESTMapper}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx, ok := contextMapper.Get(req)
		if !ok {
			http.Error(w, "Unable to find request context", http.StatusInternalServerError)
			return
		}

		if _, exists := kapi.NamespaceFrom(ctx); !exists {
			if requestInfo, err := infoResolver.GetAPIRequestInfo(req); err == nil {
				// only set the namespace if the apiRequestInfo was resolved
				// keep in mind that GetAPIRequestInfo will fail on non-api requests, so don't fail the entire http request on that
				// kind of failure.

				// TODO reconsider special casing this.  Having the special case hereallow us to fully share the kube
				// APIRequestInfoResolver without any modification or customization.
				namespace := requestInfo.Namespace
				if (requestInfo.Resource == "projects") && (len(requestInfo.Name) > 0) {
					namespace = requestInfo.Name
				}

				ctx = kapi.WithNamespace(ctx, namespace)
				contextMapper.Update(req, ctx)
			}
		}

		handler.ServeHTTP(w, req)
	})
}
