// +build windows

package i18n

// #include <stdlib.h>
import "C"

import "os"

func SetGettextLanguage(language string) {
	os.Setenv("LANGUAGE", language)

	cstr := C.CString("LANGUAGE=" + language)
	C.putenv((*C.char)(cstr))
}
