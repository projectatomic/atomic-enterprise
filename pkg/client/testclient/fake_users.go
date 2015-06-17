package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	userapi "github.com/projectatomic/appinfra-next/pkg/user/api"
)

// FakeUsers implements UsersInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeUsers struct {
	Fake *Fake
}

func (c *FakeUsers) List(label labels.Selector, field fields.Selector) (*userapi.UserList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-users"}, &userapi.UserList{})
	return obj.(*userapi.UserList), err
}

func (c *FakeUsers) Get(name string) (*userapi.User, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-user", Value: name}, &userapi.User{})
	return obj.(*userapi.User), err
}

func (c *FakeUsers) Create(user *userapi.User) (*userapi.User, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-user", Value: user}, &userapi.User{})
	return obj.(*userapi.User), err
}

func (c *FakeUsers) Update(user *userapi.User) (*userapi.User, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "update-user", Value: user}, &userapi.User{})
	return obj.(*userapi.User), err
}
