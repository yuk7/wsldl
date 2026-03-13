package install

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/yuk7/wsldl/lib/wsllib"
)

var installHTTPMu sync.Mutex

type installRoundTripFunc func(*http.Request) (*http.Response, error)

func (f installRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func withInstallMockTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()
	installHTTPMu.Lock()
	oldTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = rt
	t.Cleanup(func() {
		http.DefaultClient.Transport = oldTransport
		installHTTPMu.Unlock()
	})
}

func detectRootfsFilesFromExecutablePathForTest(executablePath string) (string, error) {
	efDir := filepath.Dir(executablePath)
	rootFile, err := detectRootfsFileName(os.DirFS(efDir))
	if err != nil {
		return "", err
	}
	if rootFile == "rootfs.tar.gz" {
		return rootFile, nil
	}
	return filepath.Join(efDir, rootFile), nil
}

func TestInstall_RoutesToTarRegister(t *testing.T) {
	t.Parallel()

	const (
		name     = "TestDistro"
		rootPath = "rootfs.tar"
	)

	called := 0
	mockWsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(gotName, gotRootPath string) error {
			called++
			if gotName != name {
				t.Fatalf("RegisterDistribution name = %q, want %q", gotName, name)
			}
			if gotRootPath != rootPath {
				t.Fatalf("RegisterDistribution rootPath = %q, want %q", gotRootPath, rootPath)
			}
			return nil
		},
	}

	err := Install(context.Background(), mockWsl, wsllib.MockWslReg{}, name, rootPath, "", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	if called != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", called)
	}
}

func TestInstall_ChecksumMismatchStopsBeforeRegister(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	rootPath := filepath.Join(tmp, "rootfs.tar")
	if err := os.WriteFile(rootPath, []byte("payload"), 0o600); err != nil {
		t.Fatalf("write rootfs: %v", err)
	}

	called := 0
	mockWsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	err := Install(context.Background(), mockWsl, wsllib.MockWslReg{}, "TestDistro", rootPath, "deadbeef", false)
	if err == nil {
		t.Fatal("Install succeeded unexpectedly; want checksum mismatch error")
	}
	if err.Error() != "checksum mismatch" {
		t.Fatalf("error = %q, want %q", err.Error(), "checksum mismatch")
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution was called %d times, want 0", called)
	}
}

func TestInstall_RoutesToExt4VhdxFlow(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	sourcePath := filepath.Join(tmp, "install.ext4.vhdx")
	basePath := filepath.Join(tmp, "distro")

	calls := make([]string, 0, 7)
	tempTarPath := ""
	mockWsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			calls = append(calls, "register")
			if filepath.Base(rootPath) != "em-vhdx-temp.tar" {
				t.Fatalf("temp tar path = %q, want basename %q", rootPath, "em-vhdx-temp.tar")
			}
			tempTarPath = rootPath
			return nil
		},
		UnregisterDistributionFunc: func(name string) error {
			calls = append(calls, "unregister")
			return nil
		},
	}

	var written wsllib.Profile
	mockReg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			calls = append(calls, "get-profile")
			return wsllib.Profile{
				BasePath: basePath,
				Flags:    0,
			}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			calls = append(calls, "write-profile")
			written = profile
			return nil
		},
	}

	var createdPath, removedPath, copiedSrc, copiedDest string
	gotCompress := true
	deps := installDeps{
		tempDir: func() string {
			return tmp
		},
		createFile: func(path string) (io.Closer, error) {
			calls = append(calls, "create-temp-tar")
			createdPath = path
			return nopCloser{}, nil
		},
		removeFile: func(path string) error {
			calls = append(calls, "remove-temp-tar")
			removedPath = path
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			calls = append(calls, "copy-vhdx")
			copiedSrc = srcPath
			copiedDest = destPath
			gotCompress = compress
			return nil
		},
	}

	if err := installWithDeps(context.Background(), mockWsl, mockReg, "TestDistro", sourcePath, "", false, deps); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if len(calls) != 7 {
		t.Fatalf("call sequence length = %d, want 7 (%v)", len(calls), calls)
	}
	wantCalls := []string{
		"create-temp-tar",
		"register",
		"remove-temp-tar",
		"get-profile",
		"unregister",
		"copy-vhdx",
		"write-profile",
	}
	for i, got := range calls {
		if got != wantCalls[i] {
			t.Fatalf("calls[%d] = %q, want %q (full=%v)", i, got, wantCalls[i], calls)
		}
	}
	if createdPath != tempTarPath {
		t.Fatalf("create temp path = %q, want %q", createdPath, tempTarPath)
	}
	if removedPath != tempTarPath {
		t.Fatalf("remove temp path = %q, want %q", removedPath, tempTarPath)
	}
	if copiedSrc != sourcePath {
		t.Fatalf("copied source path = %q, want %q", copiedSrc, sourcePath)
	}
	if copiedDest != filepath.Join(basePath, "ext4.vhdx") {
		t.Fatalf("copied dest path = %q, want %q", copiedDest, filepath.Join(basePath, "ext4.vhdx"))
	}
	if gotCompress {
		t.Fatalf("copy compress = %v, want false", gotCompress)
	}

	if written.Flags&wsllib.FlagEnableWsl2 != wsllib.FlagEnableWsl2 {
		t.Fatalf("written.Flags = %d, want WSL2 flag set", written.Flags)
	}
}

