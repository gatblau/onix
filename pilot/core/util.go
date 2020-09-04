package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/google/renameio"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// executes a command
func execute(cmd string) (string, error) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		//nolint:gosec
		c = exec.Command(strArr[0])
	} else {
		//nolint:gosec
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	// execute the command asynchronously
	if err := c.Start(); err != nil {
		return stderr.String(), fmt.Errorf("executing %s failed: %s", cmd, err)
	}
	done := make(chan error)
	// launch a go routine to wait for the command to execute
	go func() {
		// send a message to the done channel if completed or error
		done <- c.Wait()
	}()
	// wait for the done channel
	select {
	case <-done:
		// command completed
	case <-time.After(6 * time.Second):
		// command timed out after 6 secs
		return stderr.String(), fmt.Errorf("executing '%s' timed out", cmd)
	}
	return stdout.String(), nil
}

// fetch configuration
func (p *Pilot) fetch() (bool, string) {
	// retrieve configuration information
	log.Info().Msgf("fetching configuration for application with key '%s'\n", p.Cfg.EmConf.ItemInstance)
	item, err := p.Ox.GetItem(&oxc.Item{Key: p.Cfg.EmConf.ItemInstance})
	if err != nil {
		log.Warn().Msgf("cannot fetch application configuration: %s\n", err)
		log.Info().Msgf("the application configuration will be unmanaged until it is created in Onix")
	} else {
		log.Info().Msgf("application configuration retrieved successfully\n")
	}
	if item != nil {
		return true, item.Txt
	}
	return false, ""
}

// save configuration to disk
func (p *Pilot) save(cfg string) error {
	log.Info().Msgf("backing up current configuration")
	err := copyFile(p.Cfg.CfgFile, fmt.Sprintf("%s.bak", p.Cfg.CfgFile))
	if err != nil {
		log.Warn().Msgf("cannot backup configuration: %s", err)
	}
	// write retrieved configuration to disk
	if len(cfg) > 0 {
		err = ioutil.WriteFile(p.Cfg.CfgFile, []byte(cfg), 0644)
	} else {
		log.Warn().Msg("cannot write configuration to file, configuration is empty")
	}
	if err != nil {
		log.Error().Msgf("failed to write application configuration file: %s\n", err)
	} else {
		log.Info().Msgf("writing application configuration to '%s'\n", p.Cfg.CfgFile)
	}
	return err
}

// connect to the MQTT broker and subscribe for notifications
func (p *Pilot) subscribe() {
	err := p.EM.Connect()
	if err != nil {
		log.Error().Msgf("failed to connect to the notification broker: %s\n", err)
	} else {
		log.Info().Msgf("connected to notification broker, subscribed to '/II_%s' topic\n", p.Cfg.EmConf.ItemInstance)
	}
}

// instigate an application configuration reload
func (p *Pilot) reload(cfg string) {
	// if a reload command is defined
	if len(p.Cfg.ReloadCmd) > 0 {
		// execute the command
		execute(p.Cfg.ReloadCmd)
	} else
	// if a reload URI is defined
	if len(p.Cfg.ReloadURI) > 0 {
		// post the configuration to the URI
		p.postConfig(cfg)
	} else {
		// not reloading
		log.Info().Msg("skipping reloading")
	}
}

// if application is not ready after configuration reload then restore
// previous configuration
func (p *Pilot) checkRestore() {

}

// post the app configuration to the reload URI
func (p *Pilot) postConfig(cfg string) {
	// gets a reader for the payload
	reader := bytes.NewReader([]byte(cfg))
	// constructs the request
	req, err := http.NewRequest("POST", p.Cfg.ReloadURI, reader)
	if err != nil {
		log.Error().Msgf("failed to create http request for reload URI: %s", err)
	}
	// if a user name is provided then add basic authentication token
	if len(p.Cfg.ReloadURIUser) > 0 {
		req.Header.Add("Authorization", p.basicToken(p.Cfg.ReloadURIUser, p.Cfg.ReloadURIPwd))
	}
	// sets a request timeout
	http.DefaultClient.Timeout = 6 * time.Second
	// issue the request
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("failed to post configuration to reload URI: %s", err)
	} else {
		log.Info().Msgf("application configuration successfully posted to '%s'", p.Cfg.ReloadURI)
	}
}

// creates a new Basic Authentication Token
func (p *Pilot) basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// copy a file from a source to a destination
func copyFile(src, dest string) error {
	srcContent, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcContent.Close()

	data, err := ioutil.ReadAll(srcContent)
	if err != nil {
		return err
	}
	return renameio.WriteFile(dest, data, 0644)
}
