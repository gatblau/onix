package object_storage

import (
	"go/types"
	"log"
	"os"
	"os/exec"
)

func ConvertExecute(filename string) (bool, error) {
	f, extensionOriginal := TargetFileName(filename)
	extension := os.Getenv("OUTPUT_FORMAT")

	source := "source/" + f + "." + extensionOriginal
	destination := "output/" + f + "." + extension

	cmd := exec.Command("qemu-img", "convert", "-f", extensionOriginal, "-O", extension, source, destination)
	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	log.Printf("Waiting for command to finish...")
	if err := cmd.Wait(); err != nil {
		return false, err
	} else {
		return true, types.Error{}
	}
}
