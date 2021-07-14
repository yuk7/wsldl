package wtutils

import (
	"io/ioutil"
	"os"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/unicode"
	"muzzammil.xyz/jsonc"
)

const (
	// WTPackageName is package name of Windows Terminal
	WTPackageName = "Microsoft.WindowsTerminal_8wekyb3d8bbwe"
	// WTProfileNameSpaceUUID is uuid of Windows Terminal Profile NameSpace
	WTProfileNameSpaceUUID = "2bde4a90-d05f-401c-9492-e40884ead1d8"
)

// Config is json root Config
type Config struct {
	Profiles struct {
		ProfileList []Profile `json:"list"`
	} `json:"profiles"`
}

// Profile is profile of terminal
type Profile struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandline"`
	GUID        string `json:"guid"`
	Source      string `json:"source"`
}

// ReadWTConfigJSON reads Windows Terminal configuration json file
func ReadWTConfigJSON() (res string, err error) {
	json := os.Getenv("LOCALAPPDATA")
	json = json + "\\Packages\\" + WTPackageName + "\\LocalState\\settings.json"

	b, err := ioutil.ReadFile(json)
	if err != nil {
		return
	}
	res = string(b)
	return
}

// ParseWTConfigJSON parses Windows Terminal configuration json string
func ParseWTConfigJSON(str string) (res Config, err error) {
	var c Config
	err = jsonc.Unmarshal([]byte(str), &c)
	if err != nil {
		return
	}
	res = c
	return
}

// ReadParseWTConfig reads and parses Windows Terminal configuration json file
func ReadParseWTConfig() (conf Config, err error) {
	json, err := ReadWTConfigJSON()
	if err != nil {
		return
	}
	conf, err = ParseWTConfigJSON(json)
	return
}

// CreateProfileGUID creates Windows Terminal GUID (based on uuidv5 but utf16le)
func CreateProfileGUID(str string) string {
	uuidNS, _ := uuid.FromString(WTProfileNameSpaceUUID)

	utf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utfEncoder := utf16le.NewEncoder()
	ut16LeEncodedMessage, _ := utfEncoder.String(str)

	uuid := uuid.NewV5(uuidNS, ut16LeEncodedMessage)
	return uuid.String()
}
