/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package i18n

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"os"
	"strings"
)

// prints a localised message
func Printf(key I18NKey, a ...interface{}) {
	fmt.Printf(msg(key), a...)
}

// checks for the  error and if it exists prints a localised error
func Err(err error, key I18NKey, a ...interface{}) {
	if err != nil {
		fmt.Printf("%s - %s\n", fmt.Sprintf(msg(key), a...), err)
		os.Exit(1)
	}
}

// return a localised message for the current user
func msg(key I18NKey) string {
	// get the key in the user language map
	value, found := getMap()[key]
	// if the key is not found or it is empty
	if !found || len(value) == 0 {
		// if in debug mode, show if a language specific key is missing
		core.Debug("key '%d' not found in language '%s'", key, lang())
		// get the value from the english map
		value = msg_en[key]
	}
	return value
}

func getMap() map[I18NKey]string {
	var language string
	// check if internationalisation is enabled
	inter := os.Getenv("ARTISAN_I18N")
	// if not then use english
	if len(inter) == 0 {
		core.Debug("i18n is disabled, to enable set ARTISAN_I18N")
		return msg_en
	}
	// check if an overriding language has been set
	overrideLang := os.Getenv("ARTISAN_LANG")
	// if not
	if len(overrideLang) == 0 {
		// use the current user language
		language = lang()
	} else {
		// use the overriding language
		language = overrideLang
	}
	// return the map for the user language
	switch language {
	case "es":
		return msg_es
	case "de":
		return msg_de
	case "hi":
		return msg_hi
	case "zh":
		return msg_zh
	case "fr":
		return msg_fr
	default:
		// defaults to english
		return msg_en
	}
}

func splitLocale(locale string) (language string, territory string) {
	formattedLocale := strings.Split(locale, ".")[0]
	formattedLocale = strings.Replace(formattedLocale, "-", "_", -1)

	pieces := strings.Split(formattedLocale, "_")
	language = pieces[0]
	territory = ""
	if len(pieces) > 1 {
		territory = strings.Split(formattedLocale, "_")[1]
	}
	return language, territory
}
