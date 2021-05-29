package core

/*
  Onix Config Manager - REMote Host Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"fmt"
	"github.com/jackc/pgtype"
	"hash/fnv"
	"strings"
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

func toTextArray(tag string) string {
	parts := strings.Split(tag, ",")
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, part := range parts {
		buf.WriteString("\"")
		buf.WriteString(part)
		buf.WriteString("\"")
		if i < len(parts)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.String()
}
