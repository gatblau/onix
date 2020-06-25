//   Onix Config DatabaseProvider - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// )
//
// // wraps the psql command
// type Psql struct {
// 	host string
// 	user string
// 	db   string
// }
//
// func (sql *Psql) IsConnectionOK() (bool, error) {
// 	//export PGPASSWORD=onix;psql -h localhost -U onix onix
// 	return sql.run(exec.Command("pgsql", "-h", sql.host, "-U", sql.user, sql.db))
// }
//
// // generic command execution with stdout
// func (sql *Psql) run(cmd *exec.Cmd) (bool, error) {
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Printf("!!!  I cannot execute sql command: %v\n", err)
// 		return false, err
// 	}
// 	return true, nil
// }
