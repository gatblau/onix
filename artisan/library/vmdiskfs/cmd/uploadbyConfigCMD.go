package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type UploadByConfigCmd struct {
	cmd *cobra.Command
}

func NewUploadByConfigCmd() *UploadByConfigCmd {
	c := &UploadByConfigCmd{
		cmd: &cobra.Command{
			Use:   "upload-with-config",
			Short: "uploads file to object storage bucket using configuration",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *UploadByConfigCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	if err != nil {
		log.Println(err)
	}
	_, filename := object_storage.ReadObjectStorageData()
	diskConfiguration, status := object_storage.ReadConfig(filename)
	if status == true {
		imageCount := len(diskConfiguration.Disk)
		for id := 0; id < imageCount; id++ {
			targetName, size, result, msg := object_storage.Uploader(useSSL, diskConfiguration.Disk[id].Name)
			if result == true {
				log.Println("Uploaded", targetName, " of size: ", size, "Successfully")
			} else {
				log.Fatal("Something went wrong. Error: ", msg)
			}
		}
	}
}
