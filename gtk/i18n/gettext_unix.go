// +build !windows

package i18n

import "os"

// SetGettextLanguage sets language ("uk", "ru") for gettext translates
func SetGettextLanguage(language string) {
	os.Setenv("LANGUAGE", language)
}
