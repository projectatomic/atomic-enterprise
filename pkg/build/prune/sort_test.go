package prune

import (
	"sort"
	"testing"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

// TestSort verifies that builds are sorted by most recently created
func TestSort(t *testing.T) {
	present := util.Now()
	past := util.NewTime(present.Time.Add(-1 * time.Minute))
	builds := []*buildapi.Build{
		{
			ObjectMeta: kapi.ObjectMeta{
				Name:              "past",
				CreationTimestamp: past,
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name:              "present",
				CreationTimestamp: present,
			},
		},
	}
	sort.Sort(sortableBuilds(builds))
	if builds[0].Name != "present" {
		t.Errorf("Unexpected sort order")
	}
	if builds[1].Name != "past" {
		t.Errorf("Unexpected sort order")
	}
}
