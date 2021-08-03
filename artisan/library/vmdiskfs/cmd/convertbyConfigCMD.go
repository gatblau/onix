package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
)

type ConvertByConfigCmd struct {
	cmd *cobra.Command
}

func NewConvertByConfigCmd() *ConvertByConfigCmd {
	c := &ConvertByConfigCmd{
		cmd: &cobra.Command{
			Use:   "convert-with-config",
			Short: "Converts images to KubeVirt format with configuration provided",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *ConvertByConfigCmd) Run(cmd *cobra.Command, args []string) {
	_, filename := object_storage.ReadObjectStorageData()
	diskConfiguration, status := object_storage.ReadConfig(filename)
	if status == true {
		imageCount := len(diskConfiguration.Disk)
		for id := 0; id < imageCount; id++ {
			result, msg := object_storage.ConvertExecute(diskConfiguration.Disk[id].Name)
			if result == true {
				log.Println("Successfully Converted")
			} else {
				log.Fatal("Something went wrong. Error: ", msg)
			}
		}
	}
}
