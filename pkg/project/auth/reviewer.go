package auth

import (
	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	"github.com/projectatomic/appinfra-next/pkg/client"
)

// Review is a list of users and groups that can access a resource
type Review interface {
	Users() []string
	Groups() []string
}

type review struct {
	response *authorizationapi.ResourceAccessReviewResponse
}

// Users returns the users that can access a resource
func (r *review) Users() []string {
	return r.response.Users.List()
}

// Groups returns the groups that can access a resource
func (r *review) Groups() []string {
	return r.response.Groups.List()
}

// Reviewer performs access reviews for a project by name
type Reviewer interface {
	Review(name string) (Review, error)
}

// reviewer performs access reviews for a project by name
type reviewer struct {
	resourceAccessReviewsNamespacer client.ResourceAccessReviewsNamespacer
}

// NewReviewer knows how to make access control reviews for a resource by name
func NewReviewer(resourceAccessReviewsNamespacer client.ResourceAccessReviewsNamespacer) Reviewer {
	return &reviewer{
		resourceAccessReviewsNamespacer: resourceAccessReviewsNamespacer,
	}
}

// Review performs a resource access review for the given resource by name
func (r *reviewer) Review(name string) (Review, error) {
	resourceAccessReview := &authorizationapi.ResourceAccessReview{
		Verb:         "get",
		Resource:     "namespaces",
		ResourceName: name,
	}

	response, err := r.resourceAccessReviewsNamespacer.ResourceAccessReviews(name).Create(resourceAccessReview)

	if err != nil {
		return nil, err
	}
	review := &review{
		response: response,
	}
	return review, nil
}