type nopCloser struct{}

func (nopCloser) Close() error { return nil }

func TestDetectRootfsFileName_PrioritizesDefaultOrder(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"rootfs.tar.gz": {Data: []byte("rootfs")},
		"install.tar":   {Data: []byte("install")},
	}

	got, err := detectRootfsFileName(fsys)
	if err != nil {
		t.Fatalf("detectRootfsFileName failed: %v", err)
	}
	if got != "install.tar" {
		t.Fatalf("detected root file = %q, want %q", got, "install.tar")
	}
}

func TestDetectRootfsFileName_ReturnsErrorWhenNotFound(t *testing.T) {
	t.Parallel()

	_, err := detectRootfsFileName(fstest.MapFS{})
	if err == nil {
		t.Fatal("detectRootfsFileName succeeded unexpectedly")
	}
}

func TestDetectRootfsFiles_ReturnsRootfsTarGzAsRelativeName(t *testing.T) {
	tmp := t.TempDir()
	exePath := filepath.Join(tmp, "wsldl-test.exe")
	if err := os.WriteFile(exePath, []byte("exe"), 0o600); err != nil {
		t.Fatalf("write executable file failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "rootfs.tar.gz"), []byte("rootfs"), 0o600); err != nil {
		t.Fatalf("write rootfs file failed: %v", err)
	}

	got, err := detectRootfsFilesFromExecutablePathForTest(exePath)
	if err != nil {
		t.Fatalf("detectRootfsFilesFromExecutablePathForTest failed: %v", err)
	}
	if got != "rootfs.tar.gz" {
		t.Fatalf("detected path = %q, want %q", got, "rootfs.tar.gz")
	}
}

func TestDetectRootfsFiles_ReturnsAbsolutePathForInstallTar(t *testing.T) {
	tmp := t.TempDir()
	exePath := filepath.Join(tmp, "wsldl-test.exe")
	if err := os.WriteFile(exePath, []byte("exe"), 0o600); err != nil {
		t.Fatalf("write executable file failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "install.tar"), []byte("rootfs"), 0o600); err != nil {
		t.Fatalf("write install tar failed: %v", err)
	}

	got, err := detectRootfsFilesFromExecutablePathForTest(exePath)
	if err != nil {
		t.Fatalf("detectRootfsFilesFromExecutablePathForTest failed: %v", err)
	}
	want := filepath.Join(tmp, "install.tar")
	if got != want {
		t.Fatalf("detected path = %q, want %q", got, want)
	}
}

func TestInstall_ContextCanceledBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := 0
	mockWsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	err := Install(ctx, mockWsl, wsllib.MockWslReg{}, "TestDistro", "rootfs.tar", "", false)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution call count = %d, want 0", called)
	}
}

func TestInstallWithDeps_NilContext_Works(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			if rootPath != "rootfs.tar" {
				t.Fatalf("rootPath = %q, want %q", rootPath, "rootfs.tar")
			}
			return nil
		},
	}

	err := installWithDeps(nil, wsl, wsllib.MockWslReg{}, "Arch", "rootfs.tar", "", false, defaultInstallDeps())
	if err != nil {
		t.Fatalf("installWithDeps returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", called)
	}
}

