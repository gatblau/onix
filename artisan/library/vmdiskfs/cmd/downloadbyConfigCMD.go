package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type DownloadByConfigCmd struct {
	cmd *cobra.Command
}

func NewDownloadByConfigCmd() *DownloadByConfigCmd {
	c := &DownloadByConfigCmd{
		cmd: &cobra.Command{
			Use:   "download-with-config",
			Short: "downloads file from object storage and bucket to specific folder",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DownloadByConfigCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	if err != nil {
		log.Println(err)
	}
	bucket, filename := object_storage.ReadObjectStorageData()
	diskConfiguration, status := object_storage.ReadConfig(filename)
	if status == true {
		imageCount := len(diskConfiguration.Disk)
		for id := 0; id < imageCount; id++ {
			diskSize, result, msg := object_storage.Downloader(useSSL, bucket, diskConfiguration.Disk[id].Name)
			if result == true {
				log.Println("Downloaded", diskConfiguration.Disk[id].Name, " of size: ", diskSize, "Successfully")
			} else {
				log.Fatal("Something went wrong. Error: ", msg)
			}
		}
	}
}
