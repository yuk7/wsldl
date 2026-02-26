package repair

import (
	"errors"
	"os"
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
