// +build integration,!no-etcd

package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd"
	clientcmdapi "github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/serviceaccount"
	serviceaccountadmission "github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/admission/serviceaccount"

	"github.com/projectatomic/appinfra-next/pkg/cmd/admin/policy"
	"github.com/projectatomic/appinfra-next/pkg/cmd/server/bootstrappolicy"
	testutil "github.com/projectatomic/appinfra-next/test/util"
)

func TestServiceAccountAuthorization(t *testing.T) {
	saNamespace := api.NamespaceDefault
	saName := serviceaccountadmission.DefaultServiceAccountName
	saUsername := serviceaccount.MakeUsername(saNamespace, saName)

	// Start one OpenShift master as "cluster1" to play the external kube server
	cluster1MasterConfig, cluster1AdminConfigFile, err := testutil.StartTestMaster()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cluster1AdminConfig, err := testutil.GetClusterAdminClientConfig(cluster1AdminConfigFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cluster1AdminKubeClient, err := testutil.GetClusterAdminKubeClient(cluster1AdminConfigFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cluster1AdminOSClient, err := testutil.GetClusterAdminClient(cluster1AdminConfigFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get a service account token and build a client
	saToken, err := waitForServiceAccountToken(cluster1AdminKubeClient, saNamespace, saName, 20, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(saToken) == 0 {
		t.Fatalf("token was not created")
	}
	cluster1SAClientConfig := kclient.Config{
		Host:        cluster1AdminConfig.Host,
		Prefix:      cluster1AdminConfig.Prefix,
		BearerToken: saToken,
		TLSClientConfig: kclient.TLSClientConfig{
			CAFile: cluster1AdminConfig.CAFile,
			CAData: cluster1AdminConfig.CAData,
		},
	}
	cluster1SAKubeClient, err := kclient.New(&cluster1SAClientConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Make sure the service account doesn't have access
	failNS := &api.Namespace{ObjectMeta: api.ObjectMeta{Name: "test-fail"}}
	if _, err := cluster1SAKubeClient.Namespaces().Create(failNS); !errors.IsForbidden(err) {
		t.Fatalf("expected forbidden error, got %v", err)
	}

	// Make the service account a cluster admin on cluster1
	addRoleOptions := &policy.RoleModificationOptions{
		RoleName:            bootstrappolicy.ClusterAdminRoleName,
		RoleBindingAccessor: policy.NewClusterRoleBindingAccessor(cluster1AdminOSClient),
		Users:               []string{saUsername},
	}
	if err := addRoleOptions.AddRole(); err != nil {
		t.Fatalf("could not add role to service account")
	}

	// Give the policy cache a second to catch it's breath
	time.Sleep(time.Second)

	// Make sure the service account now has access
	// This tests authentication using the etcd-based token getter
	passNS := &api.Namespace{ObjectMeta: api.ObjectMeta{Name: "test-pass"}}
	if _, err := cluster1SAKubeClient.Namespaces().Create(passNS); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a kubeconfig from the serviceaccount config
	cluster1SAKubeConfigFile, err := ioutil.TempFile(testutil.GetBaseDir(), "cluster1-service-account.kubeconfig")
	if err != nil {
		t.Fatalf("error creating tmpfile: %v", err)
	}
	defer os.Remove(cluster1SAKubeConfigFile.Name())
	if err := writeClientConfigToKubeConfig(cluster1SAClientConfig, cluster1SAKubeConfigFile.Name()); err != nil {
		t.Fatalf("error creating kubeconfig: %v", err)
	}

	// Set up cluster 2 to run against cluster 1 as external kubernetes
	cluster2MasterConfig, err := testutil.DefaultMasterOptions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Don't start kubernetes in process
	cluster2MasterConfig.KubernetesMasterConfig = nil
	// Connect to cluster1 using the service account credentials
	cluster2MasterConfig.MasterClients.ExternalKubernetesKubeConfig = cluster1SAKubeConfigFile.Name()
	// Don't start etcd
	cluster2MasterConfig.EtcdConfig = nil
	// Use the same credentials as cluster1 to connect to existing etcd
	cluster2MasterConfig.EtcdClientInfo = cluster1MasterConfig.EtcdClientInfo
	// Set a custom etcd prefix to make sure data is getting sent to cluster1
	cluster2MasterConfig.EtcdStorageConfig.KubernetesStoragePrefix += "2"
	cluster2MasterConfig.EtcdStorageConfig.OpenShiftStoragePrefix += "2"
	// Don't manage any names in cluster2
	cluster2MasterConfig.ServiceAccountConfig.ManagedNames = []string{}
	// Don't create any service account tokens in cluster2
	cluster2MasterConfig.ServiceAccountConfig.PrivateKeyFile = ""
	// Use the same public keys to validate tokens as cluster1
	cluster2MasterConfig.ServiceAccountConfig.PublicKeyFiles = cluster1MasterConfig.ServiceAccountConfig.PublicKeyFiles

	// Start cluster 2 (without clearing etcd) and get admin client configs and clients
	cluster2Options := testutil.TestOptions{DeleteAllEtcdKeys: false}
	cluster2AdminConfigFile, err := testutil.StartConfiguredMasterWithOptions(cluster2MasterConfig, cluster2Options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cluster2AdminConfig, err := testutil.GetClusterAdminClientConfig(cluster2AdminConfigFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cluster2AdminOSClient, err := testutil.GetClusterAdminClient(cluster2AdminConfigFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Build a client to use the same service account token against cluster2
	cluster2SAClientConfig := cluster1SAClientConfig
	cluster2SAClientConfig.Host = cluster2AdminConfig.Host
	cluster2SAKubeClient, err := kclient.New(&cluster2SAClientConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Make sure the service account doesn't have access
	// A forbidden error makes sure the token was recognized, and policy denied us
	// This exercises the client-based token getter
	// It also makes sure we don't loop back through the cluster2 kube proxy which would cause an auth loop
	failNS2 := &api.Namespace{ObjectMeta: api.ObjectMeta{Name: "test-fail2"}}
	if _, err := cluster2SAKubeClient.Namespaces().Create(failNS2); !errors.IsForbidden(err) {
		t.Fatalf("expected forbidden error, got %v", err)
	}

	// Make the service account a cluster admin on cluster2
	addRoleOptions2 := &policy.RoleModificationOptions{
		RoleName:            bootstrappolicy.ClusterAdminRoleName,
		RoleBindingAccessor: policy.NewClusterRoleBindingAccessor(cluster2AdminOSClient),
		Users:               []string{saUsername},
	}
	if err := addRoleOptions2.AddRole(); err != nil {
		t.Fatalf("could not add role to service account")
	}

	// Give the policy cache a second to catch it's breath
	time.Sleep(time.Second)

	// Make sure the service account now has access to cluster2
	passNS2 := &api.Namespace{ObjectMeta: api.ObjectMeta{Name: "test-pass2"}}
	if _, err := cluster2SAKubeClient.Namespaces().Create(passNS2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Make sure the ns actually got created in cluster1
	if _, err := cluster1SAKubeClient.Namespaces().Get(passNS2.Name); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeClientConfigToKubeConfig(config kclient.Config, path string) error {
	kubeConfig := &clientcmdapi.Config{
		Clusters:       map[string]clientcmdapi.Cluster{"myserver": {Server: config.Host, CertificateAuthority: config.CAFile, CertificateAuthorityData: config.CAData}},
		AuthInfos:      map[string]clientcmdapi.AuthInfo{"myuser": {Token: config.BearerToken}},
		Contexts:       map[string]clientcmdapi.Context{"mycontext": {Cluster: "myserver", AuthInfo: "myuser"}},
		CurrentContext: "mycontext",
	}
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return err
	}
	if err := clientcmd.WriteToFile(*kubeConfig, path); err != nil {
		return err
	}
	return nil
}

func waitForServiceAccountToken(client *kclient.Client, ns, name string, attempts int, interval time.Duration) (string, error) {
	for i := 0; i <= attempts; i++ {
		time.Sleep(interval)
		token, err := getServiceAccountToken(client, ns, name)
		if err != nil {
			return "", err
		}
		if len(token) > 0 {
			return token, nil
		}
	}
	return "", nil
}

func getServiceAccountToken(client *kclient.Client, ns, name string) (string, error) {
	secrets, err := client.Secrets(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return "", err
	}
	for _, secret := range secrets.Items {
		if secret.Type == api.SecretTypeServiceAccountToken && secret.Annotations[api.ServiceAccountNameKey] == name {
			return string(secret.Data[api.ServiceAccountTokenKey]), nil
		}
	}
	return "", nil
}

func TestAutomaticCreationOfPullSecrets(t *testing.T) {
	saNamespace := api.NamespaceDefault
	saName := serviceaccountadmission.DefaultServiceAccountName

	_, clusterAdminConfig, err := testutil.StartTestMaster()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	clusterAdminKubeClient, err := testutil.GetClusterAdminKubeClient(clusterAdminConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get a service account token
	saToken, err := waitForServiceAccountToken(clusterAdminKubeClient, saNamespace, saName, 20, time.Second)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(saToken) == 0 {
		t.Errorf("token was not created")
	}

	// Get the matching dockercfg secret
	saPullSecret, err := waitForServiceAccountPullSecret(clusterAdminKubeClient, saNamespace, saName, 20, time.Second)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(saPullSecret) == 0 {
		t.Errorf("pull secret was not created")
	}
}

func waitForServiceAccountPullSecret(client *kclient.Client, ns, name string, attempts int, interval time.Duration) (string, error) {
	for i := 0; i <= attempts; i++ {
		time.Sleep(interval)
		token, err := getServiceAccountPullSecret(client, ns, name)
		if err != nil {
			return "", err
		}
		if len(token) > 0 {
			return token, nil
		}
	}
	return "", nil
}

func getServiceAccountPullSecret(client *kclient.Client, ns, name string) (string, error) {
	secrets, err := client.Secrets(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return "", err
	}
	for _, secret := range secrets.Items {
		if secret.Type == api.SecretTypeDockercfg && secret.Annotations[api.ServiceAccountNameKey] == name {
			return string(secret.Data[api.DockerConfigKey]), nil
		}
	}
	return "", nil
}
