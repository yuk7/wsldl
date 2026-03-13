//go:build windows

package wsllib

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	wslreg "github.com/yuk7/wslreglib-go"
)

var implHooksMu sync.Mutex

func withImplHooksLocked(t *testing.T) {
	t.Helper()
	implHooksMu.Lock()
	t.Cleanup(implHooksMu.Unlock)
}

func TestNativeWslLib_DelegatesAllCalls(t *testing.T) {
	withImplHooksLocked(t)

	origIsReg := wslIsDistributionRegisteredFunc
	origRegister := wslRegisterDistributionFunc
	origUnregister := wslUnregisterDistributionFunc
	origLaunchInteractive := wslLaunchInteractiveFunc
	origLaunch := wslLaunchFunc
	origGetConfig := wslGetDistributionConfigurationFunc
	origConfigure := wslConfigureDistributionFunc
	t.Cleanup(func() {
		wslIsDistributionRegisteredFunc = origIsReg
		wslRegisterDistributionFunc = origRegister
		wslUnregisterDistributionFunc = origUnregister
		wslLaunchInteractiveFunc = origLaunchInteractive
		wslLaunchFunc = origLaunch
		wslGetDistributionConfigurationFunc = origGetConfig
		wslConfigureDistributionFunc = origConfigure
	})

	wslIsDistributionRegisteredFunc = func(name string) bool {
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		return true
	}
	wslRegisterDistributionFunc = func(name, rootPath string) error {
		if name != "Arch" || rootPath != "rootfs.tar" {
			t.Fatalf("register args = (%q, %q)", name, rootPath)
		}
		return nil
	}
	wslUnregisterDistributionFunc = func(name string) error {
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		return nil
	}
	wslLaunchInteractiveFunc = func(name, command string, inheritPath bool) (uint32, error) {
		if name != "Arch" || command != "echo hello" || !inheritPath {
			t.Fatalf("launch interactive args = (%q, %q, %v)", name, command, inheritPath)
		}
		return 7, nil
	}
	wslLaunchFunc = func(name, command string, inheritPath bool, stdin, stdout, stderr Handle) (Handle, error) {
		if name != "Arch" || command != "id" || !inheritPath {
			t.Fatalf("launch args = (%q, %q, %v)", name, command, inheritPath)
		}
		if stdin != Handle(1) || stdout != Handle(2) || stderr != Handle(3) {
			t.Fatalf("handles = (%d, %d, %d), want (1,2,3)", stdin, stdout, stderr)
		}
		return Handle(99), nil
	}
	wslGetDistributionConfigurationFunc = func(name string) (uint32, uint64, uint32, error) {
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		return 2, 1000, 8, nil
	}
	wslConfigureDistributionFunc = func(name string, uid uint64, flags uint32) error {
		if name != "Arch" || uid != 1001 || flags != 4 {
			t.Fatalf("configure args = (%q, %d, %d)", name, uid, flags)
		}
		return nil
	}

	lib := nativeWslLib{}
	if !lib.IsDistributionRegistered("Arch") {
		t.Fatal("IsDistributionRegistered = false, want true")
	}
	if err := lib.RegisterDistribution("Arch", "rootfs.tar"); err != nil {
		t.Fatalf("RegisterDistribution returned error: %v", err)
	}
	if err := lib.UnregisterDistribution("Arch"); err != nil {
		t.Fatalf("UnregisterDistribution returned error: %v", err)
	}
	exitCode, err := lib.LaunchInteractive("Arch", "echo hello", true)
	if err != nil {
		t.Fatalf("LaunchInteractive returned error: %v", err)
	}
	if exitCode != 7 {
		t.Fatalf("exitCode = %d, want 7", exitCode)
	}
	handle, err := lib.Launch("Arch", "id", true, Handle(1), Handle(2), Handle(3))
	if err != nil {
		t.Fatalf("Launch returned error: %v", err)
	}
	if handle != Handle(99) {
		t.Fatalf("handle = %d, want 99", handle)
	}
	ver, uid, flags, err := lib.GetDistributionConfiguration("Arch")
	if err != nil {
		t.Fatalf("GetDistributionConfiguration returned error: %v", err)
	}
	if ver != 2 || uid != 1000 || flags != 8 {
		t.Fatalf("config = (%d, %d, %d), want (2, 1000, 8)", ver, uid, flags)
	}
	if err := lib.ConfigureDistribution("Arch", 1001, 4); err != nil {
		t.Fatalf("ConfigureDistribution returned error: %v", err)
	}
}

