// +build !windows

package i18n

import "os"

func SetGettextLanguage(language string) {
	os.Setenv("LANGUAGE", language)
}
