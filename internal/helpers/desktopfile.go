package helpers

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

// ValidStartupWMClass reports whether value is a usable StartupWMClass entry.
func ValidStartupWMClass(value string) bool {
	value = strings.TrimSpace(value)
	return value != "" && !strings.EqualFold(value, "undefined")
}

// ExtractJSONStringField finds the first string value for field in loosely parsed JSON bytes.
func ExtractJSONStringField(data []byte, field string) string {
	re := regexp.MustCompile(`"` + regexp.QuoteMeta(field) + `"\s*:\s*"([^"\\]*(?:\\.[^"\\]*)*)"`)
	match := re.FindSubmatch(data)
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

// ParseAppRunBinary returns the main binary path declared in an AppRun script.
func ParseAppRunBinary(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "BIN=") {
			continue
		}
		val := strings.TrimSpace(strings.TrimPrefix(line, "BIN="))
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		val = strings.TrimPrefix(val, "${APPDIR}/")
		val = strings.TrimPrefix(val, "$APPDIR/")
		return filepath.Base(val)
	}
	return ""
}

func CheckDesktopFile(desktopfile string) error {
	// Check for presence of required keys and abort otherwise
	d, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, // Do not cripple lines hat contain ";"
		desktopfile)
	PrintError("ini.load", err)
	neededKeys := []string{"Categories", "Name", "Exec", "Type", "Icon"}
	for _, k := range neededKeys {
		if d.Section("Desktop Entry").HasKey(k) == false {
			return errors.New(".desktop file is missing a '" + k + "'= key\n")
		}
	}

	val, _ := d.Section("Desktop Entry").GetKey("Icon")
	iconname := val.String()
	if strings.Contains(iconname, "/") {
		return errors.New("Desktop file contains Icon= entry with a path")
	}

	if strings.HasSuffix(filepath.Base(iconname), ".png") ||
		strings.HasSuffix(filepath.Base(iconname), ".svg") ||
		strings.HasSuffix(filepath.Base(iconname), ".svgz") ||
		strings.HasSuffix(filepath.Base(iconname), ".xpm") {
		return errors.New("Desktop file contains Icon= entry with a suffix, please remove the suffix")
	}

	return nil
}
