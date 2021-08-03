package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
)

type ConvertCmd struct {
	cmd *cobra.Command
}

func NewConvertCmd() *ConvertCmd {
	c := &ConvertCmd{
		cmd: &cobra.Command{
			Use:   "convert",
			Short: "Converts images to KubeVirt format",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *ConvertCmd) Run(cmd *cobra.Command, args []string) {
	_, filename := object_storage.ReadObjectStorageData()
	status, msg := object_storage.ConvertExecute(filename)
	if status == false {
		log.Fatal("Command finished with error: ", msg)
	} else {
		log.Println("Successfully Converted")
	}
}
