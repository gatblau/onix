/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0

  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package core

// a distinct piece of application configuration
type appCfg struct {
	// the application configuration
	config string
	// the application configuration metadata (front matter)
	meta *frontMatter
	// the type of configuration
	confType confType
	// the reload trigger type
	reloadTrigger trigger
}
