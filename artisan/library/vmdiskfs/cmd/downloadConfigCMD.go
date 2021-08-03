package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type DownloadConfigCmd struct {
	cmd *cobra.Command
}

func NewDownloadConfigCmd() *DownloadConfigCmd {
	c := &DownloadConfigCmd{
		cmd: &cobra.Command{
			Use:   "download-config",
			Short: "download config",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DownloadConfigCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	if err != nil {
		log.Println(err)
	}
	bucket, filename := object_storage.ReadObjectStorageData()
	size, status, msg := object_storage.Downloader(useSSL, bucket, filename)
	if status == true {
		log.Printf("Configuration downloaded: %v size is %i", filename, size)
	} else {
		log.Printf("Configuration %v not found! Error: %v", filename, msg)
	}
}
