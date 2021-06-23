package core

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	token := newToken("ABCDEFG")
	fmt.Printf("%s\n\n", token)
	host, ok, _ := readToken(token)
	fmt.Printf("host-> %s\n\n", host)
	fmt.Printf("ok-> %v\n\n", ok)
}
