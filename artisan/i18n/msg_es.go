/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package i18n

// spanish
var msg_es = map[I18NKey]string{
	// mensajes de error
	ERR_CANT_CREATE_REGISTRY_FOLDER: "no puedo crear la carpeta del registro local '%s', el inicio del usuario es: '%s'",
	ERR_CANT_EXEC_FUNC_IN_PACKAGE:   "no puedo ejecutar la funcion '%s' en el paquete '%s'",
	ERR_CANT_LOAD_PRIV_KEY:          "no puedo cargar la llave privada",
	ERR_CANT_PUSH_PACKAGE:           "no puedo empujar el paquete",
	ERR_INVALID_PACKAGE_NAME:        "el nombre del paquete es invalido",
	// mensajes de informacion
	INFO_PUSHED:          "empujado: %s\n",
	INFO_NOTHING_TO_PUSH: "nada que empujar\n",
	// labels
	LBL_LS_HEADER: "REPOSITORIO\tETIQUETA\tIDENTIFICADOR DE PAQUETE\tTIPO DE PAQUETE\tCREADO\tTAMANO",
}
