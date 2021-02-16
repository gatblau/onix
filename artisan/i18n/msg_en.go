/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package i18n

// english
var msg_en = map[I18NKey]string{
	// error messages
	ERR_CANT_CREATE_REGISTRY_FOLDER: "cannot create local registry folder '%s', user home: '%s'",
	ERR_CANT_DOWNLOAD_LANG:          "cannot download language dictionary from '%s'",
	ERR_CANT_EXEC_FUNC_IN_PACKAGE:   "cannot execute function '%s' in package '%s'",
	ERR_CANT_LOAD_PRIV_KEY:          "cannot load the private key",
	ERR_CANT_PUSH_PACKAGE:           "cannot push package",
	ERR_CANT_READ_RESPONSE:          "cannot read response body",
	ERR_CANT_SAVE_FILE:              "cannot save file",
	ERR_CANNOT_UPDATE_LANG_FILE:     "cannot update language file",
	ERR_INSUFFICIENT_ARGS:           "insufficient arguments",
	ERR_INVALID_PACKAGE_NAME:        "invalid package name",
	ERR_TOO_MANY_ARGS:               "too many arguments",
	// information messages
	INFO_PUSHED:          "pushed: %s\n",
	INFO_NOTHING_TO_PUSH: "nothing to push\n",
	// labels
	LBL_LS_HEADER: "REPOSITORY\tTAG\tPACKAGE ID\tPACKAGE TYPE\tCREATED\tSIZE",
}
