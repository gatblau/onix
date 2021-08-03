package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type UploadCmd struct {
	cmd *cobra.Command
}

func NewUploadCmd() *UploadCmd {
	c := &UploadCmd{
		cmd: &cobra.Command{
			Use:   "upload",
			Short: "uploads file to object storage bucket",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *UploadCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	_, filename := object_storage.ReadObjectStorageData()
	if err != nil {
		log.Println(err)
	}
	targetName, size, status, msg := object_storage.Uploader(useSSL, filename)
	if status == true {
		log.Println("Uploaded", targetName, " of size: ", size, "Successfully")
	} else {
		log.Fatal("Something went wrong. Error: ", msg)
	}
}
