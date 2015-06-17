package x509request

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	kerrors "github.com/GoogleCloudPlatform/kubernetes/pkg/util/errors"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator"
)

// UserConversion defines an interface for extracting user info from a client certificate chain
type UserConversion interface {
	User(chain []*x509.Certificate) (user.Info, bool, error)
}

// UserConversionFunc is a function that implements the UserConversion interface.
type UserConversionFunc func(chain []*x509.Certificate) (user.Info, bool, error)

// User implements x509.UserConversion
func (f UserConversionFunc) User(chain []*x509.Certificate) (user.Info, bool, error) {
	return f(chain)
}

// Authenticator implements request.Authenticator by extracting user info from verified client certificates
type Authenticator struct {
	opts x509.VerifyOptions
	user UserConversion
}

// New returns a request.Authenticator that verifies client certificates using the provided
// VerifyOptions, and converts valid certificate chains into user.Info using the provided UserConversion
func New(opts x509.VerifyOptions, user UserConversion) *Authenticator {
	return &Authenticator{opts, user}
}

// AuthenticateRequest authenticates the request using presented client certificates
func (a *Authenticator) AuthenticateRequest(req *http.Request) (user.Info, bool, error) {
	if req.TLS == nil {
		return nil, false, nil
	}

	var errlist []error
	for _, cert := range req.TLS.PeerCertificates {
		chains, err := cert.Verify(a.opts)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}

		for _, chain := range chains {
			user, ok, err := a.user.User(chain)
			if err != nil {
				errlist = append(errlist, err)
				continue
			}

			if ok {
				return user, ok, err
			}
		}
	}
	return nil, false, kerrors.NewAggregate(errlist)
}

// Verifier implements request.Authenticator by verifying a client cert on the request, then delegating to the wrapped auth
type Verifier struct {
	opts x509.VerifyOptions
	auth authenticator.Request
}

func NewVerifier(opts x509.VerifyOptions, auth authenticator.Request) authenticator.Request {
	return &Verifier{opts, auth}
}

// AuthenticateRequest verifies the presented client certificates, then delegates to the wrapped auth
func (a *Verifier) AuthenticateRequest(req *http.Request) (user.Info, bool, error) {
	if req.TLS == nil {
		return nil, false, nil
	}

	var errlist []error
	for _, cert := range req.TLS.PeerCertificates {
		_, err := cert.Verify(a.opts)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}
		return a.auth.AuthenticateRequest(req)
	}
	return nil, false, kerrors.NewAggregate(errlist)
}

// DefaultVerifyOptions returns VerifyOptions that use the system root certificates, current time,
// and requires certificates to be valid for client auth (x509.ExtKeyUsageClientAuth)
func DefaultVerifyOptions() x509.VerifyOptions {
	return x509.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
}

// SubjectToUserConversion calls SubjectToUser on the subject of the first certificate in the chain.
// If the resulting user has no name, it returns nil, false, nil
var SubjectToUserConversion = UserConversionFunc(func(chain []*x509.Certificate) (user.Info, bool, error) {
	user := SubjectToUser(chain[0].Subject)
	if len(user.GetName()) == 0 {
		return nil, false, nil
	}
	return user, true, nil
})

// CommonNameUserConversion builds user info from a certificate chain using the subject's CommonName
var CommonNameUserConversion = UserConversionFunc(func(chain []*x509.Certificate) (user.Info, bool, error) {
	if len(chain[0].Subject.CommonName) == 0 {
		return nil, false, nil
	}
	return &user.DefaultInfo{Name: chain[0].Subject.CommonName}, true, nil
})

// DNSNameUserConversion builds user info from a certificate chain using the first DNSName on the certificate
var DNSNameUserConversion = UserConversionFunc(func(chain []*x509.Certificate) (user.Info, bool, error) {
	if len(chain[0].DNSNames) == 0 {
		return nil, false, nil
	}
	return &user.DefaultInfo{Name: chain[0].DNSNames[0]}, true, nil
})

// EmailAddressUserConversion builds user info from a certificate chain using the first EmailAddress on the certificate
var EmailAddressUserConversion = UserConversionFunc(func(chain []*x509.Certificate) (user.Info, bool, error) {
	if len(chain[0].EmailAddresses) == 0 {
		return nil, false, nil
	}
	return &user.DefaultInfo{Name: chain[0].EmailAddresses[0]}, true, nil
})

func UserToSubject(u user.Info) pkix.Name {
	return pkix.Name{
		CommonName:   u.GetName(),
		SerialNumber: u.GetUID(),
		Organization: u.GetGroups(),
	}
}
func SubjectToUser(subject pkix.Name) user.Info {
	return &user.DefaultInfo{
		Name:   subject.CommonName,
		UID:    subject.SerialNumber,
		Groups: subject.Organization,
	}
}
