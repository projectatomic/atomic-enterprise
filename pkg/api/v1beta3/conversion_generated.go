package v1beta3

// AUTO-GENERATED FUNCTIONS START HERE
import (
	api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	v1beta3 "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3"
	conversion "github.com/GoogleCloudPlatform/kubernetes/pkg/conversion"
	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	apiv1beta3 "github.com/projectatomic/appinfra-next/pkg/authorization/api/v1beta3"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	buildapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/build/api/v1beta3"
	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	deployapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/deploy/api/v1beta3"
	imageapi "github.com/projectatomic/appinfra-next/pkg/image/api"
	imageapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/image/api/v1beta3"
	oauthapi "github.com/projectatomic/appinfra-next/pkg/oauth/api"
	oauthapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/oauth/api/v1beta3"
	projectapi "github.com/projectatomic/appinfra-next/pkg/project/api"
	projectapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/project/api/v1beta3"
	routeapi "github.com/projectatomic/appinfra-next/pkg/route/api"
	routeapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/route/api/v1beta3"
	sdnapi "github.com/projectatomic/appinfra-next/pkg/sdn/api"
	sdnapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/sdn/api/v1beta3"
	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
	templateapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/template/api/v1beta3"
	userapi "github.com/projectatomic/appinfra-next/pkg/user/api"
	userapiv1beta3 "github.com/projectatomic/appinfra-next/pkg/user/api/v1beta3"
	reflect "reflect"
)

func convert_api_ListMeta_To_v1beta3_ListMeta(in *api.ListMeta, out *v1beta3.ListMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.ListMeta))(in)
	}
	out.SelfLink = in.SelfLink
	out.ResourceVersion = in.ResourceVersion
	return nil
}

func convert_api_ObjectMeta_To_v1beta3_ObjectMeta(in *api.ObjectMeta, out *v1beta3.ObjectMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.ObjectMeta))(in)
	}
	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = in.UID
	out.ResourceVersion = in.ResourceVersion
	if err := s.Convert(&in.CreationTimestamp, &out.CreationTimestamp, 0); err != nil {
		return err
	}
	if in.DeletionTimestamp != nil {
		if err := s.Convert(&in.DeletionTimestamp, &out.DeletionTimestamp, 0); err != nil {
			return err
		}
	} else {
		out.DeletionTimestamp = nil
	}
	if in.Labels != nil {
		out.Labels = make(map[string]string)
		for key, val := range in.Labels {
			out.Labels[key] = val
		}
	} else {
		out.Labels = nil
	}
	if in.Annotations != nil {
		out.Annotations = make(map[string]string)
		for key, val := range in.Annotations {
			out.Annotations[key] = val
		}
	} else {
		out.Annotations = nil
	}
	return nil
}

func convert_api_ObjectReference_To_v1beta3_ObjectReference(in *api.ObjectReference, out *v1beta3.ObjectReference, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.ObjectReference))(in)
	}
	out.Kind = in.Kind
	out.Namespace = in.Namespace
	out.Name = in.Name
	out.UID = in.UID
	out.APIVersion = in.APIVersion
	out.ResourceVersion = in.ResourceVersion
	out.FieldPath = in.FieldPath
	return nil
}

func convert_api_TypeMeta_To_v1beta3_TypeMeta(in *api.TypeMeta, out *v1beta3.TypeMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.TypeMeta))(in)
	}
	out.Kind = in.Kind
	out.APIVersion = in.APIVersion
	return nil
}

func convert_v1beta3_ListMeta_To_api_ListMeta(in *v1beta3.ListMeta, out *api.ListMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1beta3.ListMeta))(in)
	}
	out.SelfLink = in.SelfLink
	out.ResourceVersion = in.ResourceVersion
	return nil
}

func convert_v1beta3_ObjectMeta_To_api_ObjectMeta(in *v1beta3.ObjectMeta, out *api.ObjectMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1beta3.ObjectMeta))(in)
	}
	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = in.UID
	out.ResourceVersion = in.ResourceVersion
	if err := s.Convert(&in.CreationTimestamp, &out.CreationTimestamp, 0); err != nil {
		return err
	}
	if in.DeletionTimestamp != nil {
		if err := s.Convert(&in.DeletionTimestamp, &out.DeletionTimestamp, 0); err != nil {
			return err
		}
	} else {
		out.DeletionTimestamp = nil
	}
	if in.Labels != nil {
		out.Labels = make(map[string]string)
		for key, val := range in.Labels {
			out.Labels[key] = val
		}
	} else {
		out.Labels = nil
	}
	if in.Annotations != nil {
		out.Annotations = make(map[string]string)
		for key, val := range in.Annotations {
			out.Annotations[key] = val
		}
	} else {
		out.Annotations = nil
	}
	return nil
}

func convert_v1beta3_ObjectReference_To_api_ObjectReference(in *v1beta3.ObjectReference, out *api.ObjectReference, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1beta3.ObjectReference))(in)
	}
	out.Kind = in.Kind
	out.Namespace = in.Namespace
	out.Name = in.Name
	out.UID = in.UID
	out.APIVersion = in.APIVersion
	out.ResourceVersion = in.ResourceVersion
	out.FieldPath = in.FieldPath
	return nil
}

func convert_v1beta3_TypeMeta_To_api_TypeMeta(in *v1beta3.TypeMeta, out *api.TypeMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1beta3.TypeMeta))(in)
	}
	out.Kind = in.Kind
	out.APIVersion = in.APIVersion
	return nil
}

func convert_api_ClusterPolicy_To_v1beta3_ClusterPolicy(in *authorizationapi.ClusterPolicy, out *apiv1beta3.ClusterPolicy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterPolicy))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.LastModified, &out.LastModified, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.Roles, &out.Roles, 0); err != nil {
		return err
	}
	return nil
}

