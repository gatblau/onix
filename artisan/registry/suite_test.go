package registry

import (
	"testing"
	"time"
)

func TestRepoDiff(t *testing.T) {
	source := &Repository{
		Repository: "group/name",
		Packages: []*Package{
			{
				Id:      "1",
				Type:    "file",
				FileRef: "123456",
				Tags:    []string{"v1", "v2"},
				Size:    "343B",
				Created: time.Now().UTC().String(),
			},
			{
				Id:      "2",
				Type:    "file",
				FileRef: "789010",
				Tags:    []string{"v3", "v4"},
				Size:    "350B",
				Created: time.Now().UTC().String(),
			},
		},
	}
	target := &Repository{
		Repository: "group/name",
		Packages: []*Package{
			{
				Id:      "1",
				Type:    "file",
				FileRef: "123456",
				Tags:    []string{"v1", "v2", "v3"},
				Size:    "343B",
				Created: time.Now().UTC().String(),
			},
		},
	}

	diff, err := source.Diff(target)
	if err != nil {
		t.Fatal(err)
	}
	if len(diff.Added) != 1 {
		t.Fatal("added list is wrong")
	}
	diff, err = target.Diff(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(diff.Removed) != 1 {
		t.Fatal("removed list is wrong")
	}
}

func TestRemoteList(t *testing.T) {
	r, err := NewRemoteRegistry("localhost:8082", "admin", "admin", "")
	if err != nil {
		t.Errorf(err.Error())
	}
	r.List(false)
}

func TestRemoteRemove(t *testing.T) {
	r, err := NewRemoteRegistry("localhost:8082", "admin", "admin", "")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = r.RemoveByNameFilter("test/testpk:*", false)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestLocalRegistry_Tag(t *testing.T) {
	l := NewLocalRegistry("")
	l.Tag("2fa75", "localhost:8082/test/my-pack:v1")
}
