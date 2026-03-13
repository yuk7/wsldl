package backup

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuk7/wsldl/lib/errutil"
	"github.com/yuk7/wsldl/lib/wsllib"
)

func TestParseArgs_NoArgs_AutoMode(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs(nil)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if !opts.auto {
		t.Fatal("auto = false, want true")
	}
}

func TestParseArgs_CustomRegFile_SetsRegPath(t *testing.T) {
	t.Parallel()

	opts, err := parseArgs([]string{"my-backup.reg"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if opts.regPath != "my-backup.reg" {
		t.Fatalf("regPath = %q, want %q", opts.regPath, "my-backup.reg")
	}
	if opts.tarPath != "" || opts.vhdxPath != "" {
		t.Fatalf("unexpected output paths: tar=%q vhdx=%q", opts.tarPath, opts.vhdxPath)
	}
}

func TestParseArgs_ExtensionRouting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		arg      string
		wantTar  string
		wantVhdx string
		wantReg  string
	}{
		{name: "tar", arg: "backup.tar", wantTar: "backup.tar"},
		{name: "tar.gz", arg: "backup.tar.gz", wantTar: "backup.tar.gz"},
		{name: "tgz", arg: "backup.tgz", wantTar: "backup.tgz"},
		{name: "vhdx", arg: "backup.ext4.vhdx", wantVhdx: "backup.ext4.vhdx"},
		{name: "vhdx.gz", arg: "backup.ext4.vhdx.gz", wantVhdx: "backup.ext4.vhdx.gz"},
		{name: "reg", arg: "backup.reg", wantReg: "backup.reg"},
		{name: "case insensitive", arg: "BACKUP.TAR.GZ", wantTar: "BACKUP.TAR.GZ"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, err := parseArgs([]string{tt.arg})
			if err != nil {
				t.Fatalf("parseArgs returned error: %v", err)
			}
			if opts.tarPath != tt.wantTar || opts.vhdxPath != tt.wantVhdx || opts.regPath != tt.wantReg {
				t.Fatalf("opts = %+v, want tar=%q vhdx=%q reg=%q", opts, tt.wantTar, tt.wantVhdx, tt.wantReg)
			}
		})
	}
}

func TestParseArgs_ShortOptionsRouting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want backupOptions
	}{
		{name: "tar", arg: "--tar", want: backupOptions{tarPath: "backup.tar"}},
		{name: "tgz", arg: "--tgz", want: backupOptions{tarPath: "backup.tar.gz"}},
		{name: "vhdx", arg: "--vhdx", want: backupOptions{vhdxPath: "backup.ext4.vhdx"}},
		{name: "vhdxgz", arg: "--vhdxgz", want: backupOptions{vhdxPath: "backup.ext4.vhdx.gz"}},
		{name: "reg", arg: "--reg", want: backupOptions{regPath: "backup.reg"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, err := parseArgs([]string{tt.arg})
			if err != nil {
				t.Fatalf("parseArgs returned error: %v", err)
			}
			if opts != tt.want {
				t.Fatalf("opts = %+v, want %+v", opts, tt.want)
			}
		})
	}
}

func TestParseArgs_InvalidInput_ReturnsInvalid(t *testing.T) {
	t.Parallel()

	_, err := parseArgs([]string{"backup.zip"})
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("error = %v, want %v", err, os.ErrInvalid)
	}
}

func TestGetCommandWithDeps_HelpVisibility(t *testing.T) {
	t.Parallel()

	wsl := wsllib.MockWslLib{
		IsDistributionRegisteredFunc: func(name string) bool {
			return false
		},
	}
	cmd := GetCommandWithDeps(wsl, wsllib.MockWslReg{})
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false")
	}
	if cmd.HelpText == nil {
		t.Fatal("HelpText is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("HelpText should not be empty")
	}
}

func TestGetCommand_WiresDefaultDeps(t *testing.T) {
	t.Parallel()

	cmd := GetCommand()
	if len(cmd.Names) != 1 || cmd.Names[0] != "backup" {
		t.Fatalf("Names = %v, want [backup]", cmd.Names)
	}
	if cmd.Visible == nil {
		t.Fatal("Visible is nil")
	}
	if cmd.Visible("Arch") {
		t.Fatal("Visible(Arch) = true, want false")
	}
	if cmd.HelpText == nil {
		t.Fatal("HelpText is nil")
	}
	if got := cmd.HelpText(); got == "" {
		t.Fatal("HelpText should not be empty")
	}
	if cmd.Run == nil {
		t.Fatal("Run is nil")
	}

	err := cmd.Run("Arch", []string{"--bad"})
	assertDisplayError(t, err)
}

func TestExecute_InvalidArgLength_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"a", "b"})
	assertDisplayError(t, err)
}

func TestExecute_InvalidSingleArg_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := execute(wsllib.MockWslLib{}, wsllib.MockWslReg{}, "Arch", []string{"--unknown"})
	assertDisplayError(t, err)
}