func convert_api_ClusterPolicyBindingList_To_v1beta3_ClusterPolicyBindingList(in *authorizationapi.ClusterPolicyBindingList, out *apiv1beta3.ClusterPolicyBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterPolicyBindingList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.ClusterPolicyBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ClusterPolicyList_To_v1beta3_ClusterPolicyList(in *authorizationapi.ClusterPolicyList, out *apiv1beta3.ClusterPolicyList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterPolicyList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.ClusterPolicy, len(in.Items))
		for i := range in.Items {
			if err := convert_api_ClusterPolicy_To_v1beta3_ClusterPolicy(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ClusterRole_To_v1beta3_ClusterRole(in *authorizationapi.ClusterRole, out *apiv1beta3.ClusterRole, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterRole))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Rules != nil {
		out.Rules = make([]apiv1beta3.PolicyRule, len(in.Rules))
		for i := range in.Rules {
			if err := s.Convert(&in.Rules[i], &out.Rules[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	return nil
}

func convert_api_ClusterRoleBindingList_To_v1beta3_ClusterRoleBindingList(in *authorizationapi.ClusterRoleBindingList, out *apiv1beta3.ClusterRoleBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterRoleBindingList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.ClusterRoleBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ClusterRoleList_To_v1beta3_ClusterRoleList(in *authorizationapi.ClusterRoleList, out *apiv1beta3.ClusterRoleList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ClusterRoleList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.ClusterRole, len(in.Items))
		for i := range in.Items {
			if err := convert_api_ClusterRole_To_v1beta3_ClusterRole(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_IsPersonalSubjectAccessReview_To_v1beta3_IsPersonalSubjectAccessReview(in *authorizationapi.IsPersonalSubjectAccessReview, out *apiv1beta3.IsPersonalSubjectAccessReview, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.IsPersonalSubjectAccessReview))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	return nil
}

func convert_api_PolicyBindingList_To_v1beta3_PolicyBindingList(in *authorizationapi.PolicyBindingList, out *apiv1beta3.PolicyBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.PolicyBindingList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.PolicyBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_PolicyList_To_v1beta3_PolicyList(in *authorizationapi.PolicyList, out *apiv1beta3.PolicyList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.PolicyList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.Policy, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ResourceAccessReview_To_v1beta3_ResourceAccessReview(in *authorizationapi.ResourceAccessReview, out *apiv1beta3.ResourceAccessReview, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.ResourceAccessReview))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Verb = in.Verb
	out.Resource = in.Resource
	if err := s.Convert(&in.Content, &out.Content, 0); err != nil {
		return err
	}
	out.ResourceName = in.ResourceName
	return nil
}

func convert_api_Role_To_v1beta3_Role(in *authorizationapi.Role, out *apiv1beta3.Role, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.Role))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Rules != nil {
		out.Rules = make([]apiv1beta3.PolicyRule, len(in.Rules))
		for i := range in.Rules {
			if err := s.Convert(&in.Rules[i], &out.Rules[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	return nil
}

func convert_api_RoleBindingList_To_v1beta3_RoleBindingList(in *authorizationapi.RoleBindingList, out *apiv1beta3.RoleBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.RoleBindingList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.RoleBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_RoleList_To_v1beta3_RoleList(in *authorizationapi.RoleList, out *apiv1beta3.RoleList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.RoleList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]apiv1beta3.Role, len(in.Items))
		for i := range in.Items {
			if err := convert_api_Role_To_v1beta3_Role(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_SubjectAccessReviewResponse_To_v1beta3_SubjectAccessReviewResponse(in *authorizationapi.SubjectAccessReviewResponse, out *apiv1beta3.SubjectAccessReviewResponse, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*authorizationapi.SubjectAccessReviewResponse))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Namespace = in.Namespace
	out.Allowed = in.Allowed
	out.Reason = in.Reason
	return nil
}

func convert_v1beta3_ClusterPolicy_To_api_ClusterPolicy(in *apiv1beta3.ClusterPolicy, out *authorizationapi.ClusterPolicy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterPolicy))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := s.Convert(&in.LastModified, &out.LastModified, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.Roles, &out.Roles, 0); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_ClusterPolicyBindingList_To_api_ClusterPolicyBindingList(in *apiv1beta3.ClusterPolicyBindingList, out *authorizationapi.ClusterPolicyBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterPolicyBindingList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.ClusterPolicyBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ClusterPolicyList_To_api_ClusterPolicyList(in *apiv1beta3.ClusterPolicyList, out *authorizationapi.ClusterPolicyList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterPolicyList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.ClusterPolicy, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_ClusterPolicy_To_api_ClusterPolicy(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ClusterRole_To_api_ClusterRole(in *apiv1beta3.ClusterRole, out *authorizationapi.ClusterRole, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterRole))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Rules != nil {
		out.Rules = make([]authorizationapi.PolicyRule, len(in.Rules))
		for i := range in.Rules {
			if err := s.Convert(&in.Rules[i], &out.Rules[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	return nil
}

func convert_v1beta3_ClusterRoleBindingList_To_api_ClusterRoleBindingList(in *apiv1beta3.ClusterRoleBindingList, out *authorizationapi.ClusterRoleBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterRoleBindingList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.ClusterRoleBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ClusterRoleList_To_api_ClusterRoleList(in *apiv1beta3.ClusterRoleList, out *authorizationapi.ClusterRoleList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ClusterRoleList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.ClusterRole, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_ClusterRole_To_api_ClusterRole(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_IsPersonalSubjectAccessReview_To_api_IsPersonalSubjectAccessReview(in *apiv1beta3.IsPersonalSubjectAccessReview, out *authorizationapi.IsPersonalSubjectAccessReview, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.IsPersonalSubjectAccessReview))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_PolicyBindingList_To_api_PolicyBindingList(in *apiv1beta3.PolicyBindingList, out *authorizationapi.PolicyBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.PolicyBindingList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.PolicyBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_PolicyList_To_api_PolicyList(in *apiv1beta3.PolicyList, out *authorizationapi.PolicyList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.PolicyList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.Policy, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ResourceAccessReview_To_api_ResourceAccessReview(in *apiv1beta3.ResourceAccessReview, out *authorizationapi.ResourceAccessReview, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.ResourceAccessReview))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Verb = in.Verb
	out.Resource = in.Resource
	if err := s.Convert(&in.Content, &out.Content, 0); err != nil {
		return err
	}
	out.ResourceName = in.ResourceName
	return nil
}

func convert_v1beta3_Role_To_api_Role(in *apiv1beta3.Role, out *authorizationapi.Role, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.Role))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Rules != nil {
		out.Rules = make([]authorizationapi.PolicyRule, len(in.Rules))
		for i := range in.Rules {
			if err := s.Convert(&in.Rules[i], &out.Rules[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	return nil
}

func convert_v1beta3_RoleBindingList_To_api_RoleBindingList(in *apiv1beta3.RoleBindingList, out *authorizationapi.RoleBindingList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.RoleBindingList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.RoleBinding, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_RoleList_To_api_RoleList(in *apiv1beta3.RoleList, out *authorizationapi.RoleList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.RoleList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]authorizationapi.Role, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_Role_To_api_Role(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_SubjectAccessReviewResponse_To_api_SubjectAccessReviewResponse(in *apiv1beta3.SubjectAccessReviewResponse, out *authorizationapi.SubjectAccessReviewResponse, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*apiv1beta3.SubjectAccessReviewResponse))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Namespace = in.Namespace
	out.Allowed = in.Allowed
	out.Reason = in.Reason
	return nil
}

func convert_api_BuildConfigList_To_v1beta3_BuildConfigList(in *buildapi.BuildConfigList, out *buildapiv1beta3.BuildConfigList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.BuildConfigList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]buildapiv1beta3.BuildConfig, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_BuildList_To_v1beta3_BuildList(in *buildapi.BuildList, out *buildapiv1beta3.BuildList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.BuildList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]buildapiv1beta3.Build, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_BuildLog_To_v1beta3_BuildLog(in *buildapi.BuildLog, out *buildapiv1beta3.BuildLog, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.BuildLog))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	return nil
}

func convert_api_BuildLogOptions_To_v1beta3_BuildLogOptions(in *buildapi.BuildLogOptions, out *buildapiv1beta3.BuildLogOptions, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.BuildLogOptions))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Follow = in.Follow
	out.NoWait = in.NoWait
	return nil
}

func convert_api_BuildRequest_To_v1beta3_BuildRequest(in *buildapi.BuildRequest, out *buildapiv1beta3.BuildRequest, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.BuildRequest))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Revision != nil {
		out.Revision = new(buildapiv1beta3.SourceRevision)
		if err := convert_api_SourceRevision_To_v1beta3_SourceRevision(in.Revision, out.Revision, s); err != nil {
			return err
		}
	} else {
		out.Revision = nil
	}
	if in.TriggeredByImage != nil {
		out.TriggeredByImage = new(v1beta3.ObjectReference)
		if err := convert_api_ObjectReference_To_v1beta3_ObjectReference(in.TriggeredByImage, out.TriggeredByImage, s); err != nil {
			return err
		}
	} else {
		out.TriggeredByImage = nil
	}
	return nil
}

func convert_api_GitSourceRevision_To_v1beta3_GitSourceRevision(in *buildapi.GitSourceRevision, out *buildapiv1beta3.GitSourceRevision, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.GitSourceRevision))(in)
	}
	out.Commit = in.Commit
	if err := convert_api_SourceControlUser_To_v1beta3_SourceControlUser(&in.Author, &out.Author, s); err != nil {
		return err
	}
	if err := convert_api_SourceControlUser_To_v1beta3_SourceControlUser(&in.Committer, &out.Committer, s); err != nil {
		return err
	}
	out.Message = in.Message
	return nil
}

func convert_api_SourceControlUser_To_v1beta3_SourceControlUser(in *buildapi.SourceControlUser, out *buildapiv1beta3.SourceControlUser, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.SourceControlUser))(in)
	}
	out.Name = in.Name
	out.Email = in.Email
	return nil
}

func convert_api_SourceRevision_To_v1beta3_SourceRevision(in *buildapi.SourceRevision, out *buildapiv1beta3.SourceRevision, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapi.SourceRevision))(in)
	}
	out.Type = buildapiv1beta3.BuildSourceType(in.Type)
	if in.Git != nil {
		out.Git = new(buildapiv1beta3.GitSourceRevision)
		if err := convert_api_GitSourceRevision_To_v1beta3_GitSourceRevision(in.Git, out.Git, s); err != nil {
			return err
		}
	} else {
		out.Git = nil
	}
	return nil
}

func convert_v1beta3_BuildConfigList_To_api_BuildConfigList(in *buildapiv1beta3.BuildConfigList, out *buildapi.BuildConfigList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.BuildConfigList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]buildapi.BuildConfig, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_BuildList_To_api_BuildList(in *buildapiv1beta3.BuildList, out *buildapi.BuildList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.BuildList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]buildapi.Build, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_BuildLog_To_api_BuildLog(in *buildapiv1beta3.BuildLog, out *buildapi.BuildLog, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.BuildLog))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_BuildLogOptions_To_api_BuildLogOptions(in *buildapiv1beta3.BuildLogOptions, out *buildapi.BuildLogOptions, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.BuildLogOptions))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	out.Follow = in.Follow
	out.NoWait = in.NoWait
	return nil
}

