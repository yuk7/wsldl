package version

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetCommand(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) != 3 {
		t.Fatalf("names length = %d, want 3", len(cmd.Names))
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}
}

func TestExecute_PrintsVersionInfo(t *testing.T) {
	oldProject, oldVersion, oldURL := project, version, url
	project = "wsldl2-test"
	version = "1.2.3"
	url = "https://example.test/wsldl2"
	t.Cleanup(func() {
		project, version, url = oldProject, oldVersion, oldURL
	})

	out := captureStdout(t, execute)
	if !strings.Contains(out, "wsldl2-test, version 1.2.3  ("+runtime.GOARCH+")") {
		t.Fatalf("output missing version line: %q", out)
	}
	if !strings.Contains(out, "https://example.test/wsldl2") {
		t.Fatalf("output missing url line: %q", out)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	return string(b)
}
