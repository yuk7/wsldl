package backup

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestBackupReg_GetProfileError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	err := backupReg(reg, "Arch", "backup.reg")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestBackupReg_RegExportCommandError(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{
				UUID: "00000000-0000-0000-0000-000000000000",
			}, nil
		},
	}

	err := backupReg(reg, "Arch", filepath.Join(t.TempDir(), "backup.reg"))
	if err == nil {
		t.Fatal("backupReg succeeded unexpectedly")
	}
}

func TestBackupTar_GzipExportCommandError(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "backup.tar.gz")
	err := backupTar("wsldl-test-missing-distro-9f85b7f2", dest)
	if err == nil {
		t.Fatal("backupTar succeeded unexpectedly for gzip destination")
	}
}

func TestBackupTar_PlainExportCommandError(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "backup.tar")
	err := backupTar("wsldl-test-missing-distro-9f85b7f2", dest)
	if err == nil {
		t.Fatal("backupTar succeeded unexpectedly for plain destination")
	}
}

func TestBackupTarWithDeps_GzipSuccessCopiesCompressedTar(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "backup.tar.gz")
	var calls []string
	var gotExportDistribution string
	var gotExportDest string
	var gotCopySrc string
	var gotCopyDest string
	var gotCopyCompress bool
	var gotRemoved string

	deps := backupTarDeps{
		tempDir: func() string { return t.TempDir() },
		export: func(distributionName, destFileName string) error {
			calls = append(calls, "export")
			gotExportDistribution = distributionName
			gotExportDest = destFileName
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			calls = append(calls, "copy")
			gotCopySrc = srcPath
			gotCopyDest = destPath
			gotCopyCompress = compress
			return nil
		},
		remove: func(path string) error {
			calls = append(calls, "remove")
			gotRemoved = path
			return nil
		},
		randIntn: func(n int) int { return 42 },
	}

	if err := backupTarWithDeps("Arch", dest, deps); err != nil {
		t.Fatalf("backupTarWithDeps returned error: %v", err)
	}

	wantCalls := []string{"export", "copy", "remove"}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
	if gotExportDistribution != "Arch" {
		t.Fatalf("distribution = %q, want %q", gotExportDistribution, "Arch")
	}
	wantTmpTar := filepath.Join(filepath.Dir(gotExportDest), "42.tar")
	if gotExportDest != wantTmpTar {
		t.Fatalf("export dest = %q, want %q", gotExportDest, wantTmpTar)
	}
	if gotCopySrc != gotExportDest {
		t.Fatalf("copy src = %q, want %q", gotCopySrc, gotExportDest)
	}
	if gotCopyDest != dest {
		t.Fatalf("copy dest = %q, want %q", gotCopyDest, dest)
	}
	if !gotCopyCompress {
		t.Fatal("copy compress = false, want true")
	}
	if gotRemoved != gotExportDest {
		t.Fatalf("remove path = %q, want %q", gotRemoved, gotExportDest)
	}
}

func TestBackupTarWithDeps_GzipTempDirEmptyReturnsError(t *testing.T) {
	t.Parallel()

	deps := backupTarDeps{
		tempDir: func() string { return "" },
		export: func(distributionName, destFileName string) error {
			t.Fatal("export should not be called when temp dir is empty")
			return nil
		},
		copyFile: func(srcPath, destPath string, compress bool) error {
			t.Fatal("copyFile should not be called when temp dir is empty")
			return nil
		},
		remove:   func(path string) error { return nil },
		randIntn: func(n int) int { return 0 },
	}

	err := backupTarWithDeps("Arch", "backup.tar.gz", deps)
	if err == nil {
		t.Fatal("backupTarWithDeps succeeded unexpectedly")
	}
	if err.Error() != "failed to create temp directory" {
		t.Fatalf("error = %q, want %q", err.Error(), "failed to create temp directory")
	}
}

func TestBackupExt4Vhdx_GetProfileError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	err := backupExt4Vhdx(reg, "Arch", "backup.ext4.vhdx")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestBackupExt4Vhdx_GetProfileErrorEvenWithBasePath(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: "C:\\WSL\\Arch"}, wantErr
		},
	}

	err := backupExt4Vhdx(reg, "Arch", "backup.ext4.vhdx")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestBackupExt4Vhdx_EmptyBasePathReturnsError(t *testing.T) {
	t.Parallel()

	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: ""}, nil
		},
	}

	err := backupExt4Vhdx(reg, "Arch", "backup.ext4.vhdx")
	if err == nil {
		t.Fatal("backupExt4Vhdx succeeded unexpectedly")
	}
	if err.Error() != "get profile failed" {
		t.Fatalf("error = %q, want %q", err.Error(), "get profile failed")
	}
}

func TestBackupExt4Vhdx_CopiesPlainFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	basePath := filepath.Join(tmp, "base")
	destPath := filepath.Join(tmp, "backup.ext4.vhdx")
	var gotSrcPath string
	var gotDestPath string
	var gotCompress bool
	called := false
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
	}

	if err := backupExt4VhdxWithCopy(reg, "Arch", destPath, func(srcPath, destPath string, compress bool) error {
		called = true
		gotSrcPath = srcPath
		gotDestPath = destPath
		gotCompress = compress
		return nil
	}); err != nil {
		t.Fatalf("backupExt4Vhdx returned error: %v", err)
	}

	if !called {
		t.Fatal("copyFile was not called")
	}
	wantSrcPath := filepath.Join(basePath, "ext4.vhdx")
	if gotSrcPath != wantSrcPath {
		t.Fatalf("srcPath = %q, want %q", gotSrcPath, wantSrcPath)
	}
	if gotDestPath != destPath {
		t.Fatalf("destPath = %q, want %q", gotDestPath, destPath)
	}
	if gotCompress {
		t.Fatal("compress = true, want false")
	}
}

func TestBackupExt4Vhdx_CompressesWhenGzSuffix(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	basePath := filepath.Join(tmp, "base")
	destPath := filepath.Join(tmp, "backup.ext4.vhdx.gz")
	var gotCompress bool
	called := false
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
	}

	if err := backupExt4VhdxWithCopy(reg, "Arch", destPath, func(_, _ string, compress bool) error {
		called = true
		gotCompress = compress
		return nil
	}); err != nil {
		t.Fatalf("backupExt4Vhdx returned error: %v", err)
	}
	if !called {
		t.Fatal("copyFile was not called")
	}
	if !gotCompress {
		t.Fatal("compress = false, want true")
	}
}
