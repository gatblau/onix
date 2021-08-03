package cmd

import (
	"github.com/gatblau/onix/artisan/library/vmdiskfs/object_storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type SendWebHookCmd struct {
	cmd *cobra.Command
}

func NewSendWebHookCmd() *SendWebHookCmd {
	c := &SendWebHookCmd{
		cmd: &cobra.Command{
			Use:   "send-webhook",
			Short: "Send webhook to Webhook receiver",
		},
	}
	c.cmd.Run = c.Run
	return c
}

func (c *SendWebHookCmd) Run(cmd *cobra.Command, args []string) {
	envSTR := os.Getenv("AWS_USE_SSL")
	useSSL, err := strconv.ParseBool(envSTR)
	if err != nil {
		log.Println(err)
	}
	webhookRCV := os.Getenv("WEBHOOK_RECEIVER")
	_, filename := object_storage.ReadObjectStorageData()
	diskConfiguration, status := object_storage.ReadConfig(filename)
	if status == true {
		msg := object_storage.WebHookMessageGenerator(useSSL, diskConfiguration)
		if failureErr, respStatus := object_storage.WebHookMessageSender(msg, webhookRCV); failureErr != nil {
			log.Println("Something went wrong. Error: ", failureErr)
		} else {
			log.Println(respStatus)
		}
	}
}
