package backup

import (
	"compress/gzip"
	"errors"
	"io"
	"os"
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
	srcPath := basePath + "\\ext4.vhdx"
	payload := []byte("plain-vhdx")
	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source vhdx: %v", err)
	}

	destPath := filepath.Join(tmp, "backup.ext4.vhdx")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
	}

	if err := backupExt4Vhdx(reg, "Arch", destPath); err != nil {
		t.Fatalf("backupExt4Vhdx returned error: %v", err)
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("destination = %q, want %q", got, payload)
	}
}

func TestBackupExt4Vhdx_CompressesWhenGzSuffix(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	basePath := filepath.Join(tmp, "base")
	srcPath := basePath + "\\ext4.vhdx"
	payload := []byte("compress-vhdx")
	if err := os.WriteFile(srcPath, payload, 0o600); err != nil {
		t.Fatalf("write source vhdx: %v", err)
	}

	destPath := filepath.Join(tmp, "backup.ext4.vhdx.gz")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{BasePath: basePath}, nil
		},
	}

	if err := backupExt4Vhdx(reg, "Arch", destPath); err != nil {
		t.Fatalf("backupExt4Vhdx returned error: %v", err)
	}

	got, err := readGzipFile(destPath)
	if err != nil {
		t.Fatalf("read gzip destination: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("gzip payload = %q, want %q", got, payload)
	}
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
