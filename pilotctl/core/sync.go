/*
  Onix Config Manager - Pilot Control
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	t "github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func SaveInfo(f multipart.File) (string, error) {
	tmp := SyncTemp()
	filePath := path.Join(tmp, "sync.xlsx")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	return filePath, nil
}

// SyncInfo syncs the content of the input spreadsheet file
// compares the logistics information in the spreadsheet and commits any differences to Onix CMDB
func SyncInfo(file string, api *API, dryRun bool) (diff *Diff, err error) {
	out := new(strings.Builder)
	file, _ = filepath.Abs(file)
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
		os.RemoveAll(filepath.Dir(file))
	}()
	og, err := loadSheet(f, out, "org-groups")
	if err != nil {
		return nil, err
	}
	or, err := loadSheet(f, out, "orgs")
	if err != nil {
		return nil, err
	}
	ar, err := loadSheet(f, out, "areas")
	if err != nil {
		return nil, err
	}
	lo, err := loadSheet(f, out, "locations")
	if err != nil {
		return nil, err
	}
	ogd, err := loadDb(api, "U_ORG_GROUP")
	if err != nil {
		return nil, err
	}
	ord, err := loadDb(api, "U_ORG")
	if err != nil {
		return nil, err
	}
	ard, err := loadDb(api, "U_AREA")
	if err != nil {
		return nil, err
	}
	lod, err := loadDb(api, "U_LOCATION")
	if err != nil {
		return nil, err
	}
	// org groups
	ogDiff := difference(*og, *ogd)
	if ogDiff.Added != nil {
		out.WriteString(fmt.Sprintf("ORG GROUPS\n"))
		out.WriteString(fmt.Sprintf("To be added:\n"))
		for _, i := range ogDiff.Added.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if ogDiff.Removed != nil {
		out.WriteString(fmt.Sprintf("To be removed:\n"))
		for _, i := range ogDiff.Removed.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if ogDiff.Updated != nil {
		out.WriteString(fmt.Sprintf("To be updated:\n"))
		for _, i := range ogDiff.Updated.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}

	// orgs
	orDiff := difference(*or, *ord)
	if orDiff.Added != nil {
		out.WriteString(fmt.Sprintf("ORGS\n"))
		out.WriteString(fmt.Sprintf("To be added:\n"))
		for _, i := range orDiff.Added.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if orDiff.Removed != nil {
		out.WriteString(fmt.Sprintf("To be removed:\n"))
		for _, i := range orDiff.Removed.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if orDiff.Updated != nil {
		out.WriteString(fmt.Sprintf("To be updated:\n"))
		for _, i := range orDiff.Updated.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}

	// areas
	arDiff := difference(*ar, *ard)
	if arDiff.Added != nil {
		out.WriteString(fmt.Sprintf("AREAS\n"))
		out.WriteString(fmt.Sprintf("To be added:\n"))
		for _, i := range arDiff.Added.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if arDiff.Updated != nil {
		out.WriteString(fmt.Sprintf("To be removed:\n"))
		for _, i := range arDiff.Updated.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if arDiff.Updated != nil {
		out.WriteString(fmt.Sprintf("To be updated:\n"))
		for _, i := range arDiff.Updated.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}

	// locations
	loDiff := difference(*lo, *lod)
	if loDiff.Added != nil {
		out.WriteString(fmt.Sprintf("LOCATIONS\n"))
		out.WriteString(fmt.Sprintf("To be added:\n"))
		for _, i := range loDiff.Added.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if loDiff.Removed != nil {
		out.WriteString(fmt.Sprintf("To be removed:\n"))
		for _, i := range loDiff.Removed.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	if loDiff.Updated != nil {
		out.WriteString(fmt.Sprintf("To be updated:\n"))
		for _, i := range loDiff.Updated.Values {
			out.WriteString(fmt.Sprintf("- %s -> %s\n", i.Key, i.Name))
		}
	}
	return &Diff{
		OG: ogDiff,
		OR: orDiff,
		AR: arDiff,
		LO: loDiff,
	}, nil
}

// load items from spreadsheet
func loadSheet(f *excelize.File, out *strings.Builder, sheetName string) (*infoList, error) {
	// loads the org group information
	rows, err := f.GetRows(sheetName)
	if err != nil {
		out.WriteString(fmt.Sprintf("ERROR: %s\n", err))
		return nil, err
	}
	return newInfo(rows)
}

// load items from cmdb
func loadDb(api *API, itemType string) (*infoList, error) {
	list := new(infoList)
	items, err := api.ox.GetItemsByType(itemType)
	if err != nil {
		return nil, err
	}
	for _, item := range items.Values {
		list.Values = append(list.Values, info{
			Key:         item.Key,
			Name:        item.Name,
			Description: item.Description,
			Info:        "",
		})
	}
	return list, nil
}

func SyncTemp() string {
	uid := t.New()
	folder := strings.Replace(uid.String(), "-", "", -1)[:12]
	tempDirPath := filepath.Join(SyncPath(), folder)
	err := os.MkdirAll(tempDirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return tempDirPath
}

func SyncPath() string {
	return path.Join(core.HomeDir(), "sync")
}

func SyncPathExists() {
	tmp := SyncPath()
	// ensure tmp folder exists for temp file operations
	_, err := os.Stat(tmp)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(tmp, os.ModePerm)
	}
}

func HomeDir() string {
	// if PILOTCTL_HOME is defined use it
	if pilotCtlHome := os.Getenv("PILOTCTL_HOME"); len(pilotCtlHome) > 0 {
		return pilotCtlHome
	}
	usr, _ := user.Current()
	return usr.HomeDir
}

func Key(iType, name string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s:%s", strings.ToUpper(iType), sha)
}

type info struct {
	Key         string `json:"key" yaml:"key"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Info        string `json:"info,omitempty" yaml:"info,omitempty"`
}

func (i *info) equals(o info) bool {
	return i.Key == o.Key && i.Name == o.Name && i.Description == o.Description && i.Info == o.Info
}

func newInfo(rows [][]string) (*infoList, error) {
	list := new(infoList)
	for ix, row := range rows {
		key := value(row, 0)
		name := value(row, 1)

		// skips the header
		if ix == 0 {
			// only skips it if the first row contains a Key header cell
			if strings.ToUpper(key) == "KEY" {
				continue
			}
		}

		// validates key/name values
		if len(key) == 0 {
			return nil, fmt.Errorf("missing key value, invalid spreadsheet format")
		}
		if len(name) == 0 {
			return nil, fmt.Errorf("missing name value, invalid spreadsheet format")
		}
		list.Values = append(list.Values, info{
			Key:         key,
			Name:        name,
			Description: value(row, 2),
			Info:        value(row, 3),
		})
	}
	return list, nil
}

type infoList struct {
	Values []info `json:"values,omitempty" yaml:"values,omitempty"`
}

func (l *infoList) contains(key string) bool {
	for _, val := range l.Values {
		if val.Key == key {
			return true
		}
	}
	return false
}

func (l *infoList) empty() bool {
	return len(l.Values) == 0
}

func (l *infoList) equals(item info) bool {
	for _, i := range l.Values {
		if i.equals(item) {
			return true
		}
	}
	return false
}

func value(row []string, ix int) string {
	if len(row) <= ix {
		return ""
	}
	return row[ix]
}

func difference(source, target infoList) *DiffReport {
	// add: are the ones in source but not in target
	add := new(infoList)
	for _, s := range source.Values {
		if !target.contains(s.Key) {
			add.Values = append(add.Values, s)
		}
	}
	// remove: are the ones in target but not in source
	remove := new(infoList)
	for _, t := range target.Values {
		if !source.contains(t.Key) {
			remove.Values = append(remove.Values, t)
		}
	}
	// update: are the ones in source and target that are different
	update := new(infoList)
	for _, s := range source.Values {
		if target.equals(s) {
			update.Values = append(update.Values, s)
		}
	}

	return &DiffReport{
		Added:   nullOnEmpty(add),
		Removed: nullOnEmpty(remove),
		Updated: nullOnEmpty(update),
	}
}

func nullOnEmpty(list *infoList) *infoList {
	if list.empty() {
		return nil
	}
	return list
}

type DiffReport struct {
	Added   *infoList `json:"added,omitempty" yaml:"added,omitempty"`
	Removed *infoList `json:"removed" yaml:"removed,omitempty"`
	Updated *infoList `json:"updated" yaml:"updated,omitempty"`
}

type Diff struct {
	OG *DiffReport `json:"org_groups" yaml:"org_groups"`
	OR *DiffReport `json:"orgs" yaml:"orgs"`
	AR *DiffReport `json:"areas" yaml:"areas"`
	LO *DiffReport `json:"locations" yaml:"locations"`
}
