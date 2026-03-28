package i18n

import (
	"os"
	"strings"
)

var currentLang = "en"

func SetLang(lang string) {
	if lang == "zh" || lang == "en" {
		currentLang = lang
	}
}

func GetLang() string {
	return currentLang
}

func T(key string) string {
	var m map[string]string
	if currentLang == "zh" {
		m = zhStrings
	} else {
		m = enStrings
	}
	if v, ok := m[key]; ok {
		return v
	}
	return key
}

func DetectLang(configLang, envLang string) string {
	if configLang == "zh" || configLang == "en" {
		return configLang
	}
	if envLang == "zh" || envLang == "en" {
		return envLang
	}
	loc := os.Getenv("LANG")
	if loc == "" {
		loc = os.Getenv("LC_ALL")
	}
	if strings.HasPrefix(loc, "zh") {
		return "zh"
	}
	return "en"
}
