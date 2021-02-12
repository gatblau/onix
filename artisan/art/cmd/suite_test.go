package cmd

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/flow"
	"github.com/gatblau/onix/artisan/tkn"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	flowPath := "s2p_bare.yaml"
	buildFilePath := "."
	// loads a bare flow from the path
	flow, err := flow.NewWithEnv(flowPath, buildFilePath, ".env")
	core.CheckErr(err, "cannot load bare flow")
	// survey for required inputs
	err = flow.Merge(false)
	core.CheckErr(err, "cannot merge bare flow")
	// if tekton format is requested
	// gets a tekton transpiler
	builder := tkn.NewBuilder(flow.Flow)
	// transpile the flow
	buf := builder.Create()
	// write to file
	err = ioutil.WriteFile(tknPath(flowPath), buf.Bytes(), os.ModePerm)
	core.CheckErr(err, "cannot write tekton file")
}

var (
	buildFilePath = "/Users/andresalos/dev/artisan/recipe/java-quarkus/project"
	flowPath      = "/Users/andresalos/dev/artisan/recipe/java-quarkus/project/_build/flows/s2p_bare.yaml"
)

func TestEnv(t *testing.T) {
	// loads a bare flow from the path
	f, err := flow.LoadFlow(flowPath)
	core.CheckErr(err, "cannot load bare flow")

	// loads the build.yaml
	var b *data.BuildFile
	// if there is a build file, load it
	if len(buildFilePath) > 0 {
		b, err = data.LoadBuildFile(path.Join(buildFilePath, "build.yaml"))
	}
	input := f.GetInputDefinition(b)

	switch strings.ToLower("env") {
	case "env":
		doEnv(input)
	}
}

func doEnv(i *data.Input) {
	buf := &bytes.Buffer{}
	for _, v := range i.Var {
		buf.WriteString(fmt.Sprintf("# %s\n", v.Description))
		if len(v.Default) > 0 {
			buf.WriteString(fmt.Sprintf("%s=%s\n", v.Name, v.Default))
		} else {
			buf.WriteString(fmt.Sprintf("%s=\n", v.Name))
		}
	}
	if true {
		fmt.Println(buf)
	} else {
		dir := filepath.Dir(flowPath)
		err := ioutil.WriteFile(path.Join(dir, ".env"), buf.Bytes(), os.ModePerm)
		core.CheckErr(err, "cannot write .env file")
	}
}

func Test(t *testing.T) {
	builder := build.NewBuilder()
	builder.Run("test", "../.", false)
}
