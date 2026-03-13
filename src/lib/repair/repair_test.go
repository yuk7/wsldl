package repair

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestIsInstalledFilesExistInDir(t *testing.T) {
	t.Parallel()

	errNotFound := errors.New("not found")

	tests := []struct {
		name string
		stat func(name string) (os.FileInfo, error)
		want bool
	}{
		{
			name: "vhdx exists",
			stat: func(name string) (os.FileInfo, error) {
				if name == "X:\\dir\\ext4.vhdx" {
					return nil, nil
				}
				return nil, errNotFound
			},
			want: true,
		},
		{
			name: "rootfs exists",
			stat: func(name string) (os.FileInfo, error) {
				if name == "X:\\dir\\rootfs" {
					return nil, nil
				}
				return nil, errNotFound
			},
			want: true,
		},
		{
			name: "neither exists",
			stat: func(name string) (os.FileInfo, error) {
				return nil, errNotFound
			},
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isInstalledFilesExistInDir("X:\\dir", tc.stat)
			if got != tc.want {
				t.Fatalf("isInstalledFilesExistInDir() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsInstalledFilesExist(t *testing.T) {
	oldExecutablePathFunc := executablePathFunc
	oldStatPathFunc := statPathFunc
	t.Cleanup(func() {
		executablePathFunc = oldExecutablePathFunc
		statPathFunc = oldStatPathFunc
	})

	executablePathFunc = func() string {
		return "X:\\dir\\wsldl.exe"
	}

	t.Run("vhdx exists", func(t *testing.T) {
		statPathFunc = func(name string) (os.FileInfo, error) {
			if name == filepath.Dir("X:\\dir\\wsldl.exe")+"\\ext4.vhdx" {
				return nil, nil
			}
			return nil, errors.New("not found")
		}
		if !IsInstalledFilesExist() {
			t.Fatal("IsInstalledFilesExist() = false, want true")
		}
	})

	t.Run("not exists", func(t *testing.T) {
		statPathFunc = func(name string) (os.FileInfo, error) {
			return nil, errors.New("not found")
		}
		if IsInstalledFilesExist() {
			t.Fatal("IsInstalledFilesExist() = true, want false")
		}
	})
}
