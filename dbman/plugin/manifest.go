//   Onix Config Manager - Dbman
//   Copyright (c) 2018-Present by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.

package plugin

import (
	"bytes"
	"encoding/json"
	"github.com/gatblau/onix/oxlib/oxc"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

// Manifest a database manifest containing the meta data required by DbMan
// to execute commands and queries
type Manifest struct {
	// the database release version
	DbVersion string `json:"dbVersion"`
	// the release description
	Description string `json:"description,omitempty"`
	// the path to where the command scripts are (if not specified use the root of the release)
	CommandsPath string `json:"commandsPath,omitempty"`
	// the path to where the query scripts are (if not specified use the root of the release)
	QueriesPath string `json:"queriesPath,omitempty"`
	// the database provider to use
	DbProvider string `json:"dbProvider"`
	// the list of commands available to execute
	Commands []Command `json:"commands"`
	// the list of commands required to create the database in the first place
	Create Action `json:"create"`
	// the list of commands required to deploy the database objects on an empty database
	Deploy Action `json:"deploy"`
	// the list of commands required to upgrade an existing database
	Upgrade Upgrade `json:"upgrade"`
	// the list of queries available to execute
	Queries []Query `json:"queries"`
}

// Action a database action containing either other sub-actions or commands
type Action struct {
	// the description for the command
	Description string `json:"description"`
	// the list of actions that comprise the command
	Actions []string `json:"actions,omitempty"`
	// the list of sub commands that comprise this command (if any)
	Commands []string `json:"commands,omitempty"`
}

// Command a set of scripts that must be executed within the same database connection
type Command struct {
	// the command identifiable name
	Name string `json:"name"`
	// the description for the action
	Description string `json:"description"`
	// whether to run this action within a database transaction
	Transactional bool `json:"transactional"`
	// whether to connect to the database as an Admin to execute this action
	AsAdmin bool `json:"asAdmin"`
	// whether to connect to the database being managed or simply connect to the server with no specific database
	UseDb bool `json:"useDb"`
	// the list of database scripts that will be executed as part of this action
	Scripts []Script `json:"scripts"`
}

// NewCommand creates a new command from a serialised json string
func NewCommand(jsonString string) (*Command, error) {
	c := &Command{}
	err := json.Unmarshal([]byte(jsonString), c)
	return c, err
}

func (c *Command) ToString() string {
	bytes, e := json.Marshal(c)
	if e != nil {
		return ""
	}
	return string(bytes)
}

// a database script and zero or more merge variables
type Script struct {
	// the script identifiable name
	Name string `json:"name"`
	// the script file name in the git repository
	File string `json:"file"`
	// a list of variables to be merged with the script prior to execution
	Vars []Var `json:"vars"`
	// the content of the script file
	// note: it is internal and automatically populated at runtime from the git repository
	Content string `json:"content,omitempty"`
}

func (c *Script) All() map[string]interface{} {
	m := map[string]interface{}{}
	m["name"] = c.Name
	m["file"] = c.File
	m["content"] = c.Content
	return m
}

// a merge variable for a script
type Var struct {
	// the name of the merge variable use as a placeholder for merging within the script
	Name        string `json:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// the name of the variable to be merged from DbMan's current configuration set
	// note: not used if omitted
	FromConf string `json:"fromConf,omitempty" yaml:"fromConf,omitempty"`
	// the value of the variable, if it is to be merged directly
	// note: not used if omitted
	FromValue string `json:"fromValue,omitempty" yaml:"fromValue,omitempty"`
	// the name of the variable to be merged from the run context
	// available values are dbVersion, appVersion, description
	// note: this is primarily intended for updating the version tracking table
	// not used if omitted
	FromContext string `json:"fromContext,omitempty" yaml:"fromContext,omitempty"`
	// the name of the input parameter
	// allows to pass query parameters via command line or query string
	FromInput string `json:"fromInput,omitempty" yaml:"fromInput,omitempty"`
}

func NewVersion(jsonString string) (*Version, error) {
	v := &Version{}
	err := json.Unmarshal([]byte(jsonString), v)
	return v, err
}

// carries version information
type Version struct {
	// the application version
	AppVersion string `json:"appVersion"`
	// the database version
	DbVersion string `json:"dbVersion"`
	// the name of the query to retrieve the version history
	Description string `json:"description"`
	// the name of the query to retrieve the version history
	Source string `json:"source"`
	// the time the version was released
	Time time.Time `json:"time"`
}

func (v *Version) ToString() string {
	bytes, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return string(bytes)
}

// Query a database query
type Query struct {
	// the identifiable name for the query
	Name string `json:"name"`
	// the description for the query
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// the name of the script file to be executed by the query
	File string `json:"file,omitempty" yaml:"file,omitempty"`
	// a list of variables to merge with the query
	Vars []Var `json:"vars,omitempty" yaml:"vars,omitempty"`
	// the content of the script file
	// note: it is internal and automatically populated at runtime from the git repository
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
}

// NewQuery creates a new query from a serialised json string
func NewQuery(jsonString string) (*Query, error) {
	q := &Query{}
	err := json.Unmarshal([]byte(jsonString), q)
	return q, err
}

func (q *Query) ToString() string {
	bytes, e := json.Marshal(q)
	if e != nil {
		return ""
	}
	return string(bytes)
}

// the commands to run at different stages in an upgrade
type Upgrade struct {
	Description string `json:"description"`
	Prepare     string `json:"prepare"`
	Alter       string `json:"alter"`
	Deploy      string `json:"deploy"`
}

// get a JSON bytes reader for the Plan
func (m *Manifest) json() (*bytes.Reader, error) {
	jsonBytes, err := m.bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(*jsonBytes), err
}

// get a []byte representing the Plan
func (m *Manifest) bytes() (*[]byte, error) {
	b, err := oxc.ToJson(m)
	return &b, err
}

// Decode get the Plan in the http Response
func (m *Manifest) Decode(content []byte) (*Manifest, error) {
	result := new(Manifest)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(result)
	return result, err
}

func (m *Manifest) getCommand(cmdName string) *Command {
	for _, cmd := range m.Commands {
		if cmdName == cmd.Name {
			return &cmd
		}
	}
	return nil
}

func (m *Manifest) findCommands(action *Action) ([]Command, error) {
	var commands []Command
	for _, cmdName := range action.Commands {
		for _, cmd := range m.Commands {
			if cmd.Name == cmdName {
				commands = append(commands, cmd)
			}
		}
	}
	return commands, nil
}

// GetQuery find the query by name
func (m *Manifest) GetQuery(queryName string) *Query {
	for _, query := range m.Queries {
		if query.Name == queryName {
			return &query
		}
	}
	return nil
}

// GetQueriesInfo get a string containing query information in the manifest
func (m *Manifest) GetQueriesInfo(format string, verbose bool) string {
	// make a copy
	queries := make([]Query, len(m.Queries))
	if verbose {
		queries = m.Queries
	} else {
		// clean unneeded info
		for ix, query := range m.Queries {
			queries[ix] = query
			queries[ix].Content = ""
			queries[ix].File = ""
			queries[ix].Vars = nil
		}
	}
	switch strings.ToLower(format) {
	case "json":
		{
			bytes, err := json.Marshal(queries)
			if err != nil {
				return "!!! I cannot convert result into JSON\n"
			}
			return string(bytes)
		}
	default:
		{
			bytes, err := yaml.Marshal(queries)
			if err != nil {
				return "!!! I cannot convert result into YAML\n"
			}
			return string(bytes)
		}
	}
}

func (m *Manifest) GetCommands(commandNames []string) []Command {
	result := make([]Command, 0)
	for _, cmdName := range commandNames {
		for _, command := range m.Commands {
			if command.Name == cmdName {
				result = append(result, command)
			}
		}
	}
	return result
}
