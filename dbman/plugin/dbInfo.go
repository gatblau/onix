//  Onix Config Manager - Dbman
//  Copyright (c) 2018-2020 by www.gatblau.org
//  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//  Contributors to this project, hereby assign copyright in this code to the project,
//  to be licensed under the same terms as the rest of the code.
package plugin

import "encoding/json"

// database server information
type DbInfo struct {
	Database        string
	OperatingSystem string
	Compiler        string
	ProcessorBits   string
}

func NewDbInfoFromJSON(jsonString string) (*DbInfo, error) {
	info := &DbInfo{}
	err := json.Unmarshal([]byte(jsonString), info)
	return info, err
}

func NewDbInfoFromMap(m map[string]interface{}) (*DbInfo, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	info := &DbInfo{}
	err = json.Unmarshal([]byte(j), info)
	return info, err
}

func (info *DbInfo) ToString() string {
	bytes, e := json.Marshal(info)
	if e != nil {
		panic(e)
	}
	return string(bytes)
}
