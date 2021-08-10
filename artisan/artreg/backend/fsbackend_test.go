package backend

import "testing"

var fs Backend = NewFsBackend()

func Test_FS_PrintUsage(t *testing.T) {
	name := fs.Name()
	if len(name) == 0 {
		t.Errorf("backend name is not implemented")
	}
}

func Test_FS_GetRepositoryInfo(t *testing.T) {
	fs.GetRepositoryInfo("", "", "", "")
}
