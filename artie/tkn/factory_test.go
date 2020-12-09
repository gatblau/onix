package tkn

import (
	"fmt"
	"testing"
)

func TestTaskYaml(t *testing.T) {
	task := NewTask("quay.io/gatblau/art-buildah")
	y := task.ToYaml()
	fmt.Print(y)
}
