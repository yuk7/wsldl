package download

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFile_WritesContentAndReturnsSHA256(t *testing.T) {
	t.Parallel()

	payload := []byte("download payload")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")

	sum, err := DownloadFile(server.URL, dest, 0)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	wantSumBytes := sha256.Sum256(payload)
	wantSum := hex.EncodeToString(wantSumBytes[:])
	if sum != wantSum {
		t.Fatalf("sha256 = %q, want %q", sum, wantSum)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read destination failed: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination payload = %q, want %q", got, payload)
	}
}

func TestDownloadFile_ProgressBarModes(t *testing.T) {
	payload := []byte("progress payload")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	tests := []struct {
		name  string
		width int
	}{
		{name: "positive width", width: 20},
		{name: "unknown size bar", width: -1},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			dest := filepath.Join(tmp, "file.bin")
			if _, err := DownloadFile(server.URL, dest, tc.width); err != nil {
				t.Fatalf("DownloadFile failed for width %d: %v", tc.width, err)
			}
		})
	}
}

func TestDownloadFile_OverwritesExistingFile(t *testing.T) {
	t.Parallel()

	payload := []byte("short")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(payload))
	}))
	defer server.Close()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")
	if err := os.WriteFile(dest, []byte("this is a much longer old content"), 0o600); err != nil {
		t.Fatalf("write old destination failed: %v", err)
	}

	if _, err := DownloadFile(server.URL, dest, 0); err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read destination failed: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination payload = %q, want %q", got, payload)
	}
}

func TestDownloadFile_InvalidURL(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")

	_, err := DownloadFile("://invalid-url", dest, 0)
	if err == nil {
		t.Fatal("DownloadFile succeeded unexpectedly")
	}
}