func convert_v1beta3_BuildRequest_To_api_BuildRequest(in *buildapiv1beta3.BuildRequest, out *buildapi.BuildRequest, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.BuildRequest))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if in.Revision != nil {
		out.Revision = new(buildapi.SourceRevision)
		if err := convert_v1beta3_SourceRevision_To_api_SourceRevision(in.Revision, out.Revision, s); err != nil {
			return err
		}
	} else {
		out.Revision = nil
	}
	if in.TriggeredByImage != nil {
		out.TriggeredByImage = new(api.ObjectReference)
		if err := convert_v1beta3_ObjectReference_To_api_ObjectReference(in.TriggeredByImage, out.TriggeredByImage, s); err != nil {
			return err
		}
	} else {
		out.TriggeredByImage = nil
	}
	return nil
}

func convert_v1beta3_GitSourceRevision_To_api_GitSourceRevision(in *buildapiv1beta3.GitSourceRevision, out *buildapi.GitSourceRevision, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.GitSourceRevision))(in)
	}
	out.Commit = in.Commit
	if err := convert_v1beta3_SourceControlUser_To_api_SourceControlUser(&in.Author, &out.Author, s); err != nil {
		return err
	}
	if err := convert_v1beta3_SourceControlUser_To_api_SourceControlUser(&in.Committer, &out.Committer, s); err != nil {
		return err
	}
	out.Message = in.Message
	return nil
}

func convert_v1beta3_SourceControlUser_To_api_SourceControlUser(in *buildapiv1beta3.SourceControlUser, out *buildapi.SourceControlUser, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.SourceControlUser))(in)
	}
	out.Name = in.Name
	out.Email = in.Email
	return nil
}

