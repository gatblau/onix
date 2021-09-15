package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/data"
	"github.com/jackc/pgtype"
	"hash/fnv"
)

func hashCode(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32())
}

func toCSV(v pgtype.TextArray) string {
	str := bytes.Buffer{}
	for i, s := range v.Elements {
		str.WriteString(s.String)
		if i < len(v.Elements)-1 {
			str.WriteString(",")
		}
	}
	return str.String()
}

// toTime converts microseconds into HH:mm:SS.ms
func toTime(microseconds int64) string {
	milliseconds := (microseconds / 1000) % 1000
	seconds := (((microseconds / 1000) - milliseconds) / 1000) % 60
	minutes := (((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) % 60
	hours := ((((((microseconds / 1000) - milliseconds) / 1000) - seconds) / 60) - minutes) / 60
	return fmt.Sprintf("%02v:%02v:%02v.%03v", hours, minutes, seconds, milliseconds)
}

func basicAuthToken(user, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// getInputFromMap transform an input in map format to input struct format
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
					return nil, fmt.Errorf("variable NAME not provided, can't process payload\n")
				}
				descValue, ok := varMap["description"].(string)
				if !ok {
					return nil, fmt.Errorf("variable DESCRIPTION not provided, can't process payload\n")
				}
				typeValue, ok := varMap["type"].(string)
				if !ok || len(typeValue) == 0 {
					typeValue = "string"
				}
				valueValue, ok := varMap["value"].(string)
				if !ok {
					return nil, fmt.Errorf("variable VALUE not provided, can't process payload\n")
				}
				vv := &data.Var{
					Name:        nameValue,
					Description: descValue,
					Value:       valueValue,
					Type:        typeValue,
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
				vv := &data.Secret{
					Name:        varMap["name"].(string),
					Description: varMap["description"].(string),
					Value:       varMap["value"].(string),
				}
				input.Secret = append(input.Secret, vv)
			}
		}
	}
	return input, nil
}
