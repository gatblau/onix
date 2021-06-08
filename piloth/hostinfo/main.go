package main

/*
  Onix Config Manager - Host Info Utility
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	hostUtil "github.com/shirou/gopsutil/host"
)

func main() {
	i, err := hostUtil.Info()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", i)
}
