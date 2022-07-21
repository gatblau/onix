/*
Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package types

import (
	"encoding/json"
	"github.com/future-architect/vuls/models"
	"golang.org/x/xerrors"
	"sort"
)

type Cve struct {
	Id               string
	Family           string
	Platform         string
	Summary          string
	AffectedPackages []models.PackageFixStatus
	CVSSScore        float64
	CVSSType         string
	CVSSVector       string
	CVSSSeverity     string
	Mitigations      []string
	PrimarySrc       []string
	PatchURLs        []string
	CPE              []string
	Confidence       []string
	References       []models.References
}

func (c *Cve) Fixed() bool {
	for _, affectedPackage := range c.AffectedPackages {
		if affectedPackage.NotFixedYet {
			return false
		}
	}
	return true
}

type CveReport struct {
	Cves []*Cve
}

func NewCveReport(file []byte) (*CveReport, error) {
	r, err := load(file)
	if err != nil {
		return nil, err
	}
	cveList := make([]*Cve, 0)
	for _, vuln := range r.ScannedCves.ToSortedSlice() {
		cve := new(Cve)
		cve.Id = vuln.CveID
		cve.Family = r.Family
		cve.Platform = r.Platform.Name
		for _, cvss := range vuln.Cvss3Scores() {
			cve.CVSSScore = cvss.Value.Score
			cve.CVSSType = string(cvss.Type)
			cve.CVSSVector = cvss.Value.Vector
			cve.CVSSSeverity = cvss.Value.Severity
			break
		}
		cve.Summary = vuln.Summaries(r.Lang, r.Family)[0].Value
		for _, m := range vuln.Mitigations {
			cve.Mitigations = append(cve.Mitigations, m.URL)
		}
		for _, m := range vuln.CveContents.References(r.Family) {
			cve.References = append(cve.References, m.Value)
		}
		links := vuln.CveContents.PrimarySrcURLs(r.Lang, r.Family, vuln.CveID, vuln.Confidences)
		for _, link := range links {
			cve.PrimarySrc = append(cve.PrimarySrc, link.Value)
		}
		for _, url := range vuln.CveContents.PatchURLs() {
			cve.PatchURLs = append(cve.PatchURLs, url)
		}
		vuln.AffectedPackages.Sort()
		for _, affected := range vuln.AffectedPackages {
			cve.AffectedPackages = append(cve.AffectedPackages, affected)
		}
		sort.Strings(vuln.CpeURIs)
		for _, name := range vuln.CpeURIs {
			cve.CPE = append(cve.CPE, name)
		}
		// for _, l := range vuln.LibraryFixedIns {
		//     libs := r.LibraryScanners.Find(l.Path, l.Name)
		//     for path, lib := range libs {
		//         cve.FixedIns = append(cve.FixedIns, []string{l.Key,
		//             fmt.Sprintf("%s-%s, FixedIn: %s (%s)",
		//                 lib.Name, lib.Version, l.FixedIn, path)})
		//     }
		// }
		for _, confidence := range vuln.Confidences {
			cve.Confidence = append(cve.Confidence, confidence.String())
		}
		cveList = append(cveList, cve)
	}
	return &CveReport{Cves: cveList}, nil
}

func (r *CveReport) Critical() int {
	return r.countBySeverity(9.0, 10.0)
}

func (r *CveReport) High() int {
	return r.countBySeverity(7.0, 8.9)
}

func (r *CveReport) Medium() int {
	return r.countBySeverity(4.0, 6.9)
}

func (r *CveReport) Low() int {
	return r.countBySeverity(0.1, 3.9)
}

func (r *CveReport) countBySeverity(from, to float64) int {
	count := 0
	for _, cve := range r.Cves {
		if cve.CVSSScore >= from && cve.CVSSScore <= to {
			count++
		}
	}
	return count
}

func load(jsonFile []byte) (*models.ScanResult, error) {
	var (
		err error
	)
	result := &models.ScanResult{}
	if err = json.Unmarshal(jsonFile, result); err != nil {
		return nil, xerrors.Errorf("Failed to parse %s: %w", jsonFile, err)
	}
	return result, nil
}
