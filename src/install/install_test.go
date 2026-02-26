package install

import (
	"os"
	"path/filepath"
	"strings"
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
	payload := []byte("fake-vhdx-content")
	if err := os.WriteFile(sourcePath, payload, 0o600); err != nil {
		t.Fatalf("write source vhdx: %v", err)
	}

	basePath := filepath.Join(tmp, "distro")
	if err := os.MkdirAll(basePath, 0o700); err != nil {
		t.Fatalf("mkdir basePath: %v", err)
	}

	calls := make([]string, 0, 4)
	mockWsl := wsllib.MockWslLib{
		RegisterDistributionFunc: func(name, rootPath string) error {
			calls = append(calls, "register")
			if !strings.HasSuffix(strings.ToLower(rootPath), "\\em-vhdx-temp.tar") {
				t.Fatalf("temp tar path = %q, want suffix %q", rootPath, "\\em-vhdx-temp.tar")
			}
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

	if err := Install(mockWsl, mockReg, "TestDistro", sourcePath, "", false); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if len(calls) != 4 {
		t.Fatalf("call sequence length = %d, want 4 (%v)", len(calls), calls)
	}
	if calls[0] != "register" || calls[1] != "get-profile" || calls[2] != "unregister" || calls[3] != "write-profile" {
		t.Fatalf("call sequence = %v, want [register get-profile unregister write-profile]", calls)
	}

	if written.Flags&wsllib.FlagEnableWsl2 != wsllib.FlagEnableWsl2 {
		t.Fatalf("written.Flags = %d, want WSL2 flag set", written.Flags)
	}

	destPath := basePath + "\\ext4.vhdx"
	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read copied vhdx: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("copied vhdx payload = %q, want %q", got, payload)
	}
}
