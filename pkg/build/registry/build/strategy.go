package build

import (
	"fmt"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/fielderrors"

	"github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/build/api/validation"
	buildutil "github.com/projectatomic/appinfra-next/pkg/build/util"
)

// strategy implements behavior for Build objects
type strategy struct {
	runtime.ObjectTyper
	kapi.NameGenerator
}

// Decorator is used to compute duration of a build since its not stored in etcd yet
var Decorator = func(obj runtime.Object) error {
	build, ok := obj.(*api.Build)
	if !ok {
		return errors.NewBadRequest(fmt.Sprintf("not a build: %v", build))
	}
	if build.StartTimestamp == nil {
		build.Duration = time.Duration(0)
	} else {
		completionTimestamp := build.CompletionTimestamp
		if completionTimestamp == nil {
			dummy := util.Now()
			completionTimestamp = &dummy
		}
		build.Duration = completionTimestamp.Rfc3339Copy().Time.Sub(build.StartTimestamp.Rfc3339Copy().Time)
	}
	return nil
}

// Strategy is the default logic that applies when creating and updating Build objects.
var Strategy = strategy{kapi.Scheme, kapi.SimpleNameGenerator}

func (strategy) NamespaceScoped() bool {
	return true
}

// AllowCreateOnUpdate is false for Build objects.
func (strategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (strategy) PrepareForCreate(obj runtime.Object) {
	build := obj.(*api.Build)
	if len(build.Status) == 0 {
		build.Status = api.BuildStatusNew
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (strategy) PrepareForUpdate(obj, old runtime.Object) {
	_ = obj.(*api.Build)
}

// Validate validates a new policy.
func (strategy) Validate(ctx kapi.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateBuild(obj.(*api.Build))
}

// ValidateUpdate is the default update validation for an end user.
func (strategy) ValidateUpdate(ctx kapi.Context, obj, old runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateBuildUpdate(obj.(*api.Build), old.(*api.Build))
}

// CheckGracefulDelete allows a build to be gracefully deleted.
func (strategy) CheckGracefulDelete(obj runtime.Object, options *kapi.DeleteOptions) bool {
	return false
}

// Matcher returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return &generic.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			build, ok := obj.(*api.Build)
			if !ok {
				return nil, nil, fmt.Errorf("not a build")
			}
			return labels.Set(build.ObjectMeta.Labels), SelectableFields(build), nil
		},
	}
}

// SelectableFields returns a label set that represents the object
func SelectableFields(build *api.Build) fields.Set {
	return fields.Set{
		"metadata.name": build.Name,
		"status":        string(build.Status),
		"podName":       buildutil.GetBuildPodName(build),
	}
}
