package preset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePresetJSON(t *testing.T) {
	t.Parallel()

	json := `{
		// jsonc comment
		"wslversion": 2,
		"installfile": "rootfs.tar.gz",
		"installfilesha256": "abc123"
	}`

	got, err := ParsePresetJSON(json)
	if err != nil {
		t.Fatalf("ParsePresetJSON failed: %v", err)
	}

	if got.WslVersion != 2 {
		t.Fatalf("WslVersion = %d, want 2", got.WslVersion)
	}
	if got.InstallFile != "rootfs.tar.gz" {
		t.Fatalf("InstallFile = %q, want %q", got.InstallFile, "rootfs.tar.gz")
	}
	if got.InstallFileSha256 != "abc123" {
		t.Fatalf("InstallFileSha256 = %q, want %q", got.InstallFileSha256, "abc123")
	}
}

func TestParsePresetJSON_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ParsePresetJSON(`{"wslversion":`); err == nil {
		t.Fatal("ParsePresetJSON succeeded unexpectedly")
	}
}

func TestReadPresetJSONFromDir(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	content := `{"wslversion":1}`
	if err := os.WriteFile(filepath.Join(tmp, "preset.json"), []byte(content), 0o600); err != nil {
		t.Fatalf("write preset.json failed: %v", err)
	}

	got, err := readPresetJSONFromDir(tmp)
	if err != nil {
		t.Fatalf("readPresetJSONFromDir failed: %v", err)
	}
	if got != content {
		t.Fatalf("readPresetJSONFromDir = %q, want %q", got, content)
	}
}

func TestReadPresetJSONFromDir_NotFound(t *testing.T) {
	t.Parallel()

	if _, err := readPresetJSONFromDir(t.TempDir()); err == nil {
		t.Fatal("readPresetJSONFromDir succeeded unexpectedly")
	}
}
