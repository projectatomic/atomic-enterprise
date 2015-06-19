package admission

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/admission"
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	apierrors "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	ktestclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client/testclient"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/client/testclient"
)

func TestBuildAdmission(t *testing.T) {
	tests := []struct {
		name             string
		kind             string
		resource         string
		object           runtime.Object
		reviewResponse   *authorizationapi.SubjectAccessReviewResponse
		expectedResource string
		expectAccept     bool
	}{
		{
			name:             "allowed source build",
			object:           testBuild(buildapi.SourceBuildStrategyType),
			kind:             "Build",
			resource:         buildsResource,
			reviewResponse:   reviewResponse(true, ""),
			expectedResource: authorizationapi.SourceBuildResource,
			expectAccept:     true,
		},
		{
			name:             "denied docker build",
			object:           testBuild(buildapi.DockerBuildStrategyType),
			kind:             "Build",
			resource:         buildsResource,
			reviewResponse:   reviewResponse(false, "cannot create build of type docker build"),
			expectAccept:     false,
			expectedResource: authorizationapi.DockerBuildResource,
		},
		{
			name:             "allowed custom build",
			object:           testBuild(buildapi.CustomBuildStrategyType),
			kind:             "Build",
			resource:         buildsResource,
			reviewResponse:   reviewResponse(true, ""),
			expectedResource: authorizationapi.CustomBuildResource,
			expectAccept:     true,
		},
		{
			name:             "allowed build config",
			object:           testBuildConfig(buildapi.DockerBuildStrategyType),
			kind:             "BuildConfig",
			resource:         buildConfigsResource,
			reviewResponse:   reviewResponse(true, ""),
			expectAccept:     true,
			expectedResource: authorizationapi.DockerBuildResource,
		},
		{
			name:             "forbidden build config",
			object:           testBuildConfig(buildapi.CustomBuildStrategyType),
			kind:             "BuildConfig",
			resource:         buildConfigsResource,
			reviewResponse:   reviewResponse(false, ""),
			expectAccept:     false,
			expectedResource: authorizationapi.CustomBuildResource,
		},
	}

	for _, test := range tests {
		c := NewBuildByStrategy(fakeClient(test.expectedResource, test.reviewResponse))
		attrs := admission.NewAttributesRecord(test.object, test.kind, "default", test.resource, admission.Create, fakeUser())
		err := c.Admit(attrs)
		if err != nil && test.expectAccept {
			t.Errorf("%s: unexpected error: %v", test.name, err)
		}
		if (err == nil || !apierrors.IsForbidden(err)) && !test.expectAccept {
			t.Errorf("%s: expecting reject error", test.name)
		}
	}
}

func fakeUser() user.Info {
	return &user.DefaultInfo{
		Name: "testuser",
	}
}

func fakeClient(expectedResource string, reviewResponse *authorizationapi.SubjectAccessReviewResponse) client.Interface {
	emptyResponse := &authorizationapi.SubjectAccessReviewResponse{}
	return &testclient.Fake{
		ReactFn: func(action ktestclient.FakeAction) (runtime.Object, error) {
			if action.Action == "create-subjectAccessReview" {
				review, ok := action.Value.(*authorizationapi.SubjectAccessReview)
				if !ok {
					return emptyResponse, fmt.Errorf("unexpected object received: %#v", review)
				}
				if review.Resource != expectedResource {
					return emptyResponse, fmt.Errorf("unexpected resource received: %s. expected: %s",
						review.Resource, expectedResource)
				}
				return reviewResponse, nil
			}
			return nil, nil
		},
	}
}

func testBuild(strategy buildapi.BuildStrategyType) *buildapi.Build {
	return &buildapi.Build{
		ObjectMeta: kapi.ObjectMeta{
			Name: "test-build",
		},
		Parameters: buildapi.BuildParameters{
			Strategy: buildapi.BuildStrategy{
				Type: strategy,
			},
		},
	}
}

func testBuildConfig(strategy buildapi.BuildStrategyType) *buildapi.BuildConfig {
	return &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{
			Name: "test-buildconfig",
		},
		Parameters: buildapi.BuildParameters{
			Strategy: buildapi.BuildStrategy{
				Type: strategy,
			},
		},
	}
}

func reviewResponse(allowed bool, msg string) *authorizationapi.SubjectAccessReviewResponse {
	return &authorizationapi.SubjectAccessReviewResponse{
		Allowed: allowed,
		Reason:  msg,
	}
}
