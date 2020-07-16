package plugin

import (
	"fmt"
	"github.com/hashicorp/go-plugin"
)

// the decorator wraps the DatabasePlugin interface and exposes it as a DatabaseProvider interface
// the DatabaseProvider is the underlying interface used by the net/rpc protocol to communicate with the plugin
// whereas the DatabasePlugin interface is a friendlier version used by plugin writers
type DatabasePluginDecorator struct {
	Plugin DatabasePlugin
}

func (db *DatabasePluginDecorator) Setup(config string) string {
	output := NewParameter()
	// parse the configuration
	c, err := NewConf(config)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	// allocate the parsed object to cfg
	db.Plugin.Setup(c)
	// return the output
	return output.ToString()
}

// RPC serialisation wrapper for getting database version information
func (db *DatabasePluginDecorator) GetVersion() string {
	// create the output struct
	output := NewParameter()
	// call the plugin operation
	version, err := db.Plugin.GetVersion()
	// if an error is found
	if err != nil {
		// set the result value to an empty map
		output.Set("result", make(map[string]interface{}))
		// return the error
		return output.ToError(err)
	}
	// set the result value
	output.Set("result", version)
	// return the serialised output back to the RPC client
	return output.ToString()
}

func (db *DatabasePluginDecorator) RunCommand(command string) string {
	output := NewParameter()
	cmd, err := NewCommand(command)
	if err != nil {
		return output.ToError(err)
	}
	log, err := db.Plugin.RunCommand(cmd)
	if log.Len() > 0 {
		output.Log(log.String())
	}
	if err != nil {
		return output.ToError(err)
	}
	return output.ToString()
}

func (db *DatabasePluginDecorator) RunQuery(queryInfo string) string {
	output := NewParameter()
	query, err := NewQuery(queryInfo)
	if err != nil {
		return output.ToError(err)
	}
	result, err := db.Plugin.RunQuery(query)
	if err != nil {
		return output.ToError(err)
	}
	output.Set("result", result)
	return output.ToString()
}

func (db *DatabasePluginDecorator) SetVersion(versionInfo string) string {
	output := NewParameter()
	v, err := NewVersion(versionInfo)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	err = db.Plugin.SetVersion(v)
	if err != nil {
		output.SetError(err)
		return output.ToError(err)
	}
	return v.ToString()
}

func (db *DatabasePluginDecorator) GetInfo() string {
	// create the output struct
	output := NewParameter()
	// call the plugin operation
	info, err := db.Plugin.GetInfo()
	// if an error is found
	if err != nil {
		// set the result value to an empty map
		output.Set("result", make(map[string]interface{}))
		// return the error
		return output.ToError(err)
	}
	// set the result value
	output.Set("result", info)
	// return the serialised output back to the RPC client
	return output.ToString()
}

// launch the database plugin
func ServeDbPlugin(pluginName string, impl DatabasePlugin) {
	// launch the plugin as an rpc server
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "dbman-db-provider",
			MagicCookieValue: fmt.Sprintf("dbman-db-%s", pluginName),
		},
		Plugins: map[string]plugin.Plugin{
			pluginName: &DatabaseProviderPlugin{
				Impl: &DatabasePluginDecorator{
					Plugin: impl,
				},
			},
		},
	})
}
