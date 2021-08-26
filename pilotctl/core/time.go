package core

/*
  Onix Pilot Host Control Service
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"math"
	"time"
)

func toElapsedValues(rfc850time string) (int, string, error) {
	if len(rfc850time) == 0 {
		return 0, "", nil
	}
	created, err := time.Parse(time.RFC850, rfc850time)
	if err != nil {
		return 0, "", err
	}
	elapsed := time.Since(created)
	seconds := elapsed.Seconds()
	minutes := elapsed.Minutes()
	hours := elapsed.Hours()
	days := hours / 24
	weeks := days / 7
	months := weeks / 4
	years := months / 12

	if math.Trunc(years) > 0 {
		return int(years), "y", nil
	} else if math.Trunc(months) > 0 {
		return int(months), "M", nil
	} else if math.Trunc(weeks) > 0 {
		return int(weeks), "w", nil
	} else if math.Trunc(days) > 0 {
		return int(days), "d", nil
	} else if math.Trunc(hours) > 0 {
		return int(hours), "H", nil
	} else if math.Trunc(minutes) > 0 {
		return int(minutes), "m", nil
	}
	return int(seconds), "s", nil
}
