package core

import (
	"fmt"
	"testing"
)

func TestFindFiles(t *testing.T) {
	files, err := FindFiles(".", "^.*\\.(go|art)$")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

func TestPackName(t *testing.T) {
	n, _ := ParseName("localhost%:9009/hh/ff/gg/hh/hh/jj/kk'|&*/testpk:v1")
	fmt.Println(n)
}

func TestUserPwd(t *testing.T) {
	u, p := UserPwd("ab.er:46567785")
	fmt.Println(u, p)
}
