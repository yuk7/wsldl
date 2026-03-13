package wtutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWTConfigJSON(t *testing.T) {
	t.Parallel()

	json := `{
		"profiles": {
			"list": [
				{
					"name": "Arch",
					"commandline": "wsl.exe -d Arch",
					"guid": "{11111111-1111-1111-1111-111111111111}",
					"source": "Windows.Terminal.Wsl"
				}
			]
		}
	}`

	conf, err := ParseWTConfigJSON(json)
	if err != nil {
		t.Fatalf("ParseWTConfigJSON failed: %v", err)
	}
	if len(conf.Profiles.ProfileList) != 1 {
		t.Fatalf("profile count = %d, want 1", len(conf.Profiles.ProfileList))
	}
	if conf.Profiles.ProfileList[0].Name != "Arch" {
		t.Fatalf("profile name = %q, want %q", conf.Profiles.ProfileList[0].Name, "Arch")
	}
}

func TestParseWTConfigJSON_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ParseWTConfigJSON(`{"profiles":`); err == nil {
		t.Fatal("ParseWTConfigJSON succeeded unexpectedly")
	}
}

func TestCreateProfileGUID(t *testing.T) {
	t.Parallel()

	got1 := CreateProfileGUID("Arch")
	got2 := CreateProfileGUID("Arch")
	if got1 != got2 {
		t.Fatalf("CreateProfileGUID is not deterministic: %q != %q", got1, got2)
	}
	if got1 != "a5a97cb8-8961-5535-816d-772efe0c6a3f" {
		t.Fatalf("CreateProfileGUID = %q, want %q", got1, "a5a97cb8-8961-5535-816d-772efe0c6a3f")
	}
}

func TestWTConfigPath(t *testing.T) {
	t.Parallel()

	got := wtConfigPath("C:\\Users\\me\\AppData\\Local")
	want := "C:\\Users\\me\\AppData\\Local\\Packages\\" + WTPackageName + "\\LocalState\\settings.json"
	if got != want {
		t.Fatalf("wtConfigPath = %q, want %q", got, want)
	}
}

func TestReadWTConfigJSONFromPath(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := tmp + "\\settings.json"
	content := `{"profiles":{"list":[]}}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write settings.json failed: %v", err)
	}

	got, err := readWTConfigJSONFromPath(path)
	if err != nil {
		t.Fatalf("readWTConfigJSONFromPath failed: %v", err)
	}
	if got != content {
		t.Fatalf("readWTConfigJSONFromPath = %q, want %q", got, content)
	}
}

func TestReadWTConfigJSONFromPath_NotFound(t *testing.T) {
	t.Parallel()

	if _, err := readWTConfigJSONFromPath(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("readWTConfigJSONFromPath succeeded unexpectedly")
	}
}

func TestReadWTConfigJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("LOCALAPPDATA", tmp)

	content := `{"profiles":{"list":[{"name":"Arch"}]}}`
	path := wtConfigPath(tmp)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write settings.json failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	got, err := ReadWTConfigJSON()
	if err != nil {
		t.Fatalf("ReadWTConfigJSON failed: %v", err)
	}
	if got != content {
		t.Fatalf("ReadWTConfigJSON = %q, want %q", got, content)
	}
}

func TestReadParseWTConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("LOCALAPPDATA", tmp)

	content := `{"profiles":{"list":[{"name":"Arch","commandline":"wsl.exe -d Arch"}]}}`
	path := wtConfigPath(tmp)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write settings.json failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	conf, err := ReadParseWTConfig()
	if err != nil {
		t.Fatalf("ReadParseWTConfig failed: %v", err)
	}
	if len(conf.Profiles.ProfileList) != 1 {
		t.Fatalf("profile count = %d, want 1", len(conf.Profiles.ProfileList))
	}
	if conf.Profiles.ProfileList[0].Name != "Arch" {
		t.Fatalf("profile name = %q, want %q", conf.Profiles.ProfileList[0].Name, "Arch")
	}
}

func TestReadParseWTConfig_ReadError(t *testing.T) {
	t.Setenv("LOCALAPPDATA", t.TempDir())

	if _, err := ReadParseWTConfig(); err == nil {
		t.Fatal("ReadParseWTConfig succeeded unexpectedly")
	}
}
