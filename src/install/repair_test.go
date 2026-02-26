package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestRepairRegistry_ProfileFoundFromBasePath_RewritesDistributionName(t *testing.T) {
	dir := executableDir(t)

	written := wsllib.Profile{}
	reg := wsllib.MockWslReg{
		GetProfileFromBasePathFunc: func(path string) (wsllib.Profile, error) {
			if path != dir {
				t.Fatalf("path = %q, want %q", path, dir)
			}
			return wsllib.Profile{
				BasePath:          dir,
				DistributionName:  "OldName",
				PackageFamilyName: "pkg",
			}, nil
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			written = profile
			return nil
		},
	}

	if err := repairRegistry(reg, "NewName"); err != nil {
		t.Fatalf("repairRegistry returned error: %v", err)
	}
	if written.BasePath != dir {
		t.Fatalf("written.BasePath = %q, want %q", written.BasePath, dir)
	}
	if written.DistributionName != "NewName" {
		t.Fatalf("written.DistributionName = %q, want %q", written.DistributionName, "NewName")
	}
	if written.PackageFamilyName != "pkg" {
		t.Fatalf("written.PackageFamilyName = %q, want %q", written.PackageFamilyName, "pkg")
	}
}

func TestRepairRegistry_Ext4VhdxExists_WritesWSL2Profile(t *testing.T) {
	dir := executableDir(t)
	ext4Path := dir + "\\ext4.vhdx"
	rootfsPath := dir + "\\rootfs"
	cleanupRepairFiles(t, ext4Path, rootfsPath)
	t.Cleanup(func() { cleanupRepairFiles(t, ext4Path, rootfsPath) })

	if err := os.WriteFile(ext4Path, []byte("vhdx"), 0o600); err != nil {
		t.Fatalf("write ext4.vhdx failed: %v", err)
	}

	written := wsllib.Profile{}
	reg := wsllib.MockWslReg{
		GetProfileFromBasePathFunc: func(path string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
		GenerateProfileFunc: func() wsllib.Profile {
			return wsllib.Profile{Flags: wsllib.FlagAppendNTPath}
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			written = profile
			return nil
		},
	}

	if err := repairRegistry(reg, "Arch"); err != nil {
		t.Fatalf("repairRegistry returned error: %v", err)
	}
	if written.DistributionName != "Arch" {
		t.Fatalf("written.DistributionName = %q, want %q", written.DistributionName, "Arch")
	}
	if written.BasePath != dir {
		t.Fatalf("written.BasePath = %q, want %q", written.BasePath, dir)
	}
	if written.Flags&wsllib.FlagEnableWsl2 != wsllib.FlagEnableWsl2 {
		t.Fatalf("written.Flags = %d, want WSL2 flag set", written.Flags)
	}
}

func TestRepairRegistry_RootfsExists_WritesWSL1Profile(t *testing.T) {
	dir := executableDir(t)
	ext4Path := dir + "\\ext4.vhdx"
	rootfsPath := dir + "\\rootfs"
	cleanupRepairFiles(t, ext4Path, rootfsPath)
	t.Cleanup(func() { cleanupRepairFiles(t, ext4Path, rootfsPath) })

	if err := os.WriteFile(rootfsPath, []byte("rootfs"), 0o600); err != nil {
		t.Fatalf("write rootfs failed: %v", err)
	}

	written := wsllib.Profile{}
	reg := wsllib.MockWslReg{
		GetProfileFromBasePathFunc: func(path string) (wsllib.Profile, error) {
			return wsllib.Profile{}, nil
		},
		GenerateProfileFunc: func() wsllib.Profile {
			return wsllib.Profile{Flags: wsllib.FlagEnableWsl2 | wsllib.FlagAppendNTPath}
		},
		WriteProfileFunc: func(profile wsllib.Profile) error {
			written = profile
			return nil
		},
	}

	if err := repairRegistry(reg, "Arch"); err != nil {
		t.Fatalf("repairRegistry returned error: %v", err)
	}
	if written.Flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
		t.Fatalf("written.Flags = %d, want WSL2 flag cleared", written.Flags)
	}
	if written.Flags&wsllib.FlagAppendNTPath != wsllib.FlagAppendNTPath {
		t.Fatalf("written.Flags = %d, want append-path preserved", written.Flags)
	}
}

func TestRepairRegistry_NoKnownInstallFiles_ReturnsError(t *testing.T) {
	dir := executableDir(t)
	ext4Path := dir + "\\ext4.vhdx"
	rootfsPath := dir + "\\rootfs"
	cleanupRepairFiles(t, ext4Path, rootfsPath)

	err := repairRegistry(wsllib.MockWslReg{}, "Arch")
	if err == nil {
		t.Fatal("repairRegistry succeeded unexpectedly")
	}
	if err.Error() != "repair failed" {
		t.Fatalf("err = %v, want %q", err, "repair failed")
	}
}

func executableDir(t *testing.T) string {
	t.Helper()

	efPath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable failed: %v", err)
	}
	return filepath.Dir(efPath)
}

func cleanupRepairFiles(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		_ = os.Remove(p)
	}
}
