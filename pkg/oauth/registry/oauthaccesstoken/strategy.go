package oauthaccesstoken

import (
	"fmt"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/fielderrors"
	"github.com/projectatomic/appinfra-next/pkg/oauth/api"
	"github.com/projectatomic/appinfra-next/pkg/oauth/api/validation"
)

// strategy implements behavior for OAuthAccessTokens
type strategy struct {
	runtime.ObjectTyper
}

// Strategy is the default logic that applies when creating OAuthAccessToken
// objects via the REST API.
var Strategy = strategy{kapi.Scheme}

func (strategy) PrepareForUpdate(obj, old runtime.Object) {}

// NamespaceScoped is false for OAuth objects
func (strategy) NamespaceScoped() bool {
	return false
}

func (strategy) GenerateName(base string) string {
	return base
}

func (strategy) PrepareForCreate(obj runtime.Object) {
}

// Validate validates a new token
func (strategy) Validate(ctx kapi.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	token := obj.(*api.OAuthAccessToken)
	return validation.ValidateAccessToken(token)
}

// AllowCreateOnUpdate is false for OAuth objects
func (strategy) AllowCreateOnUpdate() bool {
	return false
}

// Matchtoken returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return generic.MatcherFunc(func(obj runtime.Object) (bool, error) {
		tokenObj, ok := obj.(*api.OAuthAccessToken)
		if !ok {
			return false, fmt.Errorf("not a token")
		}
		fields := SelectableFields(tokenObj)
		return label.Matches(labels.Set(tokenObj.Labels)) && field.Matches(fields), nil
	})
}

// SelectableFields returns a label set that represents the object
func SelectableFields(obj *api.OAuthAccessToken) labels.Set {
	return labels.Set{
		"name":           obj.Name,
		"clientName":     obj.ClientName,
		"userName":       obj.UserName,
		"userUID":        obj.UserUID,
		"authorizeToken": obj.AuthorizeToken,
	}
}
