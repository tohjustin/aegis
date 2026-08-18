//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/tohjustin/aegis/pkg/badge"
)

var (
	binPath string
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "aegis-e2e-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "temp dir: %v\n", err)
		os.Exit(1)
	}
	binPath = filepath.Join(dir, "aegis")
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/aegis")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

type testServer struct {
	appURL   string
	adminURL string
	cmd      *exec.Cmd
}

func startServer(t *testing.T, extraArgs ...string) *testServer {
	t.Helper()

	appAddr := freeAddr(t)
	adminAddr := freeAddr(t)
	args := append([]string{
		"-app-addr", appAddr,
		"-admin-addr", adminAddr,
		"-github-access-token", "e2e-token",
		"-log-format", "text",
		"-log-level", "error",
	}, extraArgs...)

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	srv := &testServer{
		appURL:   "http://" + appAddr,
		adminURL: "http://" + adminAddr,
		cmd:      cmd,
	}
	t.Cleanup(func() { srv.stop(t, syscall.SIGTERM) })
	waitReady(t, srv.adminURL)
	return srv
}

func (s *testServer) stop(t *testing.T, sig syscall.Signal) {
	t.Helper()
	if s.cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-s.cmd.Process.Pid, sig)
	done := make(chan error, 1)
	go func() { done <- s.cmd.Wait() }()
	select {
	case <-time.After(10 * time.Second):
		_ = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
		t.Errorf("server did not exit after %s", sig)
	case err := <-done:
		if sig == syscall.SIGTERM && err != nil {
			t.Errorf("SIGTERM exit: %v", err)
		}
	}
}

func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

func waitReady(t *testing.T, adminURL string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, adminURL+"/readyz", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		if ctx.Err() != nil {
			t.Fatalf("readyz not ready: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func get(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func createBadge(t *testing.T, params *badge.Params) string {
	t.Helper()
	generated, err := badge.Create(params)
	if err != nil {
		t.Fatal(err)
	}
	return generated
}

func TestAdminEndpoints(t *testing.T) {
	srv := startServer(t)

	live := get(t, srv.adminURL+"/livez")
	if live.StatusCode != http.StatusOK {
		t.Errorf("livez status = %d", live.StatusCode)
	}
	if ct := live.Header.Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Errorf("livez content-type = %q", ct)
	}

	ready := get(t, srv.adminURL+"/readyz")
	if ready.StatusCode != http.StatusOK {
		t.Errorf("readyz status = %d", ready.StatusCode)
	}

	ver := get(t, srv.adminURL+"/version")
	if ver.StatusCode != http.StatusOK {
		t.Errorf("version status = %d", ver.StatusCode)
	}
	var info map[string]string
	if err := json.NewDecoder(ver.Body).Decode(&info); err != nil {
		t.Fatalf("version json: %v", err)
	}
	if info["go_version"] == "" {
		t.Error("version missing go_version")
	}

	_ = get(t, srv.appURL+"/static?subject=e2e&status=ok")

	metrics := get(t, srv.adminURL+"/metrics")
	if metrics.StatusCode != http.StatusOK {
		t.Errorf("metrics status = %d", metrics.StatusCode)
	}
	body := readBody(t, metrics)
	if !strings.Contains(body, "aegis_http_requests_total") {
		t.Error("metrics missing aegis_http_requests_total")
	}
	if !strings.Contains(body, "aegis_build_info") {
		t.Error("metrics missing aegis_build_info")
	}
	if !strings.Contains(body, "go_goroutines") {
		t.Error("metrics missing go_goroutines")
	}
	if runtime.GOOS == "linux" && !strings.Contains(body, "process_cpu_seconds_total") {
		t.Error("metrics missing process_cpu_seconds_total")
	}
}

func TestStaticBadge(t *testing.T) {
	srv := startServer(t)

	tests := []struct {
		name string
		path string
		want *badge.Params
	}{
		{
			name: "hex color",
			path: "/static?subject=testSubject&status=testStatus&color=ff0000",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus", Color: "ff0000"},
		},
		{
			name: "css color",
			path: "/static?subject=testSubject&status=testStatus&color=red",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus", Color: "red"},
		},
		{
			name: "default color",
			path: "/static?subject=testSubject&status=testStatus",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus", Color: badge.DefaultColor},
		},
		{
			name: "bad color",
			path: "/static?subject=testSubject&status=testStatus&color=badColor",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus", Color: badge.DefaultColor},
		},
		{
			name: "icon",
			path: "/static?subject=testSubject&status=testStatus&icon=brands/docker",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus", Icon: "brands/docker"},
		},
		{
			name: "bad icon",
			path: "/static?subject=testSubject&status=testStatus&icon=badIcon",
			want: &badge.Params{Subject: "testSubject", Status: "testStatus"},
		},
		{
			name: "style",
			path: "/static?subject=testSubject&status=testStatus&style=semaphoreci",
			want: &badge.Params{Style: badge.SemaphoreCIStyle, Subject: "testSubject", Status: "testStatus"},
		},
		{
			name: "bad style",
			path: "/static?subject=testSubject&status=testStatus&style=badStyle",
			want: &badge.Params{Style: badge.DefaultStyle, Subject: "testSubject", Status: "testStatus"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := get(t, srv.appURL+tt.path)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", resp.StatusCode)
			}
			if got := resp.Header.Get("Cache-Control"); got != "public, max-age=3600, s-maxage=3600" {
				t.Errorf("Cache-Control = %q", got)
			}
			if got := resp.Header.Get("Content-Type"); got != "image/svg+xml;utf-8" {
				t.Errorf("Content-Type = %q", got)
			}
			if got, want := readBody(t, resp), createBadge(t, tt.want); got != want {
				t.Errorf("body mismatch")
			}
		})
	}
}

func TestStaticMethodNotAllowed(t *testing.T) {
	srv := startServer(t)
	req, err := http.NewRequest(http.MethodPost, srv.appURL+"/static?subject=a&status=b", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", resp.StatusCode)
	}
}

func TestUnmatchedServiceBadge(t *testing.T) {
	srv := startServer(t)
	resp := get(t, srv.appURL+"/nope")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	want := createBadge(t, &badge.Params{Subject: "aegis", Status: "service not found"})
	if got := readBody(t, resp); got != want {
		t.Errorf("body mismatch")
	}
}

func TestRootRedirect(t *testing.T) {
	srv := startServer(t, "-root-redirect-url", "https://example.com/docs")
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(srv.appURL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("status = %d, want 302", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != "https://example.com/docs" {
		t.Errorf("Location = %q", loc)
	}
}

func TestSIGTERMExitsZero(t *testing.T) {
	appAddr := freeAddr(t)
	adminAddr := freeAddr(t)
	cmd := exec.Command(binPath,
		"-app-addr", appAddr,
		"-admin-addr", adminAddr,
		"-github-access-token", "e2e-token",
		"-log-level", "error",
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitReady(t, "http://"+adminAddr)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("exit: %v", err)
	}
}

func TestVersionFlag(t *testing.T) {
	out, err := exec.Command(binPath, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("version: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "aegis") {
		t.Errorf("output = %q", out)
	}
}
