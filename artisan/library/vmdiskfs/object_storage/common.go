package object_storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/minio/minio-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ObjectStorageProvider returns authentication to object storage
func ObjectStorageProvider(useSSL bool) *minio.Client {
	endpoint := os.Getenv("AWS_URL")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	provider, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Println(err)
	}
	return provider
}

// ReadObjectStorageData information about source bucket name and file name
func ReadObjectStorageData() (string, string) {
	var lines []string
	file, err := os.Open("scripts/parsed.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines[0], lines[1]
}

func TargetFileName(filename string) (string, string) {
	splitStr := strings.Split(filename, ".")
	targetFileName := splitStr[0]
	sourceFormatName := splitStr[1]

	return targetFileName, sourceFormatName
}

// ReadConfig - parses configuration for further processes
func ReadConfig(filename string) (*DiskConfiguration, bool) {
	diskConfiguration := DiskConfiguration{}
	file, err := ioutil.ReadFile("source/" + filename)
	if err != nil {
		log.Println(err)
	}
	if err = yaml.Unmarshal(file, &diskConfiguration); err != nil {
		return &DiskConfiguration{}, false
	} else {
		return &diskConfiguration, true
	}
}

// UnescapeUnicodeCharactersInJSON - use unicode instead of converting symbol & to u0026
func UnescapeUnicodeCharactersInJSON(jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// WebHookMessageGenerator - Generates webhook message
func WebHookMessageGenerator(useSSL bool, diskInfo *DiskConfiguration) *WebHookMessage {
	storageURL := os.Getenv("AWS_URL")
	bucket := os.Getenv("AWS_DEST_BUCKET")
	wHook := WebHookMessage{}
	var wHookDisk []WebHookDisk
	id := len(diskInfo.Disk)
	for i := 0; i < id; i++ {
		targetName, _ := TargetFileName(diskInfo.Disk[i].Name)
		temp := WebHookDisk{}
		temp.Name = targetName + ".qcow2"
		temp.DiskID = diskInfo.Disk[i].DiskID
		temp.BootDisk = diskInfo.Disk[i].BootDisk
		temp.Size = diskInfo.Disk[i].Size
		url, status, err := GenerateUrl(useSSL, bucket, temp.Name)
		furl := url.Scheme + "://" + url.Host + url.Path + "?" + url.RawQuery
		if status == true {
			temp.Url = furl
		} else {
			log.Println("Something went wrong. Error: ", err)
		}
		wHookDisk = append(wHookDisk, temp)
	}
	wHook.StorageUrl = storageURL
	wHook.Cpu = diskInfo.Cpu
	wHook.Memory = diskInfo.Memory
	wHook.Os = diskInfo.Os
	wHook.Arch = diskInfo.Arch
	wHook.Source = diskInfo.Source
	wHook.Disk = wHookDisk
	return &wHook
}

// WebHookMessageSender - Sends webhook to artisan runner webhook
func WebHookMessageSender(webHookMsg *WebHookMessage, url string) (error, string) {
	byteWebHook, err := json.Marshal(webHookMsg)
	jsonUnescaped, _ := UnescapeUnicodeCharactersInJSON(byteWebHook)
	if err != nil {
		return err, "_"
	} else {
		req, err := http.NewRequest("POST", url, bytes.NewBufferString(string(jsonUnescaped)))
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error occured: %v", err)
		}
		defer resp.Body.Close()
		return nil, resp.Status
	}
}
