package i18n

import (
	"../../core/configurator"
	"github.com/gosexy/gettext"
)

const (
	LocaleDir = "resources/locale"
)

func Init(c *configurator.Configurator, domain string) {
	//gettext.SetLocale(gettext.LcAll, "ru_RU.UTF-8")
	gettext.SetLocale(gettext.LcAll, "")
	gettext.BindTextdomain(domain, c.DataResourcePath(LocaleDir))
	gettext.BindTextdomainCodeset(domain, "UTF-8")
	gettext.Textdomain(domain)
}

func T(message string) string {
	return gettext.Gettext(message)
}
