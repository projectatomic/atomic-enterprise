package delegated

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	"github.com/projectatomic/appinfra-next/pkg/cmd/server/bootstrappolicy"
	projectapi "github.com/projectatomic/appinfra-next/pkg/project/api"
	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
)

const (
	DefaultTemplateName = "project-request"

	ProjectNameParam        = "PROJECT_NAME"
	ProjectDisplayNameParam = "PROJECT_DISPLAYNAME"
	ProjectDescriptionParam = "PROJECT_DESCRIPTION"
	ProjectAdminUserParam   = "PROJECT_ADMIN_USER"
)

var (
	parameters = []string{ProjectNameParam, ProjectDisplayNameParam, ProjectDescriptionParam, ProjectAdminUserParam}
)

func DefaultTemplate() *templateapi.Template {
	ret := &templateapi.Template{}
	ret.Name = DefaultTemplateName

	ns := "${" + ProjectNameParam + "}"

	project := &projectapi.Project{}
	project.Name = ns
	project.Annotations = map[string]string{
		projectapi.ProjectDescription: "${" + ProjectDescriptionParam + "}",
		projectapi.ProjectDisplayName: "${" + ProjectDisplayNameParam + "}",
	}
	ret.Objects = append(ret.Objects, project)

	binding := &authorizationapi.RoleBinding{}
	binding.Name = "admins"
	binding.Namespace = ns
	binding.Users = util.NewStringSet("${" + ProjectAdminUserParam + "}")
	binding.RoleRef.Name = bootstrappolicy.AdminRoleName
	ret.Objects = append(ret.Objects, binding)

	serviceAccountRoleBindings := bootstrappolicy.GetBootstrapServiceAccountProjectRoleBindings(ns)
	for i := range serviceAccountRoleBindings {
		ret.Objects = append(ret.Objects, &serviceAccountRoleBindings[i])
	}

	for _, parameterName := range parameters {
		parameter := templateapi.Parameter{}
		parameter.Name = parameterName
		ret.Parameters = append(ret.Parameters, parameter)
	}

	return ret
}
