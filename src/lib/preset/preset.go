package preset

import (
	"os"
	"path/filepath"

	"github.com/muhammadmuzzammil1998/jsonc"
	"github.com/yuk7/wsldl/lib/errutil"
)

var executablePathFunc = errutil.MustExecutable

// ReadPresetJSON reads preset.json configuration json file
func ReadPresetJSON() (res string, err error) {
	return readPresetJSONFromExecutablePath(executablePathFunc())
}

func readPresetJSONFromExecutablePath(executablePath string) (res string, err error) {
	dir := filepath.Dir(executablePath)
	return readPresetJSONFromDir(dir)
}

func readPresetJSONFromDir(dir string) (res string, err error) {
	jsonPath := filepath.Join(dir, "preset.json")
	b, err := os.ReadFile(jsonPath)
	if err != nil {
		return
	}
	res = string(b)
	return
}

// ParsePresetJSON parses preset.json configuration json string
func ParsePresetJSON(str string) (res Preset, err error) {
	var c Preset
	err = jsonc.Unmarshal([]byte(str), &c)
	if err != nil {
		return
	}
	res = c
	return
}

// ReadParsePreset reads and parses preset.json configuration file
func ReadParsePreset() (conf Preset, err error) {
	return readParsePresetFromExecutablePath(executablePathFunc())
}

func readParsePresetFromExecutablePath(executablePath string) (conf Preset, err error) {
	json, err := readPresetJSONFromExecutablePath(executablePath)
	if err != nil {
		return
	}
	conf, err = ParsePresetJSON(json)
	return
}
