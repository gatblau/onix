/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package i18n

type I18NKey int

const (
	// error messages
	ERR_CANT_EXEC_FUNC_IN_PACKAGE I18NKey = iota
	ERR_CANT_LOAD_PRIV_KEY
	ERR_CANT_PUSH_PACKAGE
	ERR_INVALID_PACKAGE_NAME
	// information messages
	INFO_PUSHED
	INFO_NOTHING_TO_PUSH
)