func convert_v1beta3_SourceRevision_To_api_SourceRevision(in *buildapiv1beta3.SourceRevision, out *buildapi.SourceRevision, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*buildapiv1beta3.SourceRevision))(in)
	}
	out.Type = buildapi.BuildSourceType(in.Type)
	if in.Git != nil {
		out.Git = new(buildapi.GitSourceRevision)
		if err := convert_v1beta3_GitSourceRevision_To_api_GitSourceRevision(in.Git, out.Git, s); err != nil {
			return err
		}
	} else {
		out.Git = nil
	}
	return nil
}

func convert_api_DeploymentConfigList_To_v1beta3_DeploymentConfigList(in *deployapi.DeploymentConfigList, out *deployapiv1beta3.DeploymentConfigList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapi.DeploymentConfigList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]deployapiv1beta3.DeploymentConfig, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_DeploymentConfigRollback_To_v1beta3_DeploymentConfigRollback(in *deployapi.DeploymentConfigRollback, out *deployapiv1beta3.DeploymentConfigRollback, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapi.DeploymentConfigRollback))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_DeploymentConfigRollbackSpec_To_v1beta3_DeploymentConfigRollbackSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func convert_api_DeploymentConfigRollbackSpec_To_v1beta3_DeploymentConfigRollbackSpec(in *deployapi.DeploymentConfigRollbackSpec, out *deployapiv1beta3.DeploymentConfigRollbackSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapi.DeploymentConfigRollbackSpec))(in)
	}
	if err := convert_api_ObjectReference_To_v1beta3_ObjectReference(&in.From, &out.From, s); err != nil {
		return err
	}
	out.IncludeTriggers = in.IncludeTriggers
	out.IncludeTemplate = in.IncludeTemplate
	out.IncludeReplicationMeta = in.IncludeReplicationMeta
	out.IncludeStrategy = in.IncludeStrategy
	return nil
}

func convert_v1beta3_DeploymentConfigList_To_api_DeploymentConfigList(in *deployapiv1beta3.DeploymentConfigList, out *deployapi.DeploymentConfigList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapiv1beta3.DeploymentConfigList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]deployapi.DeploymentConfig, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_DeploymentConfigRollback_To_api_DeploymentConfigRollback(in *deployapiv1beta3.DeploymentConfigRollback, out *deployapi.DeploymentConfigRollback, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapiv1beta3.DeploymentConfigRollback))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_DeploymentConfigRollbackSpec_To_api_DeploymentConfigRollbackSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_DeploymentConfigRollbackSpec_To_api_DeploymentConfigRollbackSpec(in *deployapiv1beta3.DeploymentConfigRollbackSpec, out *deployapi.DeploymentConfigRollbackSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*deployapiv1beta3.DeploymentConfigRollbackSpec))(in)
	}
	if err := convert_v1beta3_ObjectReference_To_api_ObjectReference(&in.From, &out.From, s); err != nil {
		return err
	}
	out.IncludeTriggers = in.IncludeTriggers
	out.IncludeTemplate = in.IncludeTemplate
	out.IncludeReplicationMeta = in.IncludeReplicationMeta
	out.IncludeStrategy = in.IncludeStrategy
	return nil
}

func convert_api_ImageList_To_v1beta3_ImageList(in *imageapi.ImageList, out *imageapiv1beta3.ImageList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*imageapi.ImageList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]imageapiv1beta3.Image, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ImageStreamList_To_v1beta3_ImageStreamList(in *imageapi.ImageStreamList, out *imageapiv1beta3.ImageStreamList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*imageapi.ImageStreamList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]imageapiv1beta3.ImageStream, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ImageList_To_api_ImageList(in *imageapiv1beta3.ImageList, out *imageapi.ImageList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*imageapiv1beta3.ImageList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]imageapi.Image, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ImageStreamList_To_api_ImageStreamList(in *imageapiv1beta3.ImageStreamList, out *imageapi.ImageStreamList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*imageapiv1beta3.ImageStreamList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]imageapi.ImageStream, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_OAuthAccessToken_To_v1beta3_OAuthAccessToken(in *oauthapi.OAuthAccessToken, out *oauthapiv1beta3.OAuthAccessToken, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthAccessToken))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.ExpiresIn = in.ExpiresIn
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	out.RedirectURI = in.RedirectURI
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	out.AuthorizeToken = in.AuthorizeToken
	out.RefreshToken = in.RefreshToken
	return nil
}

func convert_api_OAuthAccessTokenList_To_v1beta3_OAuthAccessTokenList(in *oauthapi.OAuthAccessTokenList, out *oauthapiv1beta3.OAuthAccessTokenList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthAccessTokenList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapiv1beta3.OAuthAccessToken, len(in.Items))
		for i := range in.Items {
			if err := convert_api_OAuthAccessToken_To_v1beta3_OAuthAccessToken(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_OAuthAuthorizeToken_To_v1beta3_OAuthAuthorizeToken(in *oauthapi.OAuthAuthorizeToken, out *oauthapiv1beta3.OAuthAuthorizeToken, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthAuthorizeToken))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.ExpiresIn = in.ExpiresIn
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	out.RedirectURI = in.RedirectURI
	out.State = in.State
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	return nil
}

func convert_api_OAuthAuthorizeTokenList_To_v1beta3_OAuthAuthorizeTokenList(in *oauthapi.OAuthAuthorizeTokenList, out *oauthapiv1beta3.OAuthAuthorizeTokenList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthAuthorizeTokenList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapiv1beta3.OAuthAuthorizeToken, len(in.Items))
		for i := range in.Items {
			if err := convert_api_OAuthAuthorizeToken_To_v1beta3_OAuthAuthorizeToken(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_OAuthClient_To_v1beta3_OAuthClient(in *oauthapi.OAuthClient, out *oauthapiv1beta3.OAuthClient, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthClient))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Secret = in.Secret
	out.RespondWithChallenges = in.RespondWithChallenges
	if in.RedirectURIs != nil {
		out.RedirectURIs = make([]string, len(in.RedirectURIs))
		for i := range in.RedirectURIs {
			out.RedirectURIs[i] = in.RedirectURIs[i]
		}
	} else {
		out.RedirectURIs = nil
	}
	return nil
}

