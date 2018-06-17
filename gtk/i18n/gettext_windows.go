// +build windows

package i18n

// #include <stdlib.h>
import "C"

import "os"

// SetGettextLanguage sets language ("uk", "ru") for gettext translates
func SetGettextLanguage(language string) {
	os.Setenv("LANGUAGE", language)

	// We should use putenv() because os.Setenv hasn't effect
	cstr := C.CString("LANGUAGE=" + language)
	C.putenv((*C.char)(cstr))
}
