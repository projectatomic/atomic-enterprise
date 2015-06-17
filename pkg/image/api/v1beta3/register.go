package v1beta3

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	_ "github.com/projectatomic/appinfra-next/pkg/image/api/docker10"
	_ "github.com/projectatomic/appinfra-next/pkg/image/api/dockerpre012"
)

func init() {
	api.Scheme.AddKnownTypes("v1beta3",
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