func convert_api_OAuthClientAuthorization_To_v1beta3_OAuthClientAuthorization(in *oauthapi.OAuthClientAuthorization, out *oauthapiv1beta3.OAuthClientAuthorization, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthClientAuthorization))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	return nil
}

func convert_api_OAuthClientAuthorizationList_To_v1beta3_OAuthClientAuthorizationList(in *oauthapi.OAuthClientAuthorizationList, out *oauthapiv1beta3.OAuthClientAuthorizationList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthClientAuthorizationList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapiv1beta3.OAuthClientAuthorization, len(in.Items))
		for i := range in.Items {
			if err := convert_api_OAuthClientAuthorization_To_v1beta3_OAuthClientAuthorization(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_OAuthClientList_To_v1beta3_OAuthClientList(in *oauthapi.OAuthClientList, out *oauthapiv1beta3.OAuthClientList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapi.OAuthClientList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapiv1beta3.OAuthClient, len(in.Items))
		for i := range in.Items {
			if err := convert_api_OAuthClient_To_v1beta3_OAuthClient(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_OAuthAccessToken_To_api_OAuthAccessToken(in *oauthapiv1beta3.OAuthAccessToken, out *oauthapi.OAuthAccessToken, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthAccessToken))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.ExpiresIn = in.ExpiresIn
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	out.RedirectURI = in.RedirectURI
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	out.AuthorizeToken = in.AuthorizeToken
	out.RefreshToken = in.RefreshToken
	return nil
}

func convert_v1beta3_OAuthAccessTokenList_To_api_OAuthAccessTokenList(in *oauthapiv1beta3.OAuthAccessTokenList, out *oauthapi.OAuthAccessTokenList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthAccessTokenList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapi.OAuthAccessToken, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_OAuthAccessToken_To_api_OAuthAccessToken(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_OAuthAuthorizeToken_To_api_OAuthAuthorizeToken(in *oauthapiv1beta3.OAuthAuthorizeToken, out *oauthapi.OAuthAuthorizeToken, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthAuthorizeToken))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.ExpiresIn = in.ExpiresIn
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	out.RedirectURI = in.RedirectURI
	out.State = in.State
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	return nil
}

func convert_v1beta3_OAuthAuthorizeTokenList_To_api_OAuthAuthorizeTokenList(in *oauthapiv1beta3.OAuthAuthorizeTokenList, out *oauthapi.OAuthAuthorizeTokenList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthAuthorizeTokenList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapi.OAuthAuthorizeToken, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_OAuthAuthorizeToken_To_api_OAuthAuthorizeToken(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_OAuthClient_To_api_OAuthClient(in *oauthapiv1beta3.OAuthClient, out *oauthapi.OAuthClient, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthClient))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Secret = in.Secret
	out.RespondWithChallenges = in.RespondWithChallenges
	if in.RedirectURIs != nil {
		out.RedirectURIs = make([]string, len(in.RedirectURIs))
		for i := range in.RedirectURIs {
			out.RedirectURIs[i] = in.RedirectURIs[i]
		}
	} else {
		out.RedirectURIs = nil
	}
	return nil
}

func convert_v1beta3_OAuthClientAuthorization_To_api_OAuthClientAuthorization(in *oauthapiv1beta3.OAuthClientAuthorization, out *oauthapi.OAuthClientAuthorization, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthClientAuthorization))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ClientName = in.ClientName
	out.UserName = in.UserName
	out.UserUID = in.UserUID
	if in.Scopes != nil {
		out.Scopes = make([]string, len(in.Scopes))
		for i := range in.Scopes {
			out.Scopes[i] = in.Scopes[i]
		}
	} else {
		out.Scopes = nil
	}
	return nil
}

func convert_v1beta3_OAuthClientAuthorizationList_To_api_OAuthClientAuthorizationList(in *oauthapiv1beta3.OAuthClientAuthorizationList, out *oauthapi.OAuthClientAuthorizationList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthClientAuthorizationList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapi.OAuthClientAuthorization, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_OAuthClientAuthorization_To_api_OAuthClientAuthorization(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_OAuthClientList_To_api_OAuthClientList(in *oauthapiv1beta3.OAuthClientList, out *oauthapi.OAuthClientList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*oauthapiv1beta3.OAuthClientList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]oauthapi.OAuthClient, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_OAuthClient_To_api_OAuthClient(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_Project_To_v1beta3_Project(in *projectapi.Project, out *projectapiv1beta3.Project, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapi.Project))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := convert_api_ProjectSpec_To_v1beta3_ProjectSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := convert_api_ProjectStatus_To_v1beta3_ProjectStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func convert_api_ProjectList_To_v1beta3_ProjectList(in *projectapi.ProjectList, out *projectapiv1beta3.ProjectList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapi.ProjectList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]projectapiv1beta3.Project, len(in.Items))
		for i := range in.Items {
			if err := convert_api_Project_To_v1beta3_Project(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ProjectRequest_To_v1beta3_ProjectRequest(in *projectapi.ProjectRequest, out *projectapiv1beta3.ProjectRequest, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapi.ProjectRequest))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.DisplayName = in.DisplayName
	out.Description = in.Description
	return nil
}

func convert_api_ProjectSpec_To_v1beta3_ProjectSpec(in *projectapi.ProjectSpec, out *projectapiv1beta3.ProjectSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapi.ProjectSpec))(in)
	}
	if in.Finalizers != nil {
		out.Finalizers = make([]v1beta3.FinalizerName, len(in.Finalizers))
		for i := range in.Finalizers {
			out.Finalizers[i] = v1beta3.FinalizerName(in.Finalizers[i])
		}
	} else {
		out.Finalizers = nil
	}
	return nil
}

func convert_api_ProjectStatus_To_v1beta3_ProjectStatus(in *projectapi.ProjectStatus, out *projectapiv1beta3.ProjectStatus, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapi.ProjectStatus))(in)
	}
	out.Phase = v1beta3.NamespacePhase(in.Phase)
	return nil
}

func convert_v1beta3_Project_To_api_Project(in *projectapiv1beta3.Project, out *projectapi.Project, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapiv1beta3.Project))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ProjectSpec_To_api_ProjectSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ProjectStatus_To_api_ProjectStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_ProjectList_To_api_ProjectList(in *projectapiv1beta3.ProjectList, out *projectapi.ProjectList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapiv1beta3.ProjectList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]projectapi.Project, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_Project_To_api_Project(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ProjectRequest_To_api_ProjectRequest(in *projectapiv1beta3.ProjectRequest, out *projectapi.ProjectRequest, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapiv1beta3.ProjectRequest))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.DisplayName = in.DisplayName
	out.Description = in.Description
	return nil
}

