package controller

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	osclient "github.com/projectatomic/appinfra-next/pkg/client"
	projectutil "github.com/projectatomic/appinfra-next/pkg/project/util"
)

// NamespaceController is responsible for participating in Kubernetes Namespace termination
// Use the NamespaceControllerFactory to create this controller.
type NamespaceController struct {
	// Client is an OpenShift client.
	Client osclient.Interface
	// KubeClient is a Kubernetes client.
	KubeClient kclient.Interface
}

// fatalError is an error which can't be retried.
type fatalError string

// Error implements the interface for errors
func (e fatalError) Error() string { return "fatal error handling namespace: " + string(e) }

// Handle processes a namespace and deletes content in origin if its terminating
func (c *NamespaceController) Handle(namespace *kapi.Namespace) (err error) {
	// if namespace is not terminating, ignore it
	if namespace.Status.Phase != kapi.NamespaceTerminating {
		return nil
	}

	// if we already processed this namespace, ignore it
	if projectutil.Finalized(namespace) {
		return nil
	}

	// there may still be content for us to remove
	err = deleteAllContent(c.Client, namespace.Name)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	_, err = projectutil.Finalize(c.KubeClient, namespace)
	if err != nil {
		return err
	}

	return nil
}

// deleteAllContent will purge all content in openshift in the specified namespace
func deleteAllContent(client osclient.Interface, namespace string) (err error) {
	err = deleteBuildConfigs(client, namespace)
	if err != nil {
		return err
	}
	err = deleteBuilds(client, namespace)
	if err != nil {
		return err
	}
	err = deleteDeploymentConfigs(client, namespace)
	if err != nil {
		return err
	}
	err = deleteImageStreams(client, namespace)
	if err != nil {
		return err
	}
	err = deletePolicies(client, namespace)
	if err != nil {
		return err
	}
	err = deletePolicyBindings(client, namespace)
	if err != nil {
		return err
	}
	err = deleteRoleBindings(client, namespace)
	if err != nil {
		return err
	}
	err = deleteRoles(client, namespace)
	if err != nil {
		return err
	}
	err = deleteRoutes(client, namespace)
	if err != nil {
		return err
	}
	err = deleteTemplates(client, namespace)
	if err != nil {
		return err
	}
	return nil
}

func deleteTemplates(client osclient.Interface, ns string) error {
	items, err := client.Templates(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.Templates(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteRoutes(client osclient.Interface, ns string) error {
	items, err := client.Routes(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.Routes(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteRoles(client osclient.Interface, ns string) error {
	items, err := client.Roles(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.Roles(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteRoleBindings(client osclient.Interface, ns string) error {
	items, err := client.RoleBindings(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.RoleBindings(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deletePolicyBindings(client osclient.Interface, ns string) error {
	items, err := client.PolicyBindings(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.PolicyBindings(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deletePolicies(client osclient.Interface, ns string) error {
	items, err := client.Policies(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.Policies(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteImageStreams(client osclient.Interface, ns string) error {
	items, err := client.ImageStreams(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.ImageStreams(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteDeploymentConfigs(client osclient.Interface, ns string) error {
	items, err := client.DeploymentConfigs(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.DeploymentConfigs(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteBuilds(client osclient.Interface, ns string) error {
	items, err := client.Builds(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.Builds(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func deleteBuildConfigs(client osclient.Interface, ns string) error {
	items, err := client.BuildConfigs(ns).List(labels.Everything(), fields.Everything())
	if err != nil {
		return err
	}
	for i := range items.Items {
		err := client.BuildConfigs(ns).Delete(items.Items[i].Name)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}
