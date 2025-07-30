package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

var certURL string

// CertCmd checks TLS certificate info for a given URL
var CertCmd = &cobra.Command{
	Use:   "cert --url https://example.com",
	Short: "Check TLS certificate info for an HTTPS URL",
	Run: func(cmd *cobra.Command, args []string) {
		if certURL == "" {
			fmt.Println("Please provide a URL using --url")
			return
		}
		checkTLS(certURL)
	},
}

func init() {
	CertCmd.Flags().StringVar(&certURL, "url", "", "HTTPS URL to check certificate for")
}

func checkTLS(rawURL string) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" {
		fmt.Printf("Invalid HTTPS URL: %s\n", rawURL)
		return
	}

	host := parsed.Host
	if parsed.Port() == "" {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		fmt.Printf("TLS connection error: %v\n", err)
		return
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		fmt.Println("No certificates found")
		return
	}

	cert := certs[0]
	now := time.Now()

	fmt.Printf("Certificate for %s\n", rawURL)
	fmt.Printf("Subject:        %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer:         %s\n", cert.Issuer.CommonName)
	fmt.Printf("Valid from:     %s\n", cert.NotBefore.Format(time.RFC1123))
	fmt.Printf("Valid until:    %s\n", cert.NotAfter.Format(time.RFC1123))
	fmt.Printf("Days left:      %d\n", int(cert.NotAfter.Sub(now).Hours()/24))

	roots, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("Failed to load system cert pool")
		return
	}

	opts := x509.VerifyOptions{
		DNSName:       parsed.Hostname(),
		Intermediates: x509.NewCertPool(),
		Roots:         roots,
	}

	for _, ic := range certs[1:] {
		opts.Intermediates.AddCert(ic)
	}

	if _, err := cert.Verify(opts); err != nil {
		fmt.Printf("Trusted:        No (%v)\n", err)
	} else {
		fmt.Println("Trusted:        Yes")
	}
}
