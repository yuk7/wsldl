package fileutil

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestDQEscapeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no spaces",
			in:   "hello",
			want: "hello",
		},
		{
			name: "spaces only",
			in:   "hello world",
			want: "\"hello world\"",
		},
		{
			name: "spaces and quotes",
			in:   "hello \"world\"",
			want: "\"hello \\\"world\\\"\"",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DQEscapeString(tc.in)
			if got != tc.want {
				t.Fatalf("DQEscapeString(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestCopyFile_PlainToPlain(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.txt")
	destPath := filepath.Join(tmp, "dest.txt")
	payload := []byte("plain text data")

	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := CopyFile(srcPath, destPath, false); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination = %q, want %q", got, payload)
	}
}

func TestCopyFile_GzipSourceToPlain(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.gz")
	destPath := filepath.Join(tmp, "dest.txt")
	payload := []byte("decompress me")

	if err := writeGzipFile(srcPath, payload); err != nil {
		t.Fatalf("write gzip source: %v", err)
	}

	if err := CopyFile(srcPath, destPath, false); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination = %q, want %q", got, payload)
	}
}

func TestCopyFile_CompressDestination(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.txt")
	destPath := filepath.Join(tmp, "dest.gz")
	payload := []byte("compress me")

	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := CopyFile(srcPath, destPath, true); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	got, err := readGzipFile(destPath)
	if err != nil {
		t.Fatalf("read gzip destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("gzip destination payload = %q, want %q", got, payload)
	}
}

func TestCopyFileAndCompress_CompressesByExtension(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.txt")
	destPath := filepath.Join(tmp, "dest.tgz")
	payload := []byte("auto compress by extension")

	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := CopyFileAndCompress(srcPath, destPath); err != nil {
		t.Fatalf("CopyFileAndCompress failed: %v", err)
	}

	got, err := readGzipFile(destPath)
	if err != nil {
		t.Fatalf("read gzip destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("gzip destination payload = %q, want %q", got, payload)
	}
}

func writeGzipFile(path string, payload []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	if _, err := gw.Write(payload); err != nil {
		gw.Close()
		return err
	}
	return gw.Close()
}

func readGzipFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	return io.ReadAll(gr)
}
