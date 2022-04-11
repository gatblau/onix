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
	"github.com/gatblau/onix/oxlib/oxc"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func SaveInfo(f multipart.File) (string, error) {
	tmp := SyncTemp()
	filePath := path.Join(tmp, "sync.xlsx")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		_ = os.RemoveAll(tmp)
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		_ = os.RemoveAll(tmp)
		return "", err
	}
	return filePath, nil
}

// SyncInfo syncs the content of the input spreadsheet file
// compares the logistics information in the spreadsheet and commits any differences to Onix CMDB
func SyncInfo(file string, api *API, dryRun bool) (diff *Diff, err error) {
	file, _ = filepath.Abs(file)
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err2 := f.Close(); err2 != nil {
			fmt.Println(err2)
		}
		_ = os.RemoveAll(filepath.Dir(file))
	}()
	og, err := loadSheet(f, "org-groups")
	if err != nil {
		return nil, err
	}
	or, err := loadSheet(f, "orgs")
	if err != nil {
		return nil, err
	}
	ar, err := loadSheet(f, "areas")
	if err != nil {
		return nil, err
	}
	lo, err := loadSheet(f, "locations")
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
	links, err := loadLinks(f)
	if err != nil {
		return nil, err
	}
	linksDb, err := loadLinksDb(api)
	if err != nil {
		return nil, err
	}

	// take the difference
	d := &Diff{
		OG:    difference(*og, *ogd),
		OR:    difference(*or, *ord),
		AR:    difference(*ar, *ard),
		LO:    difference(*lo, *lod),
		LINKS: diffLinks(*links, *linksDb),
	}

	// check if there are any hosts in locations to be removed
	err = canRemoveLocation(api, d.LO)
	if err != nil {
		return d, err
	}

	// if this is not a dry run then apply the changes
	if !dryRun {
		// apply org group changes
		err = applyItems(d.OG, "U_ORG_GROUP", api)
		if err != nil {
			return d, err
		}
		// apply orgs changes
		err = applyItems(d.OR, "U_ORG", api)
		if err != nil {
			return d, err
		}
		// apply areas changes
		err = applyItems(d.AR, "U_AREA", api)
		if err != nil {
			return d, err
		}
		// apply locations changes
		err = applyItems(d.LO, "U_LOCATION", api)
		if err != nil {
			return d, err
		}
		// apply links
		err = applyLinks(d.LINKS, api)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func applyLinks(links *DiffLinkReport, api *API) error {
	for _, l := range links.Added {
		r, err := api.ox.PutLink(&oxc.Link{
			Key:          l.key(),
			StartItemKey: l.From,
			EndItemKey:   l.To,
			Description:  fmt.Sprintf("link %s", l),
			Type:         "U_RELATIONSHIP",
		})
		if r.Error {
			return fmt.Errorf(r.Message)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func applyItems(d *DiffReport, itemType string, api *API) error {
	if d.Added != nil {
		for _, i := range d.Added.Values {
			r, err := api.ox.PutItem(&oxc.Item{
				Key:         i.Key,
				Name:        i.Name,
				Description: i.Description,
				Txt:         i.Info,
				Type:        itemType,
			})
			if r.Error {
				return fmt.Errorf(r.Message)
			}
			if err != nil {
				return err
			}
		}
	}
	if d.Updated != nil {
		for _, i := range d.Updated.Values {
			r, err := api.ox.PutItem(&oxc.Item{
				Key:         i.Key,
				Name:        i.Name,
				Description: i.Description,
				Txt:         i.Info,
				Type:        itemType,
			})
			if r.Error {
				return fmt.Errorf(r.Message)
			}
			if err != nil {
				return err
			}
		}
	}
	if d.Removed != nil {
		for _, i := range d.Removed.Values {
			r, err := api.ox.DeleteItem(&oxc.Item{Key: i.Key})
			if r.Error {
				return fmt.Errorf(r.Message)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func loadLinksDb(api *API) (*linkList, error) {
	links, err := api.ox.GetLinks()
	if err != nil {
		return nil, err
	}
	list := &linkList{
		Links: []linkInfo{},
	}
	for _, link := range links.Values {
		// org-group -> org link
		if strings.Contains(link.Key, "OG:") && strings.Contains(link.Key, "OR:") && !list.contains(link.StartItemKey, link.EndItemKey) {
			list.Links = append(list.Links, linkInfo{From: link.StartItemKey, To: link.EndItemKey})
		}
		// org-group -> area link
		if strings.Contains(link.Key, "OG:") && strings.Contains(link.Key, "AR:") && !list.contains(link.StartItemKey, link.EndItemKey) {
			list.Links = append(list.Links, linkInfo{From: link.StartItemKey, To: link.EndItemKey})
		}
		// org -> location link
		if strings.Contains(link.Key, "OR:") && strings.Contains(link.Key, "LO:") && !list.contains(link.StartItemKey, link.EndItemKey) {
			list.Links = append(list.Links, linkInfo{From: link.StartItemKey, To: link.EndItemKey})
		}
		// area -> location link
		if strings.Contains(link.Key, "AR:") && strings.Contains(link.Key, "LO:") && !list.contains(link.StartItemKey, link.EndItemKey) {
			list.Links = append(list.Links, linkInfo{From: link.StartItemKey, To: link.EndItemKey})
		}
	}
	return list, nil
}

func loadLinks(f *excelize.File) (*linkList, error) {
	list := &linkList{
		Links: []linkInfo{},
	}
	rows, err := f.GetRows("links")
	if err != nil {
		return nil, err
	}
	for ix, row := range rows {
		// skips the header
		if ix == 0 {
			continue
		}
		for i := 0; i < 3; i++ {
			key := value(row, i)
			if !(strings.HasPrefix(key, "OG:") || strings.HasPrefix(key, "OR:") || strings.HasPrefix(key, "AR:") || strings.HasPrefix(key, "LO:")) {
				return nil, fmt.Errorf("invalid key '%s' in links worksheet: prefix must be one of 'OG:', 'OR':, 'AR:' or 'LO:'", key)
			}
		}

		// based on indices below
		// 0: org-group
		// 1: org
		// 2: area
		// 3: location

		// org-group -> org link (0->1)
		if !list.contains(value(row, 0), value(row, 1)) {
			list.Links = append(list.Links, linkInfo{From: value(row, 0), To: value(row, 1)})
		}
		// org-group -> area link (0->2)
		if !list.contains(value(row, 0), value(row, 2)) {
			list.Links = append(list.Links, linkInfo{From: value(row, 0), To: value(row, 2)})
		}
		// org -> location link (1->3)
		if !list.contains(value(row, 1), value(row, 3)) {
			list.Links = append(list.Links, linkInfo{From: value(row, 1), To: value(row, 3)})
		}
		// area -> location link (2-3)
		if !list.contains(value(row, 2), value(row, 3)) {
			list.Links = append(list.Links, linkInfo{From: value(row, 2), To: value(row, 3)})
		}
	}
	return list, nil
}

// load items from spreadsheet
func loadSheet(f *excelize.File, sheetName string) (*infoList, error) {
	// loads the org group information
	rows, err := f.GetRows(sheetName)
	if err != nil {
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
	uid := uuid.New()
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

func canRemoveLocation(api *API, diff *DiffReport) error {
	if diff.Removed == nil {
		return nil
	}
	hosts, err := api.GetHostsAtLocations(diff.Removed.keys())
	if err != nil {
		return fmt.Errorf("cannot verify hosts in locations: %s", err)
	}
	if len(hosts) > 0 {
		buf := strings.Builder{}
		buf.WriteString("cannot remove location(s):\n")
		for _, host := range hosts {
			buf.WriteString(fmt.Sprintf("%s: host %s is active in location\n", host.HostUUID, host.Location))
		}
		return fmt.Errorf(buf.String())
	}
	return nil
}

// func SyncPathExists() {
//     tmp := SyncPath()
//     // ensure tmp folder exists for temp file operations
//     _, err := os.Stat(tmp)
//     if os.IsNotExist(err) {
//         _ = os.MkdirAll(tmp, os.ModePerm)
//     }
// }
//
// func HomeDir() string {
//     // if PILOTCTL_HOME is defined use it
//     if pilotCtlHome := os.Getenv("PILOTCTL_HOME"); len(pilotCtlHome) > 0 {
//         return pilotCtlHome
//     }
//     usr, _ := user.Current()
//     return usr.HomeDir
// }

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
		if !(strings.HasPrefix(key, "OG:") || strings.HasPrefix(key, "OR:") || strings.HasPrefix(key, "AR:") || strings.HasPrefix(key, "LO:")) {
			return nil, fmt.Errorf("invalid key '%s' in worksheet: prefix must be one of 'OG:', 'OR':, 'AR:' or 'LO:'", key)
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

func (l *infoList) keys() []string {
	var keys []string
	for _, v := range l.Values {
		keys = append(keys, v.Key)
	}
	return keys
}

func value(row []string, ix int) string {
	if len(row) <= ix {
		return ""
	}
	return row[ix]
}

func diffLinks(source, target linkList) *DiffLinkReport {
	var add, remove []linkInfo
	for _, s := range source.Links {
		if !target.contains(s.From, s.To) {
			add = append(add, linkInfo{
				From: s.From,
				To:   s.To,
			})
		}
	}
	for _, t := range target.Links {
		if !source.contains(t.From, t.To) {
			remove = append(remove, linkInfo{
				From: t.From,
				To:   t.To,
			})
		}
	}
	return &DiffLinkReport{
		Added:   add,
		Removed: remove,
	}
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
		if !target.equals(s) && !add.contains(s.Key) {
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

type DiffLinkReport struct {
	Added   []linkInfo `json:"added,omitempty" yaml:"added,omitempty"`
	Removed []linkInfo `json:"removed" yaml:"removed,omitempty"`
}

type Diff struct {
	OG    *DiffReport     `json:"org_groups" yaml:"org_groups"`
	OR    *DiffReport     `json:"orgs" yaml:"orgs"`
	AR    *DiffReport     `json:"areas" yaml:"areas"`
	LO    *DiffReport     `json:"locations" yaml:"locations"`
	LINKS *DiffLinkReport `json:"links" yaml:"links"`
}

type linkInfo struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

func (i linkInfo) key() string {
	return fmt.Sprintf("%s->%s", strings.ToUpper(i.From), strings.ToUpper(i.To))
}

func (i linkInfo) equals(from, to string) bool {
	return i.From == from && i.To == to
}

type linkList struct {
	Links []linkInfo
}

func (l *linkList) contains(from, to string) bool {
	for _, link := range l.Links {
		if link.From == from && link.To == to {
			return true
		}
	}
	return false
}

func (l *linkList) equals(from string, to string) bool {
	for _, i := range l.Links {
		if i.equals(from, to) {
			return true
		}
	}
	return false
}