func convert_v1beta3_ProjectSpec_To_api_ProjectSpec(in *projectapiv1beta3.ProjectSpec, out *projectapi.ProjectSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapiv1beta3.ProjectSpec))(in)
	}
	if in.Finalizers != nil {
		out.Finalizers = make([]api.FinalizerName, len(in.Finalizers))
		for i := range in.Finalizers {
			out.Finalizers[i] = api.FinalizerName(in.Finalizers[i])
		}
	} else {
		out.Finalizers = nil
	}
	return nil
}

func convert_v1beta3_ProjectStatus_To_api_ProjectStatus(in *projectapiv1beta3.ProjectStatus, out *projectapi.ProjectStatus, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*projectapiv1beta3.ProjectStatus))(in)
	}
	out.Phase = api.NamespacePhase(in.Phase)
	return nil
}

func convert_api_RouteList_To_v1beta3_RouteList(in *routeapi.RouteList, out *routeapiv1beta3.RouteList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*routeapi.RouteList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]routeapiv1beta3.Route, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_RouteList_To_api_RouteList(in *routeapiv1beta3.RouteList, out *routeapi.RouteList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*routeapiv1beta3.RouteList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]routeapi.Route, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_ClusterNetwork_To_v1beta3_ClusterNetwork(in *sdnapi.ClusterNetwork, out *sdnapiv1beta3.ClusterNetwork, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapi.ClusterNetwork))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Network = in.Network
	out.HostSubnetLength = in.HostSubnetLength
	return nil
}

func convert_api_ClusterNetworkList_To_v1beta3_ClusterNetworkList(in *sdnapi.ClusterNetworkList, out *sdnapiv1beta3.ClusterNetworkList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapi.ClusterNetworkList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]sdnapiv1beta3.ClusterNetwork, len(in.Items))
		for i := range in.Items {
			if err := convert_api_ClusterNetwork_To_v1beta3_ClusterNetwork(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_HostSubnet_To_v1beta3_HostSubnet(in *sdnapi.HostSubnet, out *sdnapiv1beta3.HostSubnet, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapi.HostSubnet))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Host = in.Host
	out.HostIP = in.HostIP
	out.Subnet = in.Subnet
	return nil
}

func convert_api_HostSubnetList_To_v1beta3_HostSubnetList(in *sdnapi.HostSubnetList, out *sdnapiv1beta3.HostSubnetList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapi.HostSubnetList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]sdnapiv1beta3.HostSubnet, len(in.Items))
		for i := range in.Items {
			if err := convert_api_HostSubnet_To_v1beta3_HostSubnet(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_ClusterNetwork_To_api_ClusterNetwork(in *sdnapiv1beta3.ClusterNetwork, out *sdnapi.ClusterNetwork, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapiv1beta3.ClusterNetwork))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Network = in.Network
	out.HostSubnetLength = in.HostSubnetLength
	return nil
}

func convert_v1beta3_ClusterNetworkList_To_api_ClusterNetworkList(in *sdnapiv1beta3.ClusterNetworkList, out *sdnapi.ClusterNetworkList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapiv1beta3.ClusterNetworkList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]sdnapi.ClusterNetwork, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_ClusterNetwork_To_api_ClusterNetwork(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_HostSubnet_To_api_HostSubnet(in *sdnapiv1beta3.HostSubnet, out *sdnapi.HostSubnet, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapiv1beta3.HostSubnet))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.Host = in.Host
	out.HostIP = in.HostIP
	out.Subnet = in.Subnet
	return nil
}

func convert_v1beta3_HostSubnetList_To_api_HostSubnetList(in *sdnapiv1beta3.HostSubnetList, out *sdnapi.HostSubnetList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*sdnapiv1beta3.HostSubnetList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]sdnapi.HostSubnet, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_HostSubnet_To_api_HostSubnet(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_Parameter_To_v1beta3_Parameter(in *templateapi.Parameter, out *templateapiv1beta3.Parameter, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*templateapi.Parameter))(in)
	}
	out.Name = in.Name
	out.Description = in.Description
	out.Value = in.Value
	out.Generate = in.Generate
	out.From = in.From
	return nil
}

func convert_api_TemplateList_To_v1beta3_TemplateList(in *templateapi.TemplateList, out *templateapiv1beta3.TemplateList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*templateapi.TemplateList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]templateapiv1beta3.Template, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_Parameter_To_api_Parameter(in *templateapiv1beta3.Parameter, out *templateapi.Parameter, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*templateapiv1beta3.Parameter))(in)
	}
	out.Name = in.Name
	out.Description = in.Description
	out.Value = in.Value
	out.Generate = in.Generate
	out.From = in.From
	return nil
}

func convert_v1beta3_TemplateList_To_api_TemplateList(in *templateapiv1beta3.TemplateList, out *templateapi.TemplateList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*templateapiv1beta3.TemplateList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]templateapi.Template, len(in.Items))
		for i := range in.Items {
			if err := s.Convert(&in.Items[i], &out.Items[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_Identity_To_v1beta3_Identity(in *userapi.Identity, out *userapiv1beta3.Identity, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapi.Identity))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ProviderName = in.ProviderName
	out.ProviderUserName = in.ProviderUserName
	if err := convert_api_ObjectReference_To_v1beta3_ObjectReference(&in.User, &out.User, s); err != nil {
		return err
	}
	if in.Extra != nil {
		out.Extra = make(map[string]string)
		for key, val := range in.Extra {
			out.Extra[key] = val
		}
	} else {
		out.Extra = nil
	}
	return nil
}

func convert_api_IdentityList_To_v1beta3_IdentityList(in *userapi.IdentityList, out *userapiv1beta3.IdentityList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapi.IdentityList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]userapiv1beta3.Identity, len(in.Items))
		for i := range in.Items {
			if err := convert_api_Identity_To_v1beta3_Identity(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_api_User_To_v1beta3_User(in *userapi.User, out *userapiv1beta3.User, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapi.User))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.FullName = in.FullName
	if in.Identities != nil {
		out.Identities = make([]string, len(in.Identities))
		for i := range in.Identities {
			out.Identities[i] = in.Identities[i]
		}
	} else {
		out.Identities = nil
	}
	if in.Groups != nil {
		out.Groups = make([]string, len(in.Groups))
		for i := range in.Groups {
			out.Groups[i] = in.Groups[i]
		}
	} else {
		out.Groups = nil
	}
	return nil
}

