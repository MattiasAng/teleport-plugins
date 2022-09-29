package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/trace"
	"time"
)

// CheckIfExpired returns true if there's at least 1 non-expired credential in the creds list. It also returns
// an error aggregate with an individual error for each invalid credential
func CheckIfExpired(credentials []client.Credentials) (bool, error) {
	var validCredentials []client.Credentials
	var errors []error

	for _, credential := range credentials {
		tlsConfig, err := credential.TLSConfig()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		// If tlsConfig is nil, it means this is a credential for an insecure client, we let it pass
		if tlsConfig == nil {
			continue
		}

		isValid := true
		// client.Credentials.TLSConfig() does not populate tlsConfig.Certificate:
		// it only sets tlsConfig.GetClientCertificate.
		// We have to invoke the function to retrieve the certificate chain.
		certificateChain, _ := tlsConfig.GetClientCertificate(&tls.CertificateRequestInfo{})
		if len(certificateChain.Certificate) == 0 {
			isValid = false
		}

		// We consider a chain valid if all its certs are not expired
		for _, certificate := range certificateChain.Certificate {
			parsedCert, err := x509.ParseCertificate(certificate)
			if err != nil {
				errors = append(errors, trace.WrapWithMessage(err, "failed to parse certificate while checking credential validity"))
				isValid = false
				break
			}

			if time.Now().After(parsedCert.NotAfter) {
				isValid = false
				errors = append(
					errors,
					trace.CompareFailed(
						"expired credential found: the certificate for '%s', issued by '%s' is not valid after '%s'",
						parsedCert.Subject.CommonName,
						parsedCert.Issuer.CommonName,
						parsedCert.NotAfter,
					),
				)
			}
		}
		if isValid {
			validCredentials = append(validCredentials, credential)
		}
	}

	return len(validCredentials) != 0, trace.NewAggregate(errors...)
}
