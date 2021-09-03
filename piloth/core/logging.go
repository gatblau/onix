package core

import (
	"log"
	"os"
)

/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "PILOT INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	WarningLogger = log.New(os.Stdout, "PILOT WARNING: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	ErrorLogger = log.New(os.Stderr, "PILOT ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	DebugLogger = log.New(os.Stdout, "PILOT DEBUG: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
}
