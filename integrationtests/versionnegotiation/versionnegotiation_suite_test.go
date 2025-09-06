package versionnegotiation

import (
	"context"
	"crypto/x509"
	"flag"
	"testing"

	quic "github.com/enetx/uquic"
	tls "github.com/refraction-networking/utls"

	"github.com/enetx/uquic/integrationtests/tools"
	"github.com/enetx/uquic/logging"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	enableQlog      bool
	tlsConfig       *tls.Config
	tlsClientConfig *tls.Config
)

func init() {
	flag.BoolVar(&enableQlog, "qlog", false, "enable qlog")

	ca, caPrivateKey, err := tools.GenerateCA()
	if err != nil {
		panic(err)
	}
	leafCert, leafPrivateKey, err := tools.GenerateLeafCert(ca, caPrivateKey)
	if err != nil {
		panic(err)
	}
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leafCert.Raw},
			PrivateKey:  leafPrivateKey,
		}},
		NextProtos: []string{tools.ALPN},
	}

	root := x509.NewCertPool()
	root.AddCert(ca)
	tlsClientConfig = &tls.Config{
		ServerName: "localhost",
		RootCAs:    root,
		NextProtos: []string{tools.ALPN},
	}
}

func getTLSConfig() *tls.Config       { return tlsConfig }
func getTLSClientConfig() *tls.Config { return tlsClientConfig }

func TestQuicVersionNegotiation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Negotiation Suite")
}

func maybeAddQLOGTracer(c *quic.Config) *quic.Config {
	if c == nil {
		c = &quic.Config{}
	}
	if !enableQlog {
		return c
	}
	qlogger := tools.NewQlogConnectionTracer(GinkgoWriter)
	if c.Tracer == nil {
		c.Tracer = qlogger
	} else if qlogger != nil {
		origTracer := c.Tracer
		c.Tracer = func(ctx context.Context, p logging.Perspective, connID quic.ConnectionID) *logging.ConnectionTracer {
			return logging.NewMultiplexedConnectionTracer(
				qlogger(ctx, p, connID),
				origTracer(ctx, p, connID),
			)
		}
	}
	return c
}