func TestNativeWslReg_DelegatesAndConvertsProfiles(t *testing.T) {
	withImplHooksLocked(t)

	origGetByName := regGetProfileFromNameFunc
	origGetByPath := regGetProfileFromBasePathFunc
	origWrite := regWriteProfileFunc
	origSetVersion := regSetWslVersionFunc
	origGenerate := regGenerateProfileFunc
	t.Cleanup(func() {
		regGetProfileFromNameFunc = origGetByName
		regGetProfileFromBasePathFunc = origGetByPath
		regWriteProfileFunc = origWrite
		regSetWslVersionFunc = origSetVersion
		regGenerateProfileFunc = origGenerate
	})

	src := wslreg.Profile{
		UUID:              "{1234}",
		BasePath:          "C:\\WSL\\Arch",
		DistributionName:  "Arch",
		DefaultUid:        1000,
		Flags:             8,
		State:             1,
		Version:           2,
		PackageFamilyName: "pkg",
		WsldlTerm:         1,
	}
	regGetProfileFromNameFunc = func(name string) (wslreg.Profile, error) {
		if name != "Arch" {
			t.Fatalf("name = %q, want %q", name, "Arch")
		}
		return src, nil
	}
	regGetProfileFromBasePathFunc = func(path string) (wslreg.Profile, error) {
		if path != "C:\\WSL\\Arch" {
			t.Fatalf("path = %q, want %q", path, "C:\\WSL\\Arch")
		}
		return src, nil
	}
	regWriteProfileFunc = func(profile wslreg.Profile) error {
		if !reflect.DeepEqual(profile, src) {
			t.Fatalf("profile = %+v, want %+v", profile, src)
		}
		return nil
	}
	regSetWslVersionFunc = func(name string, version int) error {
		if name != "Arch" || version != 2 {
			t.Fatalf("set version args = (%q, %d)", name, version)
		}
		return nil
	}
	regGenerateProfileFunc = func() wslreg.Profile {
		return src
	}

	reg := nativeWslReg{}
	p1, err := reg.GetProfileFromName("Arch")
	if err != nil {
		t.Fatalf("GetProfileFromName returned error: %v", err)
	}
	if !reflect.DeepEqual(p1, toProfile(src)) {
		t.Fatalf("profile = %+v, want %+v", p1, toProfile(src))
	}
	p2, err := reg.GetProfileFromBasePath("C:\\WSL\\Arch")
	if err != nil {
		t.Fatalf("GetProfileFromBasePath returned error: %v", err)
	}
	if !reflect.DeepEqual(p2, toProfile(src)) {
		t.Fatalf("profile = %+v, want %+v", p2, toProfile(src))
	}
	if err := reg.WriteProfile(toProfile(src)); err != nil {
		t.Fatalf("WriteProfile returned error: %v", err)
	}
	if err := reg.SetWslVersion("Arch", 2); err != nil {
		t.Fatalf("SetWslVersion returned error: %v", err)
	}
	if got := reg.GenerateProfile(); !reflect.DeepEqual(got, toProfile(src)) {
		t.Fatalf("GenerateProfile = %+v, want %+v", got, toProfile(src))
	}
}

func TestProfileConverters_RoundTrip(t *testing.T) {
	withImplHooksLocked(t)

	src := wslreg.Profile{
		UUID:              "{abcd}",
		BasePath:          "C:\\WSL\\Ubuntu",
		DistributionName:  "Ubuntu",
		DefaultUid:        1001,
		Flags:             2,
		State:             1,
		Version:           2,
		PackageFamilyName: "pkg2",
		WsldlTerm:         2,
	}

	got := toProfile(src)
	back := fromProfile(got)
	if !reflect.DeepEqual(back, src) {
		t.Fatalf("round trip = %+v, want %+v", back, src)
	}
}

func TestNewNativeProviders_ReturnExpectedTypes(t *testing.T) {
	withImplHooksLocked(t)

	if _, ok := NewNativeWslLib().(nativeWslLib); !ok {
		t.Fatalf("NewNativeWslLib type = %T, want nativeWslLib", NewNativeWslLib())
	}
	if _, ok := NewNativeWslReg().(nativeWslReg); !ok {
		t.Fatalf("NewNativeWslReg type = %T, want nativeWslReg", NewNativeWslReg())
	}
}

func TestNativeWslReg_ReturnsWrappedErrors(t *testing.T) {
	withImplHooksLocked(t)

	origGetByName := regGetProfileFromNameFunc
	origGetByPath := regGetProfileFromBasePathFunc
	origWrite := regWriteProfileFunc
	origSetVersion := regSetWslVersionFunc
	t.Cleanup(func() {
		regGetProfileFromNameFunc = origGetByName
		regGetProfileFromBasePathFunc = origGetByPath
		regWriteProfileFunc = origWrite
		regSetWslVersionFunc = origSetVersion
	})

	wantErr := errors.New("boom")
	regGetProfileFromNameFunc = func(name string) (wslreg.Profile, error) { return wslreg.Profile{}, wantErr }
	regGetProfileFromBasePathFunc = func(path string) (wslreg.Profile, error) { return wslreg.Profile{}, wantErr }
	regWriteProfileFunc = func(profile wslreg.Profile) error { return wantErr }
	regSetWslVersionFunc = func(name string, version int) error { return wantErr }

	reg := nativeWslReg{}
	if _, err := reg.GetProfileFromName("Arch"); !errors.Is(err, wantErr) {
		t.Fatalf("GetProfileFromName err = %v, want %v", err, wantErr)
	}
	if _, err := reg.GetProfileFromBasePath("C:\\WSL\\Arch"); !errors.Is(err, wantErr) {
		t.Fatalf("GetProfileFromBasePath err = %v, want %v", err, wantErr)
	}
	if err := reg.WriteProfile(Profile{}); !errors.Is(err, wantErr) {
		t.Fatalf("WriteProfile err = %v, want %v", err, wantErr)
	}
	if err := reg.SetWslVersion("Arch", 2); !errors.Is(err, wantErr) {
		t.Fatalf("SetWslVersion err = %v, want %v", err, wantErr)
	}
}
