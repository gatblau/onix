package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type DownloadCmd struct {
	cmd *cobra.Command
}

func NewDownloadCmd() *DownloadCmd {
	c := &DownloadCmd{
		cmd: &cobra.Command{
			Use:   "download",
			Short: "downloads file from object storage and bucket to specific folder",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *DownloadCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	if err != nil {
		log.Println(err)
	}
	bucket, filename := object_storage.ReadObjectStorageData()
	diskSize, status, msg := object_storage.Downloader(useSSL, bucket, filename)
	if status == true {
		log.Println("Downloaded", filename, " of size: ", diskSize, "Successfully")
	} else {
		log.Fatal("Something went wrong. Error: ", msg)
	}
}
