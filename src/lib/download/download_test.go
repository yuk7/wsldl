package download

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var defaultClientMu sync.Mutex

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func withDefaultClientTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()
	defaultClientMu.Lock()
	oldTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = rt
	t.Cleanup(func() {
		http.DefaultClient.Transport = oldTransport
		defaultClientMu.Unlock()
	})
}

func newHTTPResponse(payload []byte, contentLength int64) *http.Response {
	return &http.Response{
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewReader(payload)),
		ContentLength: contentLength,
		Header:        make(http.Header),
	}
}

type errReadCloser struct {
	err error
}

func (e errReadCloser) Read(p []byte) (int, error) {
	return 0, e.err
}

func (e errReadCloser) Close() error {
	return nil
}

func TestDownloadFile_WritesContentAndReturnsSHA256(t *testing.T) {
	payload := []byte("download payload")
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(payload, int64(len(payload))), nil
	}))

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")

	sum, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, 0)
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
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(payload, int64(len(payload))), nil
	}))

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
			if _, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, tc.width); err != nil {
				t.Fatalf("DownloadFile failed for width %d: %v", tc.width, err)
			}
		})
	}
}

func TestDownloadFile_OverwritesExistingFile(t *testing.T) {
	payload := []byte("short")
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(payload, int64(len(payload))), nil
	}))

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")
	if err := os.WriteFile(dest, []byte("this is a much longer old content"), 0o600); err != nil {
		t.Fatalf("write old destination failed: %v", err)
	}

	if _, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, 0); err != nil {
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

	_, err := DownloadFile(context.Background(), "://invalid-url", dest, 0)
	if err == nil {
		t.Fatal("DownloadFile succeeded unexpectedly")
	}
}

func TestDownloadFile_ContextCanceled(t *testing.T) {
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, req.Context().Err()
	}))

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "file.bin")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := DownloadFile(ctx, "http://example.com/file.bin", dest, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}
}

func TestDownloadFile_NilContext_UsesBackground(t *testing.T) {
	payload := []byte("nil context payload")
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(payload, int64(len(payload))), nil
	}))

	dest := filepath.Join(t.TempDir(), "file.bin")
	if _, err := DownloadFile(nil, "http://example.com/file.bin", dest, 0); err != nil {
		t.Fatalf("DownloadFile(nil, ...) returned error: %v", err)
	}
}

func TestDownloadFile_OpenDestinationError(t *testing.T) {
	payload := []byte("payload")
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(payload, int64(len(payload))), nil
	}))

	dest := filepath.Join(t.TempDir(), "missing", "file.bin")
	_, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, 0)
	if err == nil {
		t.Fatal("DownloadFile succeeded unexpectedly")
	}
}

func TestDownloadFile_ContentLengthZero_WithProgressBar(t *testing.T) {
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return newHTTPResponse(nil, 0), nil
	}))

	dest := filepath.Join(t.TempDir(), "file.bin")
	sum, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, 10)
	if err != nil {
		t.Fatalf("DownloadFile returned error: %v", err)
	}
	if sum != hex.EncodeToString(sha256.New().Sum(nil)) {
		t.Fatalf("sha256 = %q, want empty payload sha256", sum)
	}
}

func TestDownloadFile_ReadError_ReturnsError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("read failed")
	withDefaultClientTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          errReadCloser{err: wantErr},
			ContentLength: 10,
			Header:        make(http.Header),
		}, nil
	}))

	dest := filepath.Join(t.TempDir(), "file.bin")
	_, err := DownloadFile(context.Background(), "http://example.com/file.bin", dest, 0)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
