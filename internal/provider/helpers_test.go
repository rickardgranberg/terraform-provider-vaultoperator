package provider

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"text/template"
	"time"
)

func startVault(t *testing.T, enableTLS bool) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "vaultoperator-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "key")
	certPath := filepath.Join(tempDir, "cert")
	configPath := filepath.Join(tempDir, "vault.hcl")
	disableTLS := "1"
	protocol := "http"

	if enableTLS {
		disableTLS = "0"
		protocol = "https"

		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}

		keyFile, err := os.Create(keyPath)
		if err != nil {
			t.Fatal(err)
		}
		defer keyFile.Close()

		err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
		if err != nil {
			t.Fatal(err)
		}

		certTemplate := x509.Certificate{
			SerialNumber:          big.NewInt(1),
			Subject:               pkix.Name{CommonName: "localhost"},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(1, 0, 0),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
			IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &key.PublicKey, key)
		if err != nil {
			t.Fatal(err)
		}

		certFile, err := os.Create(certPath)
		if err != nil {
			t.Fatal(err)
		}
		defer certFile.Close()

		err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		if err != nil {
			t.Fatal(err)
		}
	}

	configTemplate, err := template.ParseFiles("../../vault.hcl")
	if err != nil {
		t.Fatal(err)
	}

	config := struct {
		CertFile   string
		KeyFile    string
		DisableTLS string
	}{
		CertFile:   certPath,
		KeyFile:    keyPath,
		DisableTLS: disableTLS,
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer configFile.Close()

	err = configTemplate.Execute(configFile, config)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("vault", "server", "-config", configPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	tee := io.TeeReader(stdout, &buf)

	scanner := bufio.NewScanner(tee)

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	clusterAddress := regexp.MustCompile("cluster address: \"(.*?)\"")
	vaultStarted := regexp.MustCompile("Vault server started!")

	for scanner.Scan() {
		if match := clusterAddress.FindStringSubmatch(scanner.Text()); match != nil {
			_, clusterPort, err := net.SplitHostPort(match[1])
			if err != nil {
				t.Fatal(err)
			}

			port, err := strconv.Atoi(clusterPort)
			if err != nil {
				t.Fatal(err)
			}

			t.Setenv("VAULT_ADDR", fmt.Sprintf("%s://localhost:%d", protocol, port-1))
		}

		if vaultStarted.MatchString(scanner.Text()) {
			t.Cleanup(stopVault(t, cmd))
			return
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	output, err := io.ReadAll(&buf)

	if err != nil {
		t.Error(err)
	} else {
		t.Error(string(output))
	}

	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	t.Error("Unable to start Vault server")
}

func stopVault(t *testing.T, cmd *exec.Cmd) func() {
	return func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Error(err)
		}

		if err := cmd.Wait(); err != nil {
			exitErr, ok := err.(*exec.ExitError)

			if !ok || exitErr.String() != "signal: killed" {
				t.Error(err)
			}
		}
	}
}