func convert_api_UserIdentityMapping_To_v1beta3_UserIdentityMapping(in *userapi.UserIdentityMapping, out *userapiv1beta3.UserIdentityMapping, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapi.UserIdentityMapping))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectMeta_To_v1beta3_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := convert_api_ObjectReference_To_v1beta3_ObjectReference(&in.Identity, &out.Identity, s); err != nil {
		return err
	}
	if err := convert_api_ObjectReference_To_v1beta3_ObjectReference(&in.User, &out.User, s); err != nil {
		return err
	}
	return nil
}

func convert_api_UserList_To_v1beta3_UserList(in *userapi.UserList, out *userapiv1beta3.UserList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapi.UserList))(in)
	}
	if err := convert_api_TypeMeta_To_v1beta3_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_api_ListMeta_To_v1beta3_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]userapiv1beta3.User, len(in.Items))
		for i := range in.Items {
			if err := convert_api_User_To_v1beta3_User(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_Identity_To_api_Identity(in *userapiv1beta3.Identity, out *userapi.Identity, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapiv1beta3.Identity))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.ProviderName = in.ProviderName
	out.ProviderUserName = in.ProviderUserName
	if err := convert_v1beta3_ObjectReference_To_api_ObjectReference(&in.User, &out.User, s); err != nil {
		return err
	}
	if in.Extra != nil {
		out.Extra = make(map[string]string)
		for key, val := range in.Extra {
			out.Extra[key] = val
		}
	} else {
		out.Extra = nil
	}
	return nil
}

func convert_v1beta3_IdentityList_To_api_IdentityList(in *userapiv1beta3.IdentityList, out *userapi.IdentityList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapiv1beta3.IdentityList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]userapi.Identity, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_Identity_To_api_Identity(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1beta3_User_To_api_User(in *userapiv1beta3.User, out *userapi.User, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapiv1beta3.User))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	out.FullName = in.FullName
	if in.Identities != nil {
		out.Identities = make([]string, len(in.Identities))
		for i := range in.Identities {
			out.Identities[i] = in.Identities[i]
		}
	} else {
		out.Identities = nil
	}
	if in.Groups != nil {
		out.Groups = make([]string, len(in.Groups))
		for i := range in.Groups {
			out.Groups[i] = in.Groups[i]
		}
	} else {
		out.Groups = nil
	}
	return nil
}

func convert_v1beta3_UserIdentityMapping_To_api_UserIdentityMapping(in *userapiv1beta3.UserIdentityMapping, out *userapi.UserIdentityMapping, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapiv1beta3.UserIdentityMapping))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectReference_To_api_ObjectReference(&in.Identity, &out.Identity, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ObjectReference_To_api_ObjectReference(&in.User, &out.User, s); err != nil {
		return err
	}
	return nil
}

