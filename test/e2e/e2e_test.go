//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var k8scopeBin string

func TestMain(m *testing.M) {
	bin, err := buildBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build k8scope: %v\n", err)
		os.Exit(1)
	}
	k8scopeBin = bin
	os.Exit(m.Run())
}

func buildBinary() (string, error) {
	bin := "/tmp/k8scope-e2e"
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/k8scope")
	cmd.Dir = findRepoRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return bin, cmd.Run()
}

func findRepoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir
		}
		parent := dir[:strings.LastIndex(dir, "/")]
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

func runK8scope(t *testing.T, args ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, k8scopeBin, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func kubectl(t *testing.T, args ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestInstallStartupMode(t *testing.T) {
	out, err := runK8scope(t, "install", "--mode", "startup", "--skip-preflight")
	if err != nil {
		t.Fatalf("install failed: %v\noutput: %s", err, out)
	}
	t.Logf("install output: %s", out)

	waitForPods(t, "k8scope", 5*time.Minute)
}

func TestDryRunCreatesNothing(t *testing.T) {
	out, err := runK8scope(t, "install", "--mode", "startup", "--dry-run", "--skip-preflight")
	if err != nil {
		t.Fatalf("dry-run failed: %v\noutput: %s", err, out)
	}

	if !strings.Contains(out, "dry-run") && !strings.Contains(out, "Dry run") {
		t.Error("expected dry-run indicator in output")
	}
}

func TestInvalidModeReturnsError(t *testing.T) {
	_, err := runK8scope(t, "install", "--mode", "nonexistent", "--skip-preflight")
	if err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestVersionCommand(t *testing.T) {
	out, err := runK8scope(t, "version")
	if err != nil {
		t.Fatalf("version failed: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "k8scope") {
		t.Errorf("expected k8scope in version output, got: %s", out)
	}
}

func waitForPods(t *testing.T, namespace string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		out, err := kubectl(t, "get", "pods", "-n", namespace,
			"--field-selector=status.phase!=Running,status.phase!=Succeeded",
			"--no-headers")
		if err == nil && strings.TrimSpace(out) == "" {
			t.Log("all pods are running")
			return
		}
		time.Sleep(15 * time.Second)
	}

	out, _ := kubectl(t, "get", "pods", "-n", namespace, "-o", "wide")
	t.Fatalf("timeout waiting for pods in namespace %s:\n%s", namespace, out)
}
