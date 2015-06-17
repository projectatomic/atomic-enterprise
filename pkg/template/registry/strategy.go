package registry

import (
	"fmt"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/fielderrors"

	"github.com/projectatomic/appinfra-next/pkg/template/api"
	"github.com/projectatomic/appinfra-next/pkg/template/api/validation"
)

// templateStrategy implements behavior for Templates
type templateStrategy struct {
	runtime.ObjectTyper
	kapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Template
// objects via the REST API.
var Strategy = templateStrategy{kapi.Scheme, kapi.SimpleNameGenerator}

// NamespaceScoped is true for templates.
func (templateStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (templateStrategy) PrepareForUpdate(obj, old runtime.Object) {}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (templateStrategy) PrepareForCreate(obj runtime.Object) {
}

// Validate validates a new template.
func (templateStrategy) Validate(ctx kapi.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateTemplate(obj.(*api.Template))
}

// AllowCreateOnUpdate is false for templates.
func (templateStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (templateStrategy) ValidateUpdate(ctx kapi.Context, obj, old runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateTemplateUpdate(obj.(*api.Template), old.(*api.Template))
}

// MatchTemplate returns a generic matcher for a given label and field selector.
func MatchTemplate(label labels.Selector, field fields.Selector) generic.Matcher {
	return generic.MatcherFunc(func(obj runtime.Object) (bool, error) {
		o, ok := obj.(*api.Template)
		if !ok {
			return false, fmt.Errorf("not a pod")
		}
		return label.Matches(labels.Set(o.Labels)), nil
	})
}