func convert_v1beta3_UserList_To_api_UserList(in *userapiv1beta3.UserList, out *userapi.UserList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*userapiv1beta3.UserList))(in)
	}
	if err := convert_v1beta3_TypeMeta_To_api_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := convert_v1beta3_ListMeta_To_api_ListMeta(&in.ListMeta, &out.ListMeta, s); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]userapi.User, len(in.Items))
		for i := range in.Items {
			if err := convert_v1beta3_User_To_api_User(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func init() {
	err := api.Scheme.AddGeneratedConversionFuncs(
		convert_api_BuildConfigList_To_v1beta3_BuildConfigList,
		convert_api_BuildList_To_v1beta3_BuildList,
		convert_api_BuildLogOptions_To_v1beta3_BuildLogOptions,
		convert_api_BuildLog_To_v1beta3_BuildLog,
		convert_api_BuildRequest_To_v1beta3_BuildRequest,
		convert_api_ClusterNetworkList_To_v1beta3_ClusterNetworkList,
		convert_api_ClusterNetwork_To_v1beta3_ClusterNetwork,
		convert_api_ClusterPolicyBindingList_To_v1beta3_ClusterPolicyBindingList,
		convert_api_ClusterPolicyList_To_v1beta3_ClusterPolicyList,
		convert_api_ClusterPolicy_To_v1beta3_ClusterPolicy,
		convert_api_ClusterRoleBindingList_To_v1beta3_ClusterRoleBindingList,
		convert_api_ClusterRoleList_To_v1beta3_ClusterRoleList,
		convert_api_ClusterRole_To_v1beta3_ClusterRole,
		convert_api_DeploymentConfigList_To_v1beta3_DeploymentConfigList,
		convert_api_DeploymentConfigRollbackSpec_To_v1beta3_DeploymentConfigRollbackSpec,
		convert_api_DeploymentConfigRollback_To_v1beta3_DeploymentConfigRollback,
		convert_api_GitSourceRevision_To_v1beta3_GitSourceRevision,
		convert_api_HostSubnetList_To_v1beta3_HostSubnetList,
		convert_api_HostSubnet_To_v1beta3_HostSubnet,
		convert_api_IdentityList_To_v1beta3_IdentityList,
		convert_api_Identity_To_v1beta3_Identity,
		convert_api_ImageList_To_v1beta3_ImageList,
		convert_api_ImageStreamList_To_v1beta3_ImageStreamList,
		convert_api_IsPersonalSubjectAccessReview_To_v1beta3_IsPersonalSubjectAccessReview,
		convert_api_ListMeta_To_v1beta3_ListMeta,
		convert_api_OAuthAccessTokenList_To_v1beta3_OAuthAccessTokenList,
		convert_api_OAuthAccessToken_To_v1beta3_OAuthAccessToken,
		convert_api_OAuthAuthorizeTokenList_To_v1beta3_OAuthAuthorizeTokenList,
		convert_api_OAuthAuthorizeToken_To_v1beta3_OAuthAuthorizeToken,
		convert_api_OAuthClientAuthorizationList_To_v1beta3_OAuthClientAuthorizationList,
		convert_api_OAuthClientAuthorization_To_v1beta3_OAuthClientAuthorization,
		convert_api_OAuthClientList_To_v1beta3_OAuthClientList,
		convert_api_OAuthClient_To_v1beta3_OAuthClient,
		convert_api_ObjectMeta_To_v1beta3_ObjectMeta,
		convert_api_ObjectReference_To_v1beta3_ObjectReference,
		convert_api_Parameter_To_v1beta3_Parameter,
		convert_api_PolicyBindingList_To_v1beta3_PolicyBindingList,
		convert_api_PolicyList_To_v1beta3_PolicyList,
		convert_api_ProjectList_To_v1beta3_ProjectList,
		convert_api_ProjectRequest_To_v1beta3_ProjectRequest,
		convert_api_ProjectSpec_To_v1beta3_ProjectSpec,
		convert_api_ProjectStatus_To_v1beta3_ProjectStatus,
		convert_api_Project_To_v1beta3_Project,
		convert_api_ResourceAccessReview_To_v1beta3_ResourceAccessReview,
		convert_api_RoleBindingList_To_v1beta3_RoleBindingList,
		convert_api_RoleList_To_v1beta3_RoleList,
		convert_api_Role_To_v1beta3_Role,
		convert_api_RouteList_To_v1beta3_RouteList,
		convert_api_SourceControlUser_To_v1beta3_SourceControlUser,
		convert_api_SourceRevision_To_v1beta3_SourceRevision,
		convert_api_SubjectAccessReviewResponse_To_v1beta3_SubjectAccessReviewResponse,
		convert_api_TemplateList_To_v1beta3_TemplateList,
		convert_api_TypeMeta_To_v1beta3_TypeMeta,
		convert_api_UserIdentityMapping_To_v1beta3_UserIdentityMapping,
		convert_api_UserList_To_v1beta3_UserList,
		convert_api_User_To_v1beta3_User,
		convert_v1beta3_BuildConfigList_To_api_BuildConfigList,
		convert_v1beta3_BuildList_To_api_BuildList,
		convert_v1beta3_BuildLogOptions_To_api_BuildLogOptions,
		convert_v1beta3_BuildLog_To_api_BuildLog,
		convert_v1beta3_BuildRequest_To_api_BuildRequest,
		convert_v1beta3_ClusterNetworkList_To_api_ClusterNetworkList,
		convert_v1beta3_ClusterNetwork_To_api_ClusterNetwork,
		convert_v1beta3_ClusterPolicyBindingList_To_api_ClusterPolicyBindingList,
		convert_v1beta3_ClusterPolicyList_To_api_ClusterPolicyList,
		convert_v1beta3_ClusterPolicy_To_api_ClusterPolicy,
		convert_v1beta3_ClusterRoleBindingList_To_api_ClusterRoleBindingList,
		convert_v1beta3_ClusterRoleList_To_api_ClusterRoleList,
		convert_v1beta3_ClusterRole_To_api_ClusterRole,
		convert_v1beta3_DeploymentConfigList_To_api_DeploymentConfigList,
		convert_v1beta3_DeploymentConfigRollbackSpec_To_api_DeploymentConfigRollbackSpec,
		convert_v1beta3_DeploymentConfigRollback_To_api_DeploymentConfigRollback,
		convert_v1beta3_GitSourceRevision_To_api_GitSourceRevision,
		convert_v1beta3_HostSubnetList_To_api_HostSubnetList,
		convert_v1beta3_HostSubnet_To_api_HostSubnet,
		convert_v1beta3_IdentityList_To_api_IdentityList,
		convert_v1beta3_Identity_To_api_Identity,
		convert_v1beta3_ImageList_To_api_ImageList,
		convert_v1beta3_ImageStreamList_To_api_ImageStreamList,
		convert_v1beta3_IsPersonalSubjectAccessReview_To_api_IsPersonalSubjectAccessReview,
		convert_v1beta3_ListMeta_To_api_ListMeta,
		convert_v1beta3_OAuthAccessTokenList_To_api_OAuthAccessTokenList,
		convert_v1beta3_OAuthAccessToken_To_api_OAuthAccessToken,
		convert_v1beta3_OAuthAuthorizeTokenList_To_api_OAuthAuthorizeTokenList,
		convert_v1beta3_OAuthAuthorizeToken_To_api_OAuthAuthorizeToken,
		convert_v1beta3_OAuthClientAuthorizationList_To_api_OAuthClientAuthorizationList,
		convert_v1beta3_OAuthClientAuthorization_To_api_OAuthClientAuthorization,
		convert_v1beta3_OAuthClientList_To_api_OAuthClientList,
		convert_v1beta3_OAuthClient_To_api_OAuthClient,
		convert_v1beta3_ObjectMeta_To_api_ObjectMeta,
		convert_v1beta3_ObjectReference_To_api_ObjectReference,
		convert_v1beta3_Parameter_To_api_Parameter,
		convert_v1beta3_PolicyBindingList_To_api_PolicyBindingList,
		convert_v1beta3_PolicyList_To_api_PolicyList,
		convert_v1beta3_ProjectList_To_api_ProjectList,
		convert_v1beta3_ProjectRequest_To_api_ProjectRequest,
		convert_v1beta3_ProjectSpec_To_api_ProjectSpec,
		convert_v1beta3_ProjectStatus_To_api_ProjectStatus,
		convert_v1beta3_Project_To_api_Project,
		convert_v1beta3_ResourceAccessReview_To_api_ResourceAccessReview,
		convert_v1beta3_RoleBindingList_To_api_RoleBindingList,
		convert_v1beta3_RoleList_To_api_RoleList,
		convert_v1beta3_Role_To_api_Role,
		convert_v1beta3_RouteList_To_api_RouteList,
		convert_v1beta3_SourceControlUser_To_api_SourceControlUser,
		convert_v1beta3_SourceRevision_To_api_SourceRevision,
		convert_v1beta3_SubjectAccessReviewResponse_To_api_SubjectAccessReviewResponse,
		convert_v1beta3_TemplateList_To_api_TemplateList,
		convert_v1beta3_TypeMeta_To_api_TypeMeta,
		convert_v1beta3_UserIdentityMapping_To_api_UserIdentityMapping,
		convert_v1beta3_UserList_To_api_UserList,
		convert_v1beta3_User_To_api_User,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}

// AUTO-GENERATED FUNCTIONS END HERE
