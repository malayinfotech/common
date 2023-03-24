// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"common/sync2"
	"common/testcontext"
	"drpc"
	"drpc/drpcmigrate"
	"drpc/drpcserver"
)

func TestDialerUnencrypted(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	d := NewDefaultPooledDialer(nil)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ctx.Check(lis.Close)

	conn, err := d.DialAddressUnencrypted(ctx, lis.Addr().String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

type goodHandler struct{}

func (goodHandler) HandleRPC(stream drpc.Stream, rpc string) error { return nil }

func TestDialHostnameVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	certificatePEM, privateKeyPEM := createTestingCertificate(t, "localhost")

	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	require.NoError(t, err)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	// start a server with the certificate
	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	serverAddr := tcpListener.Addr().String()

	listenMux := drpcmigrate.NewListenMux(tcpListener, len(drpcmigrate.DRPCHeader))

	listenCtx, listenCancel := context.WithCancel(ctx)
	defer listenCancel()
	ctx.Go(func() error {
		return listenMux.Run(listenCtx)
	})

	drpcListener := tls.NewListener(listenMux.Route(drpcmigrate.DRPCHeader), serverTLSConfig)
	defer ctx.Check(drpcListener.Close)

	acceptConnectionSuccess := func() (err error) {
		conn, err := drpcListener.Accept()
		if err != nil {
			return errs.Wrap(err)
		}
		defer func() { _ = conn.Close() }()

		_ = drpcserver.New(new(goodHandler)).ServeOne(ctx, conn)
		return nil
	}

	acceptConnectionFailure := func() error {
		conn, err := drpcListener.Accept()
		if err != nil {
			return errs.Wrap(err)
		}
		defer func() { _ = conn.Close() }()

		buffer := make([]byte, 256)
		_, err = conn.Read(buffer)
		if err == nil {
			return errs.New("expected connection failure, but there is no error")
		}
		return nil
	}

	useConn := func(ctx context.Context, conn drpc.Conn) error {
		stream, err := conn.NewStream(ctx, "test-rpc", nil)
		if err != nil {
			return errs.Wrap(err)
		}
		return errs.Wrap(stream.Close())
	}

	// create client
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certificatePEM)

	// happy scenario 1 get hostname from address
	requireNoErrors(t, sync2.Concurrently(
		acceptConnectionSuccess,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs: certPool,
			}

			// use a domain name to ensure we can get hostname from the address
			localAddr := strings.ReplaceAll(serverAddr, "127.0.0.1", "localhost")
			conn, err := dialer.DialAddressHostnameVerification(ctx, localAddr)
			if err != nil {
				return errs.Wrap(err)
			}
			defer func() { _ = conn.Close() }()

			if err := useConn(ctx, conn); err != nil {
				return errs.Wrap(err)
			}

			return errs.Wrap(conn.Close())
		},
	))

	// happy scenario 2
	requireNoErrors(t, sync2.Concurrently(
		acceptConnectionSuccess,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs:    certPool,
				ServerName: "localhost",
			}
			// Can't verify IPv6 during ci because of docker default
			// connection, err = dialer.DialAddressHostnameVerification(ctx, "[::1]:22111", clientTLSConfig)
			conn, err := dialer.DialAddressHostnameVerification(ctx, serverAddr)
			if err != nil {
				return errs.Wrap(err)
			}
			defer func() { _ = conn.Close() }()

			if err := useConn(ctx, conn); err != nil {
				return errs.Wrap(err)
			}

			return errs.Wrap(conn.Close())
		},
	))

	// failure scenario invalid certificate
	requireNoErrors(t, sync2.Concurrently(
		acceptConnectionFailure,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs:    certPool,
				ServerName: "storx.test",
			}
			conn, err := dialer.DialAddressHostnameVerification(ctx, serverAddr)
			if err != nil {
				return errs.Wrap(err)
			}
			defer func() { _ = conn.Close() }()

			err = useConn(ctx, conn)
			if err == nil {
				return errs.New("expected an error")
			}
			if !strings.Contains(err.Error(), "certificate is valid for localhost, not storx.test") {
				return errs.New("expected an error, got: %w", err)
			}

			return nil
		},
	))

	// test invalid hostname
	dialer := NewDefaultDialer(nil)
	_, err = dialer.DialAddressHostnameVerification(ctx, "storx.test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing port in address")
}

func createTestingCertificate(t *testing.T, hostname string) (certificatePEM []byte, privateKeyPEM []byte) {
	notAfter := time.Now().Add(1 * time.Minute)

	// first create a server certificate
	template := x509.Certificate{
		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames:              []string{hostname},
		SerialNumber:          big.NewInt(1337),
		BasicConstraintsValid: false,
		IsCA:                  true,
		NotAfter:              notAfter,
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	certificateDERBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	require.NoError(t, err)

	certificatePEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDERBytes})

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})

	return certificatePEM, privateKeyPEM
}

func requireNoErrors(t *testing.T, errs []error) {
	t.Helper()
	if len(errs) > 0 {
		for _, err := range errs {
			assert.NoError(t, err)
		}
		t.Fatal()
	}
}
