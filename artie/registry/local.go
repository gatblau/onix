/*
  Onix Config Manager - Artie
  Copyright (c) 2018-2020 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package registry

import "github.com/gatblau/onix/artie/core"

// the interface implemented by a local registry
// creds = user[:pwd]
type Local interface {
	Push(name *core.ArtieName, remote Remote, creds string)
	Pull(name *core.ArtieName, remote Remote)
	Remove(names []core.ArtieName)
}
