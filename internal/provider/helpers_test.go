package provider

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

func startVault(t *testing.T) {
	t.Helper()

	configPath, err := filepath.Abs("../../vault.hcl")
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
			clusterHost, clusterPort, err := net.SplitHostPort(match[1])
			if err != nil {
				t.Fatal(err)
			}

			port, err := strconv.Atoi(clusterPort)
			if err != nil {
				t.Fatal(err)
			}

			t.Setenv("VAULT_ADDR", fmt.Sprintf("http://%s:%d", clusterHost, port-1))
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
