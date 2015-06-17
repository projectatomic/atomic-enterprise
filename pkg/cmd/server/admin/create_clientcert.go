package admin

import (
	"errors"
	"io"

	"github.com/golang/glog"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/projectatomic/appinfra-next/pkg/cmd/server/crypto"
)

type CreateClientCertOptions struct {
	GetSignerCertOptions *GetSignerCertOptions

	CertFile string
	KeyFile  string

	User   string
	Groups util.StringList

	Overwrite bool
	Output    io.Writer
}

func (o CreateClientCertOptions) Validate(args []string) error {
	if len(args) != 0 {
		return errors.New("no arguments are supported")
	}
	if len(o.CertFile) == 0 {
		return errors.New("cert must be provided")
	}
	if len(o.KeyFile) == 0 {
		return errors.New("key must be provided")
	}
	if len(o.User) == 0 {
		return errors.New("user must be provided")
	}

	if o.GetSignerCertOptions == nil {
		return errors.New("signer options are required")
	}
	if err := o.GetSignerCertOptions.Validate(); err != nil {
		return err
	}

	return nil
}

func (o CreateClientCertOptions) CreateClientCert() (*crypto.TLSCertificateConfig, error) {
	glog.V(4).Infof("Creating a client cert with: %#v and %#v", o, o.GetSignerCertOptions)

	signerCert, err := o.GetSignerCertOptions.GetSignerCert()
	if err != nil {
		return nil, err
	}

	var cert *crypto.TLSCertificateConfig
	written := true
	userInfo := &user.DefaultInfo{Name: o.User, Groups: o.Groups}
	if o.Overwrite {
		cert, err = signerCert.MakeClientCertificate(o.CertFile, o.KeyFile, userInfo)
	} else {
		cert, written, err = signerCert.EnsureClientCertificate(o.CertFile, o.KeyFile, userInfo)
	}
	if written {
		glog.V(3).Infof("Generated new client cert as %s and key as %s\n", o.CertFile, o.KeyFile)
	} else {
		glog.V(3).Infof("Keeping existing client cert at %s and key at %s\n", o.CertFile, o.KeyFile)
	}
	return cert, err
}
