package i18n

import (
	"github.com/gosexy/gettext"
)

// Init initializes gettext with domain and language
func Init(localeDir, domain, language string) {
	if language != "" {
		SetGettextLanguage(language)
	}

	gettext.SetLocale(gettext.LcAll, "")
	gettext.BindTextdomain(domain, localeDir)
	gettext.BindTextdomainCodeset(domain, "UTF-8")
	gettext.Textdomain(domain)
}

// T translates message and returns translated string
func T(message string) string {
	return gettext.Gettext(message)
}
