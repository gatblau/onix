package core

import (
	"bytes"
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/google/renameio"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

	err := c.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("executing %s failed: %s", cmd, err)
	}
	return stdout.String(), nil
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

// download configuration
func (p *Pilot) fetchCfg() bool {
	retrieved := false
	// retrieve configuration information
	log.Info().Msgf("searching for application with key '%s'\n", p.Cfg.EmConf.ItemInstance)
	item, err := p.Ox.GetItem(&oxc.Item{Key: p.Cfg.EmConf.ItemInstance})
	if err != nil {
		log.Warn().Msgf("cannot fetch application configuration: %s\n", err)
		log.Info().Msgf("the application configuration will be unmanaged until it is created in Onix")
	} else {
		log.Info().Msgf("application configuration found\n")
	}
	// if we managed to get an item (e.g. Onix is up & running)
	if item != nil {
		// backup configuration
		log.Info().Msgf("backing up current configuration")
		err = copyFile(p.Cfg.CfgFile, fmt.Sprintf("%s.bak", p.Cfg.CfgFile))
		if err != nil {
			log.Warn().Msgf("cannot backup configuration: %s", err)
		}
		// write retrieved configuration to disk
		if len(item.Txt) > 0 {
			err = ioutil.WriteFile(p.Cfg.CfgFile, []byte(item.Txt), 0644)
		} else {
			log.Warn().Msg("cannot write configuration to file, configuration is empty")
		}
		if err != nil {
			log.Error().Msgf("failed to write application configuration file: %s\n", err)
		} else {
			log.Info().Msgf("writing application configuration to '%s'\n", p.Cfg.CfgFile)
			retrieved = true
		}
	}
	return retrieved
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
func (p *Pilot) reload() {
	if len(p.Cfg.ReloadCmd) > 0 {
		execute(p.Cfg.ReloadCmd)
	} else if len(p.Cfg.ReloadURI) > 0 {
		log.Warn().Msgf("reloading configuration by calling URI is currently not implemented")
	} else {
		log.Info().Msg("skipping reloading")
	}
}

// if application is not ready after configuration reload then restore
// previous configuration
func (p *Pilot) checkRestore() {

}
