package generator

import (
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/testclient"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/build/generator"
	mocks "github.com/projectatomic/appinfra-next/pkg/build/generator/test"
	imageapi "github.com/projectatomic/appinfra-next/pkg/image/api"
)

func TestCreateInstantiate(t *testing.T) {
	imageStream := mocks.MockImageStream("testImageStream", "registry.com/namespace/imagename", map[string]string{"test": "newImageID123"})
	image := mocks.MockImage("testImage@id", "registry.com/namespace/imagename@id")
	fakeSecrets := []runtime.Object{}
	for _, s := range mocks.MockBuilderSecrets() {
		fakeSecrets = append(fakeSecrets, s)
	}
	rest := InstantiateREST{&generator.BuildGenerator{
		Secrets:         testclient.NewSimpleFake(fakeSecrets...),
		ServiceAccounts: mocks.MockBuilderServiceAccount(mocks.MockBuilderSecrets()),
		Client: generator.Client{
			GetBuildConfigFunc: func(ctx kapi.Context, name string) (*buildapi.BuildConfig, error) {
				return mocks.MockBuildConfig(mocks.MockSource(), mocks.MockSourceStrategyForImageRepository(), mocks.MockOutput()), nil
			},
			UpdateBuildConfigFunc: func(ctx kapi.Context, buildConfig *buildapi.BuildConfig) error {
				return nil
			},
			CreateBuildFunc: func(ctx kapi.Context, build *buildapi.Build) error {
				return nil
			},
			GetBuildFunc: func(ctx kapi.Context, name string) (*buildapi.Build, error) {
				return &buildapi.Build{}, nil
			},
			GetImageStreamFunc: func(ctx kapi.Context, name string) (*imageapi.ImageStream, error) {
				return imageStream, nil
			},
			GetImageStreamTagFunc: func(ctx kapi.Context, name string) (*imageapi.ImageStreamTag, error) {
				return &imageapi.ImageStreamTag{Image: *image}, nil
			},
			GetImageStreamImageFunc: func(ctx kapi.Context, name string) (*imageapi.ImageStreamImage, error) {
				return &imageapi.ImageStreamImage{Image: *image}, nil
			},
		}}}

	_, err := rest.Create(kapi.NewDefaultContext(), &buildapi.BuildRequest{ObjectMeta: kapi.ObjectMeta{Name: "name"}})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestCreateInstantiateValidationError(t *testing.T) {
	rest := InstantiateREST{&generator.BuildGenerator{}}
	_, err := rest.Create(kapi.NewDefaultContext(), &buildapi.BuildRequest{})
	if err == nil {
		t.Error("Expected object got none!")
	}
}
