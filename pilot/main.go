/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package main

import "gatblau.org/onix/pilot/cmd"

func main() {
	// if len(os.Args) < 2 {
	//   return
	// }
	// var (
	//   cmdStr = os.Args[1]
	//   vars   string
	// )
	// for i := 2; i < len(os.Args); i++ {
	//   vars += os.Args[i] + " "
	// }
	// child := exec.Command(cmdStr, strings.Trim(vars, " "))
	// child.Stdout = os.Stdout
	// child.Stderr = os.Stderr
	// err := child.Start()
	//
	// if err != nil {
	//   log.Fatal(err)
	// }
	//
	// child.Wait()

	rootCmd := cmd.InitialiseRootCmd()

	// Execute adds all child commands to the root command and sets flags appropriately.
	rootCmd.Execute()
}