func TestInstallWithDeps_HttpTempDirEmpty_ReturnsError(t *testing.T) {
	t.Parallel()

	deps := installDeps{
		tempDir: func() string {
			return ""
		},
		createFile: func(path string) (io.Closer, error) {
			return nopCloser{}, nil
		},
		removeFile: func(path string) error {
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", "https://example.com/rootfs.tar.gz", "", false, deps)
	if err == nil {
		t.Fatal("installWithDeps succeeded unexpectedly")
	}
	if err.Error() != "failed to create temp directory" {
		t.Fatalf("error = %q, want %q", err.Error(), "failed to create temp directory")
	}
}

func TestInstallExt4Vhdx_WrapperSuccess(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	sourcePath := filepath.Join(tmp, "install.ext4.vhdx")
	if err := os.WriteFile(sourcePath, []byte("vhdx"), 0o600); err != nil {
		t.Fatalf("write source vhdx failed: %v", err)
	}
	basePath := filepath.Join(tmp, "distro")
	if err := os.MkdirAll(basePath, 0o700); err != nil {
		t.Fatalf("mkdir basePath failed: %v", err)
	}

	registerCalled := 0
	unregisterCalled := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			registerCalled++
			return nil
		},
		UnregisterDistributionFunc: func(name string) error {
			unregisterCalled++
			return nil
		},
	}

	written := wsllib.Profile{}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			written = profile
			return nil
		},
	}

	err := InstallExt4Vhdx(wsl, reg, "Arch", sourcePath)
	if err != nil {
		t.Fatalf("InstallExt4Vhdx returned error: %v", err)
	}
	if registerCalled != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", registerCalled)
	}
	if unregisterCalled != 1 {
		t.Fatalf("UnregisterDistribution call count = %d, want 1", unregisterCalled)
	}
	if written.Flags&wsllib.FlagEnableWsl2 != wsllib.FlagEnableWsl2 {
		t.Fatalf("written.Flags = %d, want WSL2 flag set", written.Flags)
	}
	if _, err := os.Stat(filepath.Join(basePath, "ext4.vhdx")); err != nil {
		t.Fatalf("ext4.vhdx was not copied: %v", err)
	}
}

func TestInstallExt4VhdxWithDeps_ProfileLookupError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile lookup failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc:   func(name, rootPath string) error { return nil },
		UnregisterDistributionFunc: func(name string) error { return nil },
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	deps := installDeps{
		tempDir: func() string {
			return t.TempDir()
		},
		createFile: func(path string) (io.Closer, error) {
			return nopCloser{}, nil
		},
		removeFile: func(path string) error {
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			t.Fatal("copyFile should not be called when profile lookup fails")
			return nil
		},
	}

	err := installExt4VhdxWithDeps(wsl, reg, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallWithDeps_HTTPDownloadPath_Success(t *testing.T) {
	payload := []byte("rootfs")
	withInstallMockTransport(t, installRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(string(payload))),
			ContentLength: int64(len(payload)),
			Header:        make(http.Header),
		}, nil
	}))

	tmp := t.TempDir()
	var gotRootPath string
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			gotRootPath = rootPath
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return nil
		},
	}

	removed := ""
	deps := installDeps{
		tempDir: tmpDirConst(tmp),
		createFile: func(path string) (io.Closer, error) {
			return nopCloser{}, nil
		},
		removeFile: func(path string) error {
			removed = path
			return os.Remove(path)
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", "http://example.com/rootfs.tar.gz", "", false, deps)
	if err != nil {
		t.Fatalf("installWithDeps returned error: %v", err)
	}
	if gotRootPath == "" {
		t.Fatal("RegisterDistribution did not receive downloaded root path")
	}
	if filepath.Dir(gotRootPath) != tmp {
		t.Fatalf("download path dir = %q, want %q", filepath.Dir(gotRootPath), tmp)
	}
	if !strings.HasSuffix(gotRootPath, "rootfs.tar.gz") {
		t.Fatalf("download path = %q, want suffix %q", gotRootPath, "rootfs.tar.gz")
	}
	if removed == "" {
		t.Fatal("temporary file was not scheduled for removal")
	}
	if _, statErr := os.Stat(removed); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("temporary file still exists: stat err = %v", statErr)
	}
}

