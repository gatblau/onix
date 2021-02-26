/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package cmd

import (
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/i18n"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// list local packages
type LangFetchCmd struct {
	cmd *cobra.Command
}

func NewLangFetchCmd() *LangFetchCmd {
	c := &LangFetchCmd{
		cmd: &cobra.Command{
			Use:   "fetch [language code]",
			Short: "fetches a language dictionary and installs it in the local registry",
			Long:  `fetches a language dictionary and installs it in the local registry`,
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *LangFetchCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		i18n.Raise(i18n.ERR_INSUFFICIENT_ARGS)
	}
	if len(args) > 1 {
		i18n.Raise(i18n.ERR_TOO_MANY_ARGS)
	}
	// checks the lang path exists within the registry
	core.LangExists()
	// try and fetch the language dictionary
	url := fmt.Sprintf("https://raw.githubusercontent.com/gatblau/artisan/master/lang/%s_i18n.toml", args[0])
	resp, err := http.Get(url)
	i18n.Err(err, i18n.ERR_CANT_DOWNLOAD_LANG, url)
	if resp.StatusCode != 200 {
		i18n.Err(fmt.Errorf(resp.Status), i18n.ERR_CANT_DOWNLOAD_LANG, url)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	i18n.Err(err, i18n.ERR_CANT_READ_RESPONSE)
	err = ioutil.WriteFile(path.Join(core.LangPath(), fmt.Sprintf("%s_i18n.toml", args[0])), bodyBytes, os.ModePerm)
	i18n.Err(err, i18n.ERR_CANT_SAVE_FILE)
}
