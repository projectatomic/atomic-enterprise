package v1

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	_ "github.com/projectatomic/atomic-enterprise/pkg/image/api/docker10"
	_ "github.com/projectatomic/atomic-enterprise/pkg/image/api/dockerpre012"
)

func init() {
	api.Scheme.AddKnownTypes("v1",
		&Image{},
		&ImageList{},
		&ImageStream{},
		&ImageStreamList{},
		&ImageStreamMapping{},
		&ImageStreamTag{},
		&ImageStreamImage{},
	)
}

func (*Image) IsAnAPIObject()              {}
func (*ImageList) IsAnAPIObject()          {}
func (*ImageStream) IsAnAPIObject()        {}
func (*ImageStreamList) IsAnAPIObject()    {}
func (*ImageStreamMapping) IsAnAPIObject() {}
func (*ImageStreamTag) IsAnAPIObject()     {}
func (*ImageStreamImage) IsAnAPIObject()   {}