func TestExecute_RegOption_ProfileError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("profile failed")
	reg := wsllib.MockWslReg{
		GetProfileFromNameFunc: func(name string) (wsllib.Profile, error) {
			return wsllib.Profile{}, wantErr
		},
	}

	err := execute(wsllib.MockWslLib{}, reg, "Arch", []string{"--reg"})
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecute_VhdxCustomPath_Success(t *testing.T) {
	t.Parallel()

	destPath := filepath.Join(t.TempDir(), "backup.ext4.vhdx")
	gotName := ""
	gotDestPath := ""
	called := false

	err := executeWithBackups(
		wsllib.MockWslLib{},
		wsllib.MockWslReg{},
		"Arch",
		[]string{destPath},
		func(wsllib.WslReg, string, string) error {
			t.Fatal("backupReg should not be called")
			return nil
		},
		func(string, string) error {
			t.Fatal("backupTar should not be called")
			return nil
		},
		func(_ wsllib.WslReg, name, dest string) error {
			called = true
			gotName = name
			gotDestPath = dest
			return nil
		},
	)
	if err != nil {
		t.Fatalf("executeWithBackups returned error: %v", err)
	}
	if !called {
		t.Fatal("backupExt4Vhdx should be called")
	}
	if gotName != "Arch" {
		t.Fatalf("name = %q, want %q", gotName, "Arch")
	}
	if gotDestPath != destPath {
		t.Fatalf("dest path = %q, want %q", gotDestPath, destPath)
	}
}

func TestExecuteWithBackups_InvalidArgs_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	err := executeWithBackups(
		wsllib.MockWslLib{},
		wsllib.MockWslReg{},
		"Arch",
		[]string{"a", "b"},
		func(wsllib.WslReg, string, string) error { return nil },
		func(string, string) error { return nil },
		func(wsllib.WslReg, string, string) error { return nil },
	)
	assertDisplayError(t, err)
}

func TestExecuteWithBackupsOptions_AutoWSL2_UsesRegAndVhdxGz(t *testing.T) {
	t.Parallel()

	gotReg := ""
	gotVhdx := ""
	tarCalled := 0
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			if name != "Arch" {
				t.Fatalf("name = %q, want %q", name, "Arch")
			}
			return 0, 0, uint32(wsllib.FlagEnableWsl2), nil
		},
	}

	err := executeWithBackupsOptions(
		wsl,
		wsllib.MockWslReg{},
		"Arch",
		backupOptions{auto: true},
		func(_ wsllib.WslReg, _ string, path string) error {
			gotReg = path
			return nil
		},
		func(string, string) error {
			tarCalled++
			return nil
		},
		func(_ wsllib.WslReg, _ string, path string) error {
			gotVhdx = path
			return nil
		},
	)
	if err != nil {
		t.Fatalf("executeWithBackupsOptions returned error: %v", err)
	}
	if gotReg != "backup.reg" {
		t.Fatalf("reg path = %q, want %q", gotReg, "backup.reg")
	}
	if gotVhdx != "backup.ext4.vhdx.gz" {
		t.Fatalf("vhdx path = %q, want %q", gotVhdx, "backup.ext4.vhdx.gz")
	}
	if tarCalled != 0 {
		t.Fatalf("backupTar call count = %d, want 0", tarCalled)
	}
}

func TestExecuteWithBackupsOptions_AutoWSL1_UsesRegAndTarGz(t *testing.T) {
	t.Parallel()

	regCalled := 0
	vhdxCalled := 0
	gotTar := ""
	wsl := wsllib.MockWslLib{
		GetDistributionConfigurationFunc: func(name string) (uint32, uint64, uint32, error) {
			return 0, 0, 0, nil
		},
	}

	err := executeWithBackupsOptions(
		wsl,
		wsllib.MockWslReg{},
		"Arch",
		backupOptions{auto: true},
		func(wsllib.WslReg, string, string) error {
			regCalled++
			return nil
		},
		func(_ string, path string) error {
			gotTar = path
			return nil
		},
		func(wsllib.WslReg, string, string) error {
			vhdxCalled++
			return nil
		},
	)
	if err != nil {
		t.Fatalf("executeWithBackupsOptions returned error: %v", err)
	}
	if regCalled != 1 {
		t.Fatalf("backupReg call count = %d, want 1", regCalled)
	}
	if gotTar != "backup.tar.gz" {
		t.Fatalf("tar path = %q, want %q", gotTar, "backup.tar.gz")
	}
	if vhdxCalled != 0 {
		t.Fatalf("backupExt4Vhdx call count = %d, want 0", vhdxCalled)
	}
}

func TestExecuteWithBackupsOptions_TarError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("tar failed")
	err := executeWithBackupsOptions(
		wsllib.MockWslLib{},
		wsllib.MockWslReg{},
		"Arch",
		backupOptions{tarPath: "backup.tar.gz"},
		func(wsllib.WslReg, string, string) error { return nil },
		func(string, string) error { return wantErr },
		func(wsllib.WslReg, string, string) error { return nil },
	)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func TestExecuteWithBackupsOptions_VhdxError_ReturnsDisplayError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("vhdx failed")
	err := executeWithBackupsOptions(
		wsllib.MockWslLib{},
		wsllib.MockWslReg{},
		"Arch",
		backupOptions{vhdxPath: "backup.ext4.vhdx.gz"},
		func(wsllib.WslReg, string, string) error { return nil },
		func(string, string) error { return nil },
		func(wsllib.WslReg, string, string) error { return wantErr },
	)
	de := assertDisplayError(t, err)
	if !errors.Is(de, wantErr) {
		t.Fatalf("wrapped error = %v, want %v", de.Unwrap(), wantErr)
	}
}

func assertDisplayError(t *testing.T, err error) *errutil.DisplayError {
	t.Helper()
	var de *errutil.DisplayError
	if !errors.As(err, &de) {
		t.Fatalf("error type = %T, want *errutil.DisplayError", err)
	}
	return de
}
