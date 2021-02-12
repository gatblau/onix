// +build windows

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.

  Portions taken from: jibber_jabber - Copyright 2014 Pivotal
  https://github.com/cloudfoundry-attic/jibber_jabber
*/
package i18n

import (
	"strings"
	"syscall"
	"unsafe"
)

const LOCALE_NAME_MAX_LENGTH uint32 = 85

func lang() string {
	language, _ := splitLocale(getLocale())
	return strings.ToLower(language)
}

func getLocale() (locale string) {
	locale, err = getWindowsLocale()
	if err != nil {
		locale = ""
	}
	return locale
}

func getWindowsLocaleFrom(sysCall string) (locale string, err error) {
	buffer := make([]uint16, LOCALE_NAME_MAX_LENGTH)

	dll := syscall.MustLoadDLL("kernel32")
	proc := dll.MustFindProc(sysCall)
	r, _, dllError := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(LOCALE_NAME_MAX_LENGTH))
	if r == 0 {
		err = errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
		return
	}

	locale = syscall.UTF16ToString(buffer)

	return
}

func getAllWindowsLocaleFrom(sysCall string) (string, error) {
	dll, err := syscall.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("Could not find kernel32 dll")
	}

	proc, err := dll.FindProc(sysCall)
	if err != nil {
		return "", err
	}

	locale, _, dllError := proc.Call()
	if locale == 0 {
		return "", errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
	}

	return SUPPORTED_LOCALES[locale], nil
}

func getWindowsLocale() (locale string, err error) {
	dll, err := syscall.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("Could not find kernel32 dll")
	}
	proc, err := dll.FindProc("GetVersion")
	if err != nil {
		return "", err
	}
	v, _, _ := proc.Call()
	windowsVersion := byte(v)
	isVistaOrGreater := (windowsVersion >= 6)

	if isVistaOrGreater {
		locale, err = getWindowsLocaleFrom("GetUserDefaultLocaleName")
		if err != nil {
			locale, err = getWindowsLocaleFrom("GetSystemDefaultLocaleName")
		}
	} else if !isVistaOrGreater {
		locale, err = getAllWindowsLocaleFrom("GetUserDefaultLCID")
		if err != nil {
			locale, err = getAllWindowsLocaleFrom("GetSystemDefaultLCID")
		}
	} else {
		panic(v)
	}
	return
}
