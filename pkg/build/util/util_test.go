package util

import (
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

func TestGetBuildPodName(t *testing.T) {
	if expected, actual := "mybuild-build", GetBuildPodName(&buildapi.Build{ObjectMeta: kapi.ObjectMeta{Name: "mybuild"}}); expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