func TestInstallWithDeps_HTTPDownloadPath_DownloadError(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	tmp := t.TempDir()
	deps := installDeps{
		tempDir: tmpDirConst(tmp),
		createFile: func(path string) (io.Closer, error) {
			return nopCloser{}, nil
		},
		removeFile: func(path string) error {
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", "http://127.0.0.1:1/rootfs.tar.gz", "", false, deps)
	if err == nil {
		t.Fatal("installWithDeps succeeded unexpectedly")
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution call count = %d, want 0", called)
	}
}

func TestInstallWithDeps_SHA256Path_FileOpenError(t *testing.T) {
	t.Parallel()

	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", filepath.Join(t.TempDir(), "missing.tar"), "abcd", false, defaultInstallDeps())
	if err == nil {
		t.Fatal("installWithDeps succeeded unexpectedly")
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution call count = %d, want 0", called)
	}
}

func TestInstallWithDeps_SHA256Path_Success(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	rootPath := filepath.Join(tmp, "rootfs.tar")
	payload := []byte("payload")
	if err := os.WriteFile(rootPath, payload, 0o600); err != nil {
		t.Fatalf("write rootfs failed: %v", err)
	}
	sumRaw := sha256.Sum256(payload)
	sum := hex.EncodeToString(sumRaw[:])

	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, gotRootPath string) error {
			called++
			if gotRootPath != rootPath {
				t.Fatalf("rootPath = %q, want %q", gotRootPath, rootPath)
			}
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", rootPath, sum, false, defaultInstallDeps())
	if err != nil {
		t.Fatalf("installWithDeps returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", called)
	}
}

func TestInstallExt4VhdxWithDeps_TempDirEmpty(t *testing.T) {
	t.Parallel()

	deps := installDeps{
		tempDir: func() string { return "" },
		createFile: func(path string) (io.Closer, error) {
			t.Fatal("createFile should not be called when tempDir is empty")
			return nil, nil
		},
		removeFile: func(path string) error { return nil },
		copyFile:   func(srcPath, destPath string, compress bool) error { return nil },
	}

	err := installExt4VhdxWithDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", "install.ext4.vhdx", deps)
	if err == nil || err.Error() != "failed to create temp directory" {
		t.Fatalf("error = %v, want %q", err, "failed to create temp directory")
	}
}

func TestInstallExt4VhdxWithDeps_CreateFileError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("create failed")
	deps := installDeps{
		tempDir: func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) {
			return nil, wantErr
		},
		removeFile: func(path string) error { return nil },
		copyFile:   func(srcPath, destPath string, compress bool) error { return nil },
	}

	err := installExt4VhdxWithDeps(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallExt4VhdxWithDeps_RegisterError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("register failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error { return wantErr },
	}
	deps := installDeps{
		tempDir:    func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return nil },
		copyFile:   func(srcPath, destPath string, compress bool) error { return nil },
	}

	err := installExt4VhdxWithDeps(wsl, wsllib.MockWslReg{}, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallExt4VhdxWithDeps_UnregisterError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("unregister failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc:   func(name, rootPath string) error { return nil },
		UnregisterDistributionFunc: func(name string) error { return wantErr },
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: t.TempDir()}, nil
		},
	}
	deps := installDeps{
		tempDir:    func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return nil },
		copyFile: func(srcPath, destPath string, compress bool) error {
			t.Fatal("copyFile should not be called when unregister fails")
			return nil
		},
	}

	err := installExt4VhdxWithDeps(wsl, reg, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallExt4VhdxWithDeps_CopyError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("copy failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc:   func(name, rootPath string) error { return nil },
		UnregisterDistributionFunc: func(name string) error { return nil },
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: t.TempDir()}, nil
		},
	}
	deps := installDeps{
		tempDir:    func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return nil },
		copyFile: func(srcPath, destPath string, compress bool) error {
			return wantErr
		},
	}

	err := installExt4VhdxWithDeps(wsl, reg, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallExt4VhdxWithDeps_WriteProfileError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("write profile failed")
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc:   func(name, rootPath string) error { return nil },
		UnregisterDistributionFunc: func(name string) error { return nil },
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: t.TempDir()}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			return wantErr
		},
	}
	deps := installDeps{
		tempDir:    func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return nil },
		copyFile:   func(srcPath, destPath string, compress bool) error { return nil },
	}

	err := installExt4VhdxWithDeps(wsl, reg, "Arch", "install.ext4.vhdx", deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestInstallExt4VhdxWithDeps_EmptyBasePathAndNilErr_ReturnsNil(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error { return nil },
	}
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
	}
	deps := installDeps{
		tempDir:    func() string { return t.TempDir() },
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return nil },
		copyFile: func(srcPath, destPath string, compress bool) error {
			t.Fatal("copyFile should not be called when base path is empty")
			return nil
		},
	}

	if err := installExt4VhdxWithDeps(wsl, reg, "Arch", "install.ext4.vhdx", deps); err != nil {
		t.Fatalf("error = %v, want nil", err)
	}
}

