package describe

import (
	"strings"
	"testing"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	ktestclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client/testclient"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	"github.com/projectatomic/appinfra-next/pkg/client/testclient"
	projectapi "github.com/projectatomic/appinfra-next/pkg/project/api"
)

func mustParseTime(t string) time.Time {
	out, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return out
}

func TestProjectStatus(t *testing.T) {
	testCases := map[string]struct {
		Path     string
		Extra    []runtime.Object
		ErrFn    func(error) bool
		Contains []string
		Time     time.Time
	}{
		"missing project": {
			ErrFn: func(err error) bool { return errors.IsNotFound(err) },
		},
		"empty project with display name": {
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{
						Name:      "example",
						Namespace: "",
						Annotations: map[string]string{
							projectapi.ProjectDisplayName: "Test",
						},
					},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project Test (example)\n",
				"You have no Services, DeploymentConfigs, or BuildConfigs.",
			},
		},
		"empty service": {
			Path: "../../../../test/fixtures/app-scenarios/k8s-service-with-nothing.json",
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{Name: "example", Namespace: ""},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project example\n",
				"service empty-service",
				"(<initializing>:5432", "To see more information",
			},
		},
		"unstarted build": {
			Path: "../../../../test/fixtures/app-scenarios/new-project-no-build.yaml",
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{Name: "example", Namespace: ""},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project example\n",
				"service sinatra-example-2 (172.30.17.48:8080)",
				"builds git://github.com",
				"with docker.io/openshift/ruby-20-centos7:latest",
				"not built yet",
				"#1 deployment waiting on image or update",
				"To see more information",
			},
		},
		"running build": {
			Path: "../../../../test/fixtures/app-scenarios/new-project-one-build.yaml",
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{Name: "example", Namespace: ""},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project example\n",
				"service sinatra-example-1 (172.30.17.47:8080)",
				"builds git://github.com",
				"with docker.io/openshift/ruby-20-centos7:latest",
				"build 1 running for about a minute",
				"#1 deployment waiting on image or update",
				"To see more information",
			},
			Time: mustParseTime("2015-04-06T21:20:03Z"),
		},
		"a/b test DeploymentConfig": {
			Path: "../../../../test/fixtures/app-scenarios/new-project-two-deployment-configs.yaml",
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{Name: "example", Namespace: ""},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project example\n",
				"service sinatra-app-example (172.30.17.49:8080)",
				"sinatra-app-example-a deploys",
				"sinatra-app-example-b deploys",
				"with docker.io/openshift/ruby-20-centos7:latest",
				"build 1 running for about a minute",
				"- 7a4f354: Prepare v1beta3 Template types (Roy Programmer <someguy@outhere.com>)",
				"To see more information",
			},
			Time: mustParseTime("2015-04-06T21:20:03Z"),
		},
		"with real deployments": {
			Path: "../../../../test/fixtures/app-scenarios/new-project-deployed-app.yaml",
			Extra: []runtime.Object{
				&projectapi.Project{
					ObjectMeta: kapi.ObjectMeta{Name: "example", Namespace: ""},
				},
			},
			ErrFn: func(err error) bool { return err == nil },
			Contains: []string{
				"In project example\n",
				"service database (172.30.17.240:5434 -> 3306)",
				"service frontend (172.30.17.154:5432 -> 8080)",
				"database deploys",
				"frontend deploys",
				"with docker.io/openshift/ruby-20-centos7:latest",
				"#2 deployment failed less than a second ago: unable to contact server - 0/1 pods",
				"#2 deployment running for 7 seconds - 2/1 pods",
				"#1 deployed 8 seconds ago",
				"#1 deployed less than a second ago",
				"To see more information",
			},
			Time: mustParseTime("2015-04-07T04:12:25Z"),
		},
	}
	oldTimeFn := timeNowFn
	defer func() { timeNowFn = oldTimeFn }()
	for k, test := range testCases {
		timeNowFn = func() time.Time {
			if !test.Time.IsZero() {
				return test.Time
			}
			return time.Now()
		}
		o := ktestclient.NewObjects(kapi.Scheme, kapi.Scheme)
		if len(test.Path) > 0 {
			if err := ktestclient.AddObjectsFromPath(test.Path, o, kapi.Scheme); err != nil {
				t.Fatal(err)
			}
		}
		for _, obj := range test.Extra {
			o.Add(obj)
		}
		oc, kc := testclient.NewFixtureClients(o)
		d := ProjectStatusDescriber{C: oc, K: kc}
		out, err := d.Describe("example", "")
		if !test.ErrFn(err) {
			t.Errorf("%s: unexpected error: %v", k, err)
		}
		if err != nil {
			continue
		}
		for _, s := range test.Contains {
			if !strings.Contains(out, s) {
				t.Errorf("%s: did not have %q:\n%s\n---", k, s, out)
			}
		}
		t.Logf("\n%s", out)
	}
}
