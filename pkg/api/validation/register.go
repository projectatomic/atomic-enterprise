package validation

import (
	authorizationvalidation "github.com/projectatomic/appinfra-next/pkg/authorization/api/validation"
	buildvalidation "github.com/projectatomic/appinfra-next/pkg/build/api/validation"
	deployvalidation "github.com/projectatomic/appinfra-next/pkg/deploy/api/validation"
	imagevalidation "github.com/projectatomic/appinfra-next/pkg/image/api/validation"
	oauthvalidation "github.com/projectatomic/appinfra-next/pkg/oauth/api/validation"
	projectvalidation "github.com/projectatomic/appinfra-next/pkg/project/api/validation"
	routevalidation "github.com/projectatomic/appinfra-next/pkg/route/api/validation"
	sdnvalidation "github.com/projectatomic/appinfra-next/pkg/sdn/api/validation"
	templatevalidation "github.com/projectatomic/appinfra-next/pkg/template/api/validation"
	uservalidation "github.com/projectatomic/appinfra-next/pkg/user/api/validation"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	imageapi "github.com/projectatomic/appinfra-next/pkg/image/api"
	oauthapi "github.com/projectatomic/appinfra-next/pkg/oauth/api"
	projectapi "github.com/projectatomic/appinfra-next/pkg/project/api"
	routeapi "github.com/projectatomic/appinfra-next/pkg/route/api"
	sdnapi "github.com/projectatomic/appinfra-next/pkg/sdn/api"
	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
	userapi "github.com/projectatomic/appinfra-next/pkg/user/api"
)

func init() {
	Validator.Register(&authorizationapi.SubjectAccessReview{}, authorizationvalidation.ValidateSubjectAccessReview, nil)
	Validator.Register(&authorizationapi.ResourceAccessReview{}, authorizationvalidation.ValidateResourceAccessReview, nil)

	Validator.Register(&authorizationapi.Policy{}, authorizationvalidation.ValidateLocalPolicy, authorizationvalidation.ValidateLocalPolicyUpdate)
	Validator.Register(&authorizationapi.PolicyBinding{}, authorizationvalidation.ValidateLocalPolicyBinding, authorizationvalidation.ValidateLocalPolicyBindingUpdate)
	Validator.Register(&authorizationapi.Role{}, authorizationvalidation.ValidateLocalRole, authorizationvalidation.ValidateLocalRoleUpdate)
	Validator.Register(&authorizationapi.RoleBinding{}, authorizationvalidation.ValidateLocalRoleBinding, authorizationvalidation.ValidateLocalRoleBindingUpdate)

	Validator.Register(&authorizationapi.ClusterPolicy{}, authorizationvalidation.ValidateClusterPolicy, authorizationvalidation.ValidateClusterPolicyUpdate)
	Validator.Register(&authorizationapi.ClusterPolicyBinding{}, authorizationvalidation.ValidateClusterPolicyBinding, authorizationvalidation.ValidateClusterPolicyBindingUpdate)
	Validator.Register(&authorizationapi.ClusterRole{}, authorizationvalidation.ValidateClusterRole, authorizationvalidation.ValidateClusterRoleUpdate)
	Validator.Register(&authorizationapi.ClusterRoleBinding{}, authorizationvalidation.ValidateClusterRoleBinding, authorizationvalidation.ValidateClusterRoleBindingUpdate)

	Validator.Register(&buildapi.Build{}, buildvalidation.ValidateBuild, buildvalidation.ValidateBuildUpdate)
	Validator.Register(&buildapi.BuildConfig{}, buildvalidation.ValidateBuildConfig, buildvalidation.ValidateBuildConfigUpdate)
	Validator.Register(&buildapi.BuildRequest{}, buildvalidation.ValidateBuildRequest, nil)

	Validator.Register(&deployapi.DeploymentConfig{}, deployvalidation.ValidateDeploymentConfig, deployvalidation.ValidateDeploymentConfigUpdate)
	Validator.Register(&deployapi.DeploymentConfigRollback{}, deployvalidation.ValidateDeploymentConfigRollback, nil)

	Validator.Register(&imageapi.Image{}, imagevalidation.ValidateImage, nil)
	Validator.Register(&imageapi.ImageStream{}, imagevalidation.ValidateImageStream, imagevalidation.ValidateImageStreamUpdate)
	Validator.Register(&imageapi.ImageStreamMapping{}, imagevalidation.ValidateImageStreamMapping, nil)

	Validator.Register(&oauthapi.OAuthAccessToken{}, oauthvalidation.ValidateAccessToken, nil)
	Validator.Register(&oauthapi.OAuthAuthorizeToken{}, oauthvalidation.ValidateAuthorizeToken, nil)
	Validator.Register(&oauthapi.OAuthClient{}, oauthvalidation.ValidateClient, oauthvalidation.ValidateClientUpdate)
	Validator.Register(&oauthapi.OAuthClientAuthorization{}, oauthvalidation.ValidateClientAuthorization, oauthvalidation.ValidateClientAuthorizationUpdate)

	Validator.Register(&projectapi.Project{}, projectvalidation.ValidateProject, projectvalidation.ValidateProjectUpdate)
	Validator.Register(&projectapi.ProjectRequest{}, projectvalidation.ValidateProjectRequest, nil)

	Validator.Register(&routeapi.Route{}, routevalidation.ValidateRoute, routevalidation.ValidateRouteUpdate)

	Validator.Register(&sdnapi.ClusterNetwork{}, sdnvalidation.ValidateClusterNetwork, sdnvalidation.ValidateClusterNetworkUpdate)
	Validator.Register(&sdnapi.HostSubnet{}, sdnvalidation.ValidateHostSubnet, sdnvalidation.ValidateHostSubnetUpdate)

	Validator.Register(&templateapi.Template{}, templatevalidation.ValidateTemplate, templatevalidation.ValidateTemplateUpdate)

	Validator.Register(&userapi.User{}, uservalidation.ValidateUser, uservalidation.ValidateUserUpdate)
	Validator.Register(&userapi.Identity{}, uservalidation.ValidateIdentity, uservalidation.ValidateIdentityUpdate)
	Validator.Register(&userapi.UserIdentityMapping{}, uservalidation.ValidateUserIdentityMapping, uservalidation.ValidateUserIdentityMappingUpdate)
}
