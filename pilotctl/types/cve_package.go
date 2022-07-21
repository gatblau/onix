/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

type CvePackage struct {
	HostUUID    string
	CveID       string
	PackageName string
	FixedIn     string
	CvssScore   float64
}