func tmpDirConst(path string) func() string {
	return func() string { return path }
}

type secondErrContext struct {
	calls int
}

func (c *secondErrContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c *secondErrContext) Done() <-chan struct{} {
	return nil
}

func (c *secondErrContext) Err() error {
	c.calls++
	if c.calls >= 2 {
		return context.Canceled
	}
	return nil
}

func (c *secondErrContext) Value(key any) any {
	return nil
}

func TestInstallWithDeps_ShowProgress_SHA256Success(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	rootPath := filepath.Join(tmp, "rootfs.tar")
	payload := []byte("payload")
	if err := os.WriteFile(rootPath, payload, 0o600); err != nil {
		t.Fatalf("write rootfs failed: %v", err)
	}
	sumRaw := sha256.Sum256(payload)
	sum := hex.EncodeToString(sumRaw[:])

	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, gotRootPath string) error {
			called++
			if gotRootPath != rootPath {
				t.Fatalf("rootPath = %q, want %q", gotRootPath, rootPath)
			}
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", rootPath, sum, true, defaultInstallDeps())
	if err != nil {
		t.Fatalf("installWithDeps returned error: %v", err)
	}
	if called != 1 {
		t.Fatalf("RegisterDistribution call count = %d, want 1", called)
	}
}

func TestInstallWithDeps_ShowProgress_HTTPDownloadSuccess(t *testing.T) {
	payload := []byte("rootfs")
	withInstallMockTransport(t, installRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(string(payload))),
			ContentLength: int64(len(payload)),
			Header:        make(http.Header),
		}, nil
	}))

	tmp := t.TempDir()
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			return nil
		},
	}
	deps := installDeps{
		tempDir:    tmpDirConst(tmp),
		createFile: func(path string) (io.Closer, error) { return nopCloser{}, nil },
		removeFile: func(path string) error { return os.Remove(path) },
		copyFile:   func(srcPath, destPath string, compress bool) error { return nil },
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", "http://example.com/rootfs.tar.gz", "", true, deps)
	if err != nil {
		t.Fatalf("installWithDeps returned error: %v", err)
	}
}

func TestInstallWithDeps_SHA256Path_CopyError(t *testing.T) {
	t.Parallel()

	dirPath := t.TempDir()
	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	err := installWithDeps(context.Background(), wsl, wsllib.MockWslReg{}, "Arch", dirPath, "abcd", false, defaultInstallDeps())
	if err == nil {
		t.Fatal("installWithDeps succeeded unexpectedly")
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution call count = %d, want 0", called)
	}
}

func TestInstallWithDeps_ContextCanceledAfterPreparation(t *testing.T) {
	t.Parallel()

	ctx := &secondErrContext{}
	called := 0
	wsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			called++
			return nil
		},
	}

	err := installWithDeps(ctx, wsl, wsllib.MockWslReg{}, "Arch", "rootfs.tar", "", false, defaultInstallDeps())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}
	if called != 0 {
		t.Fatalf("RegisterDistribution call count = %d, want 0", called)
	}
}

func TestDetectRootfsFiles_ReturnsErrorWhenNotFound(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	exePath := filepath.Join(tmp, "wsldl-test.exe")
	if err := os.WriteFile(exePath, []byte("exe"), 0o600); err != nil {
		t.Fatalf("write executable file failed: %v", err)
	}

	_, err := detectRootfsFilesFromExecutablePathForTest(exePath)
	if err == nil {
		t.Fatal("detectRootfsFiles succeeded unexpectedly")
	}
}
