package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
)

// FakeRoleBindings implements RoleBindingInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeRoleBindings struct {
	Fake *Fake
}

func (c *FakeRoleBindings) List(label labels.Selector, field fields.Selector) (*authorizationapi.RoleBindingList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-roleBinding"}, &authorizationapi.RoleBindingList{})
	return obj.(*authorizationapi.RoleBindingList), err
}

func (c *FakeRoleBindings) Get(name string) (*authorizationapi.RoleBinding, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-roleBinding"}, &authorizationapi.RoleBinding{})
	return obj.(*authorizationapi.RoleBinding), err
}

func (c *FakeRoleBindings) Create(roleBinding *authorizationapi.RoleBinding) (*authorizationapi.RoleBinding, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-roleBinding", Value: roleBinding}, &authorizationapi.RoleBinding{})
	return obj.(*authorizationapi.RoleBinding), err
}

func (c *FakeRoleBindings) Update(roleBinding *authorizationapi.RoleBinding) (*authorizationapi.RoleBinding, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "update-roleBinding"}, &authorizationapi.RoleBinding{})
	return obj.(*authorizationapi.RoleBinding), err
}

func (c *FakeRoleBindings) Delete(name string) error {
	c.Fake.Actions = append(c.Fake.Actions, FakeAction{Action: "delete-roleBinding", Value: name})
	return nil
}
