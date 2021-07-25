package merge

// Set a group of values for a variable group identified by a group name for a specific index
type Set struct {
	Context *Context
	// a list of values associated with a set name
	// eg: [ "NAME" ] [ "port a" ]
	//     [ "DESC" ] [ "this is port a" ]
	//     [ "VALUE" ] [ "80" ]
	Value map[string]string
}

func (s *Set) Get(name string) string {
	return s.Value[name]
}
