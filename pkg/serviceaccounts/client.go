package serviceaccounts

import (
	"fmt"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"

	"github.com/projectatomic/appinfra-next/pkg/client"
)

// TokenRetriever defined an interface for getting an API token for a service account
type TokenRetriever interface {
	GetToken(serviceAccountNamespace, serviceAccountName string) (token string, err error)
}

// ClientLookupTokenRetriever uses its client to look up a service account token
type ClientLookupTokenRetriever struct {
	Client kclient.Interface
}

// GetToken returns a token for the named service account or an error if none existed after a timeout
func (s *ClientLookupTokenRetriever) GetToken(namespace, name string) (string, error) {
	for i := 0; i < 30; i++ {
		// Wait on subsequent retries
		if i != 0 {
			time.Sleep(time.Second)
		}

		// Get the service account
		serviceAccount, err := s.Client.ServiceAccounts(namespace).Get(name)
		if err != nil {
			continue
		}

		// Get the secrets
		// TODO: JTL: create one directly once we have that ability
		for _, secretRef := range serviceAccount.Secrets {
			secret, err := s.Client.Secrets(namespace).Get(secretRef.Name)
			if err != nil {
				// Tolerate fetch errors on a particular secret
				continue
			}
			if IsValidServiceAccountToken(serviceAccount, secret) {
				// Return a valid token
				return string(secret.Data[kapi.ServiceAccountTokenKey]), nil
			}
		}
	}

	return "", fmt.Errorf("Could not get token for %s/%s", namespace, name)
}

// Clients returns an OpenShift and Kubernetes client with the credentials of the named service account
// TODO: change return types to client.Interface/kclient.Interface to allow auto-reloading credentials
func Clients(config kclient.Config, tokenRetriever TokenRetriever, namespace, name string) (*client.Client, *kclient.Client, error) {
	// Clear existing auth info
	config.Username = ""
	config.Password = ""
	config.CertFile = ""
	config.CertData = []byte{}
	config.KeyFile = ""
	config.KeyData = []byte{}

	// For now, just initialize the token once
	// TODO: refetch the token if the client encounters 401 errors
	token, err := tokenRetriever.GetToken(namespace, name)
	if err != nil {
		return nil, nil, err
	}
	config.BearerToken = token

	c, err := client.New(&config)
	if err != nil {
		return nil, nil, err
	}

	kc, err := kclient.New(&config)
	if err != nil {
		return nil, nil, err
	}

	return c, kc, nil
}

// IsValidServiceAccountToken returns true if the given secret contains a service account token valid for the given service account
func IsValidServiceAccountToken(serviceAccount *kapi.ServiceAccount, secret *kapi.Secret) bool {
	if secret.Type != kapi.SecretTypeServiceAccountToken {
		return false
	}
	if secret.Namespace != serviceAccount.Namespace {
		return false
	}
	if secret.Annotations[kapi.ServiceAccountNameKey] != serviceAccount.Name {
		return false
	}
	if secret.Annotations[kapi.ServiceAccountUIDKey] != string(serviceAccount.UID) {
		return false
	}
	if len(secret.Data[kapi.ServiceAccountTokenKey]) == 0 {
		return false
	}
	return true
}
