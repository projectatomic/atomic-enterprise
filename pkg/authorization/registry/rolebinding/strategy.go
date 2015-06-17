package rolebinding

import (
	"fmt"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/fielderrors"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	"github.com/projectatomic/appinfra-next/pkg/authorization/api/validation"
)

// strategy implements behavior for nodes
type strategy struct {
	namespaced bool

	runtime.ObjectTyper
}

var ClusterStrategy = strategy{false, kapi.Scheme}
var LocalStrategy = strategy{true, kapi.Scheme}

// NamespaceScoped is false for rolebindings.
func (s strategy) NamespaceScoped() bool {
	return s.namespaced
}

// AllowCreateOnUpdate is false for rolebindings.
func (s strategy) AllowCreateOnUpdate() bool {
	return false
}

func (s strategy) GenerateName(base string) string {
	return base
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (s strategy) PrepareForCreate(obj runtime.Object) {
	_ = obj.(*authorizationapi.RoleBinding)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (s strategy) PrepareForUpdate(obj, old runtime.Object) {
	_ = obj.(*authorizationapi.RoleBinding)
}

// Validate validates a new role.
func (s strategy) Validate(ctx kapi.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateRoleBinding(obj.(*authorizationapi.RoleBinding), s.namespaced)
}

// ValidateUpdate is the default update validation for an end user.
func (s strategy) ValidateUpdate(ctx kapi.Context, obj, old runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateRoleBindingUpdate(obj.(*authorizationapi.RoleBinding), old.(*authorizationapi.RoleBinding), s.namespaced)
}

// Matcher returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return &generic.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			roleBinding, ok := obj.(*authorizationapi.RoleBinding)
			if !ok {
				return nil, nil, fmt.Errorf("not a rolebinding")
			}
			return labels.Set(roleBinding.ObjectMeta.Labels), SelectableFields(roleBinding), nil
		},
	}
}

// SelectableFields returns a label set that represents the object
func SelectableFields(roleBinding *authorizationapi.RoleBinding) fields.Set {
	return fields.Set{
		"name": roleBinding.Name,
	}
}
