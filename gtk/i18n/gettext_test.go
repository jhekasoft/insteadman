package i18n

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestT(t *testing.T) {
	translates := map[string]map[string]string{
		"uk": {
			"About": "Про програму",
			"%s Installing...": "%s Встановлення...",
		},
		"ru": {
			"About": "О программе",
			"%s Installing...": "%s Установка...",
		},
	}

	for lang, langTranslates := range translates {
		Init("../../resources/locale", "insteadman", lang)

		for key, mustBeTranslate := range langTranslates {
			assert.Equal(t, T(key), mustBeTranslate)
		}
	}


}