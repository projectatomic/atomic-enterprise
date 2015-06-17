package generator

import (
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/build/generator"
)

func TestCreateClone(t *testing.T) {
	rest := CloneREST{&generator.BuildGenerator{Client: generator.Client{
		CreateBuildFunc: func(ctx kapi.Context, build *buildapi.Build) error {
			return nil
		},
		GetBuildFunc: func(ctx kapi.Context, name string) (*buildapi.Build, error) {
			return &buildapi.Build{}, nil
		},
	}}}

	_, err := rest.Create(kapi.NewDefaultContext(), &buildapi.BuildRequest{ObjectMeta: kapi.ObjectMeta{Name: "name"}})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestCreateCloneValidationError(t *testing.T) {
	rest := CloneREST{&generator.BuildGenerator{}}
	_, err := rest.Create(kapi.NewDefaultContext(), &buildapi.BuildRequest{})
	if err == nil {
		t.Error("Expected object got none!")
	}
}
