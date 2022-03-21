/*
  Onix Config Manager - Artisan Host Runner
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package handlers

import (
	"fmt"
	"net/http"
)

func checkErr(w http.ResponseWriter, msg string, err error) bool {
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", msg, err)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	return err != nil
}
