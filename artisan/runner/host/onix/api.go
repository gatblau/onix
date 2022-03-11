package onix

/*
  Onix Pilot Host Control
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"fmt"
	"log"

	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/oxlib/oxc"
)

// API backend services API
type API struct {
	conf *Conf
	ox   *oxc.Client
}

var (
	api *API
)

func Api() *API {
	var (
		err    error
		newapi *API
	)
	if api == nil {
		newapi, err = newAPI(new(Conf))
		if err != nil {
			log.Fatalf("ERROR: fail to create backend services API: %s", err)
		}
		api = newapi
	}
	return api
}

func newAPI(cfg *Conf) (*API, error) {
	oxcfg := &oxc.ClientConf{
		BaseURI:            cfg.getOxWapiUrl(),
		Username:           cfg.getOxWapiUsername(),
		Password:           cfg.getOxWapiPassword(),
		InsecureSkipVerify: cfg.getOxWapiInsecureSkipVerify(),
	}
	oxcfg.SetAuthMode("basic")
	ox, err := oxc.NewClient(oxcfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create onix http client: %s", err)
	}
	return &API{
		conf: cfg,
		ox:   ox,
	}, nil
}

func (r *API) GetCommand(cmdName string) (*Cmd, error) {

	item, err := r.ox.GetItem(&oxc.Item{Key: cmdName})
	if err != nil {
		return nil, fmt.Errorf("Failed to get command with key '%s' from Onix: %s", cmdName, err)
	}
	input, err := getInputFromMap(item.Meta)
	if err != nil {
		return nil, fmt.Errorf("Failed to transform input map: %s", err)
	}
	return &Cmd{
		Key:         item.Key,
		Description: item.Description,
		Package:     item.GetStringAttr("PACKAGE"),
		Function:    item.GetStringAttr("FX"),
		User:        item.GetStringAttr("USER"),
		Pwd:         item.GetStringAttr("PWD"),
		Verbose:     item.GetBoolAttr("VERBOSE"),
		Input:       input,
	}, nil
}

func getInputFromMap(inputMap map[string]interface{}) (*data.Input, error) {
	input := &data.Input{}
	in := inputMap["input"]
	if in != nil {
		// load vars
		vars := in.(map[string]interface{})["var"]
		if vars != nil {
			input.Var = data.Vars{}
			v := vars.([]interface{})
			for _, i := range v {
				varMap := i.(map[string]interface{})
				nameValue, ok := varMap["name"].(string)
				if !ok {
					return nil, fmt.Errorf("variable NAME must be provided, can't process payload\n")
				}
				descValue, ok := varMap["description"].(string)
				if !ok {
					descValue = ""
				}
				typeValue, ok := varMap["type"].(string)
				if !ok || len(typeValue) == 0 {
					typeValue = "string"
				}
				requiredValue, ok := varMap["required"].(bool)
				if !ok {
					requiredValue = false
				}
				valueValue, ok := varMap["value"].(string)
				if !ok && requiredValue {
					return nil, fmt.Errorf("variable %s VALUE not provided, can't process payload\n", nameValue)
				}
				vv := &data.Var{
					Name:        nameValue,
					Description: descValue,
					Value:       valueValue,
					Type:        typeValue,
					Required:    requiredValue,
				}
				input.Var = append(input.Var, vv)
			}
		}
		// load secrets
		secrets := in.(map[string]interface{})["secret"]
		if secrets != nil {
			input.Secret = data.Secrets{}
			v := secrets.([]interface{})
			for _, i := range v {
				varMap := i.(map[string]interface{})
				nameValue, ok := varMap["name"].(string)
				if !ok {
					return nil, fmt.Errorf("secret NAME must be provided, can't process payload\n")
				}
				descValue, ok := varMap["description"].(string)
				if !ok {
					descValue = ""
				}
				requiredValue, ok := varMap["required"].(bool)
				if !ok {
					requiredValue = false
				}
				valueValue, ok := varMap["value"].(string)
				if !ok {
					return nil, fmt.Errorf("secret %s VALUE not provided, can't process payload\n", nameValue)
				}
				vv := &data.Secret{
					Name:        nameValue,
					Description: descValue,
					Value:       valueValue,
					Required:    requiredValue,
				}
				input.Secret = append(input.Secret, vv)
			}
		}
	}
	return input, nil
}
