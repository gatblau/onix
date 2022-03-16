/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package registry

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

// RemoteRegistry enables admin operations on a remote registry
type RemoteRegistry struct {
	domain string
	user   string
	pwd    string
	api    *Api
}

// NewRemoteRegistry creates an object to manage a remote registry
func NewRemoteRegistry(domain, user, pwd string) (*RemoteRegistry, error) {
	if strings.HasPrefix(domain, "http") {
		return nil, fmt.Errorf("remote registry domain '%s' should not specify protocol scheme", domain)
	}
	if strings.Contains(domain, "/") {
		return nil, fmt.Errorf("remote registry domain '%s' should not contain slashes", domain)
	}
	return &RemoteRegistry{
		domain: domain,
		user:   user,
		pwd:    pwd,
		api:    newGenericAPI(domain),
	}, nil
}

// List all packages in the remote registry
func (r *RemoteRegistry) List(quiet bool) {
	// get a reference to the remote registry
	repos, err, _, _ := r.api.GetAllRepositoryInfo(r.user, r.pwd)
	core.CheckErr(err, "cannot list remote registry packages")
	var w *tabwriter.Writer
	if quiet {
		// get a table writer for the stdout
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
		// repository, tag, package id, created, size
		for _, repo := range repos {
			for _, a := range repo.Packages {
				_, err = fmt.Fprintln(w, fmt.Sprintf("%s", a.Id[0:12]))
				core.CheckErr(err, "failed to write package Id")
			}
		}
	} else {
		// get a table writer for the stdout
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 12, ' ', 0)
		// print the header row
		_, err = fmt.Fprintln(w, i18n.String(i18n.LBL_LS_HEADER))
		core.CheckErr(err, "failed to write table header")
		// repository, tag, package id, created, size
		for _, repo := range repos {
			for _, a := range repo.Packages {
				for _, tag := range a.Tags {
					_, err = fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
						fmt.Sprintf("%s:%s", r.domain, repo.Repository),
						tag,
						a.Id[0:12],
						a.Type,
						toElapsedLabel(a.Created),
						a.Size),
					)
					core.CheckErr(err, "failed to write output")
				}
			}
		}
	}
	err = w.Flush()
	core.CheckErr(err, "failed to flush output")
}

// Remove one or more packages from a remote registry
func (r *RemoteRegistry) Remove(filter string) error {
	repos, err, _, tls := r.api.GetAllRepositoryInfo(r.user, r.pwd)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		for _, p := range repo.Packages {
			var tagCount = len(p.Tags)
			for _, tag := range p.Tags {
				name, err1 := core.ParseName(fmt.Sprintf("%s/%s:%s", r.domain, repo.Repository, tag))
				if err1 != nil {
					return err1
				}
				matched, err2 := regexp.MatchString(filter, name.String())
				if err2 != nil {
					return fmt.Errorf("invalid filter expression '%s': %s", filter, err2)
				}
				if matched {
					// if more than one tag exist, remove the tag
					if tagCount > 1 {
						// get the package metadata
						pInfo, err3 := r.api.GetPackageInfo(name.Group, name.Name, p.Id, r.user, r.pwd, tls)
						if err3 != nil {
							return err3
						}
						// remove the tag
						pInfo.RemoveTag(tag)
						// push the metadata back to the remote
						err3 = r.api.UpsertPackageInfo(name, pInfo, r.user, r.pwd, tls)
						if err3 != nil {
							return err3
						}
					}
					// if we are hitting the last tag
					if tagCount == 1 {
						// remove the package files
						if err = r.api.DeletePackage(name.Group, name.Name, tag, r.user, r.pwd, tls); err != nil {
							return err
						}
						if err = r.api.DeletePackageInfo(name.Group, name.Name, p.Id, r.user, r.pwd, tls); err != nil {
							return err
						}
					}
					tagCount--
				}
			}
		}
	}
	return nil
}
