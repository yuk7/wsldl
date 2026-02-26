package backup

import (
	"errors"
	"path/filepath"
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
