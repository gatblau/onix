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

func TestExtract(t *testing.T) {
	content := `
Praesent tristique magna sit amet. 
Etiam tempor orci eu lobortis elementum nibh tellus. 
Eros donec ac odio tempor orci. 
Nulla at volutpat diam ut venenatis tellus in metus. 
Enim ut sem viverra aliquet. 
Consequat mauris nunc congue nisi vitae suscipit. 
Enim ut sem viverra aliquet 234.
Nunc scelerisque viverra mauris in aliquam sem fringilla ut morbi. 
Dui accumsan sit amet nulla facilisi morbi.
`
	// extracts the content between prefix "Enim ut" and end of line "$"
	matches := Extract(content, "Enim ut", "$", -1)
	for _, match := range matches {
		fmt.Println(match)
	}
}
