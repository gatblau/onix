package build

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"testing"
)

func Test(t *testing.T) {
	cmd := `printf "hello %s and %s" "my dear" "yyy aas .w.w"`
	parts := breakdown(cmd)
	for _, part := range parts {
		fmt.Println(part)
	}
}

func breakdown(str string) []string {
	// first break down by double quotes
	// strings.Split(str, )
	p := shellwords.NewParser()
	args, err := p.Parse(str)
	if err != nil {
		panic(err)
	}
	return args
}
