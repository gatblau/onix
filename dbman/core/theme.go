//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package core

import (
	"fmt"
)

// an HTTP query theme used to style the presentation of the results in a web page
type Theme struct {
	// the content of the CSS stylesheet to embed
	Style *string
	// the content of the stylesheet to embed
	Header *string
	// the content of the footer to embed
	Footer *string
}

// creates a new theme fetching any relevant resources from the remote repository
func NewTheme(name string, sm *ScriptManager) *Theme {
	if len(name) > 0 {
		// note if reading fails assume no file available
		style, _ := sm.FetchFile(fmt.Sprintf("/theme/%s/style.css", name))
		header, _ := sm.FetchFile(fmt.Sprintf("/theme/%s/header.html", name))
		footer, _ := sm.FetchFile(fmt.Sprintf("/theme/%s/footer.html", name))
		return &Theme{
			Style:  style,
			Header: header,
			Footer: footer,
		}
	} else {
		// return an empty
		return new(Theme)
	}
}
