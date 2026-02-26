package install

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

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

	err := Install(mockWsl, wsllib.MockWslReg{}, name, rootPath, "", false)
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

	err := Install(mockWsl, wsllib.MockWslReg{}, "TestDistro", rootPath, "deadbeef", false)
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

	if err := installWithDeps(mockWsl, mockReg, "TestDistro", sourcePath, "", false, deps); err != nil {
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
