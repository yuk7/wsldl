package fileutil

import (
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getenvFromMap(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}

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

func TestGetWindowsDirectoryFromEnv(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{
			name: "prefer SYSTEMROOT",
			env: map[string]string{
				"SYSTEMROOT": "D:\\Windows",
				"WINDIR":     "E:\\WinDir",
			},
			want: "D:\\Windows",
		},
		{
			name: "fallback to WINDIR",
			env: map[string]string{
				"WINDIR": "E:\\WinDir",
			},
			want: "E:\\WinDir",
		},
		{
			name: "fallback to default",
			env:  map[string]string{},
			want: "C:\\WINDOWS",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := getWindowsDirectoryFromEnv(getenvFromMap(tc.env))
			if got != tc.want {
				t.Fatalf("getWindowsDirectoryFromEnv() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetWindowsDirectory(t *testing.T) {
	t.Setenv("SYSTEMROOT", "D:\\Windows")
	t.Setenv("WINDIR", "E:\\WinDir")

	if got := GetWindowsDirectory(); got != "D:\\Windows" {
		t.Fatalf("GetWindowsDirectory() = %q, want %q", got, "D:\\Windows")
	}
}

func TestIsSpecialDir(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	systemRoot := filepath.Join(tmp, "system-root")
	userProfile := filepath.Join(tmp, "profile")

	userProfileSpecial, err := filepath.Abs(userProfile + "\\")
	if err != nil {
		t.Fatalf("filepath.Abs(userProfile) failed: %v", err)
	}
	system32Special, err := filepath.Abs(systemRoot + "\\System32")
	if err != nil {
		t.Fatalf("filepath.Abs(system32Special) failed: %v", err)
	}

	tests := []struct {
		name string
		cdir string
		env  map[string]string
		want bool
	}{
		{
			name: "matches userprofile special dir",
			cdir: userProfileSpecial,
			env: map[string]string{
				"USERPROFILE": userProfile,
			},
			want: true,
		},
		{
			name: "matches systemroot system32 special dir",
			cdir: system32Special,
			env: map[string]string{
				"SystemRoot": systemRoot,
			},
			want: true,
		},
		{
			name: "not special dir",
			cdir: filepath.Join(tmp, "regular"),
			env: map[string]string{
				"SystemRoot":  filepath.Join(tmp, "other"),
				"USERPROFILE": filepath.Join(tmp, "another"),
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isSpecialDir(tc.cdir, getenvFromMap(tc.env), filepath.Abs)
			if got != tc.want {
				t.Fatalf("isSpecialDir(%q) = %v, want %v", tc.cdir, got, tc.want)
			}
		})
	}
}

func TestIsCurrentDirSpecial_AbsError(t *testing.T) {
	if got := isCurrentDirSpecial(func(string) (string, error) {
		return "", errors.New("abs failed")
	}, getenvFromMap(map[string]string{})); !got {
		t.Fatal("isCurrentDirSpecial() = false, want true when abs fails")
	}
}

func TestIsSpecialDir_AbsError(t *testing.T) {
	if got := isSpecialDir("C:\\anywhere", getenvFromMap(map[string]string{}), func(string) (string, error) {
		return "", errors.New("abs failed")
	}); !got {
		t.Fatal("isSpecialDir() = false, want true when abs fails")
	}
}

func TestIsCurrentDirSpecial(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	base := t.TempDir()
	specialDir := filepath.Join(base, "profile\\")
	if err := os.MkdirAll(specialDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.Chdir(specialDir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Errorf("restore working directory failed: %v", err)
		}
	})
	cdir, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("filepath.Abs failed: %v", err)
	}
	userProfile := strings.TrimSuffix(cdir, "\\")

	t.Setenv("SystemDrive", "")
	t.Setenv("SystemRoot", "")
	t.Setenv("USERPROFILE", userProfile)

	if !IsCurrentDirSpecial() {
		candidate, _ := filepath.Abs(userProfile + "\\")
		t.Fatalf("IsCurrentDirSpecial() = false, want true (cdir=%q candidate=%q)", cdir, candidate)
	}
}

func TestIsCurrentDirSpecial_NotSpecial(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Errorf("restore working directory failed: %v", err)
		}
	})

	t.Setenv("SystemDrive", "")
	t.Setenv("SystemRoot", "")
	t.Setenv("USERPROFILE", "")

	if IsCurrentDirSpecial() {
		t.Fatal("IsCurrentDirSpecial() = true, want false")
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

func TestCopyFile_SourceOpenError(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "not-found.txt")
	destPath := filepath.Join(tmp, "dest.txt")

	if err := CopyFile(srcPath, destPath, false); err == nil {
		t.Fatal("CopyFile succeeded unexpectedly for missing source")
	}
}

func TestCopyFile_DestinationCreateError(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.txt")
	destPath := filepath.Join(tmp, "missing-dir", "dest.txt")
	if err := os.WriteFile(srcPath, []byte("payload"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := CopyFile(srcPath, destPath, false); err == nil {
		t.Fatal("CopyFile succeeded unexpectedly for invalid destination path")
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

func TestCopyFile_GzipSourceInvalid(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "src.gz")
	destPath := filepath.Join(tmp, "dest.txt")
	if err := os.WriteFile(srcPath, []byte("not a valid gzip stream"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := CopyFile(srcPath, destPath, false); err == nil {
		t.Fatal("CopyFile succeeded unexpectedly for invalid gzip source")
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
