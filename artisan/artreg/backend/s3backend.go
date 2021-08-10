package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/hashicorp/go-uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
)

// S3Backend package registry S3 backend implementation
type S3Backend struct {
	client *minio.Client
}

// NewS3Backend create a new backend
func NewS3Backend(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *S3Backend {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	s3 := &S3Backend{
		client: minioClient,
	}
	return s3
}

// UploadPackage upload an package to the remote repository
func (s3 *S3Backend) UploadPackage(group, name, packageRef string, zipfile multipart.File, jsonFile multipart.File, repo multipart.File, user string, pwd string) error {
	// ensure files are properly closed
	defer zipfile.Close()
	defer jsonFile.Close()
	defer repo.Close()

	seal := new(data.Seal)
	sealBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("cannot read package seal file: %s", err)
	}
	err = json.Unmarshal(sealBytes, seal)
	if err != nil {
		return fmt.Errorf("cannot unmarshal package seal file: %s", err)
	}
	packageBytes, err := ioutil.ReadAll(zipfile)
	if err != nil {
		return fmt.Errorf("cannot read package file: %s", err)
	}
	err = s3.put(sealBytes, group, name, fmt.Sprintf("%s.json", packageRef))
	if err != nil {
		return fmt.Errorf("cannot upload seal file to s3 bucket: %s", err)
	}
	err = s3.put(packageBytes, group, name, fmt.Sprintf("%s.zip", packageRef))
	if err != nil {
		return fmt.Errorf("cannot upload package file to s3 bucket: %s", err)
	}
	return nil
}

// GetAllRepositoryInfo get information for all repositories in the remote repository
func (s3 *S3Backend) GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error) {
	panic("implement me")
}

// GetRepositoryInfo get information for a specific repository in the remote repository
func (s3 *S3Backend) GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error) {
	// try and get the repository.json file from the package bucket
	bn := s3.bucketName(group, name)
	object, err := s3.client.GetObject(context.Background(), bn, "repository.json", minio.GetObjectOptions{})
	// if we have an error
	if err != nil {
		// assumes file not found and returns an empty repository
		return &registry.Repository{
			Repository: fmt.Sprintf("%s/%s", group, name),
			Packages:   make([]*registry.Package, 0),
		}, nil
	}
	// create a bytes buffer to capture the content of the object
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, object); err != nil {
		return nil, fmt.Errorf("cannot copy repository.json bytes to buffer: %s", err)
	}
	repo := new(registry.Repository)
	err = json.Unmarshal(buf.Bytes(), repo)
	return repo, err
}

// GetPackageInfo get package information
func (s3 *S3Backend) GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error) {
	repo, err := s3.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return nil, err
	}
	if repo != nil {
		return repo.FindPackage(id), nil
	}
	return nil, nil
}

// UpdatePackageInfo update package information
func (s3 *S3Backend) UpdatePackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error {
	// get the repository info
	repo, err := s3.GetRepositoryInfo(group, name, user, pwd)
	if err != nil {
		return err
	}
	// update the repository
	updated := repo.UpdatePackage(packageInfo)
	if !updated {
		return fmt.Errorf("package not found in remote repository, no update was made\n")
	}
	// marshal the repository object
	repoBytes, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("cannot marshal repository object: %s\n", err)
	}
	// push bytes to S3
	return s3.put(repoBytes, group, name, "repository.json")
}

// GetManifest get the package manifest
func (s3 *S3Backend) GetManifest(group, name, tag, user, pwd string) (*data.Manifest, error) {
	panic("implement me")
}

// Download open a file for download
func (s3 *S3Backend) Download(repoGroup, repoName, fileName, user, pwd string) (*os.File, error) {
	panic("implement me")
}

// Name print usage info
func (s3 *S3Backend) Name() string {
	panic("implement me")
}

func (s3 *S3Backend) bucketName(repoGroup, repoName string) string {
	return fmt.Sprintf("%s_%s", repoGroup, repoName)
}

// toFile converts the bytes into a file
func (s3 *S3Backend) toFile(bytes []byte) (*os.File, error) {
	// create an UUId
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	// generate an internal random and transient name based on the UUId
	name := path.Join(core.TmpPath(), fmt.Sprintf("%s.file", uuid))
	// create a transient temp file
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	// write the bytes into it
	_, err = file.Write(bytes)
	if err != nil {
		return nil, err
	}
	// closes the file
	file.Close()
	// open the created file
	file, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	// remove the file from the file system
	err = os.Remove(name)
	return file, nil
}

// upload the passed in bytes as a file in the remote s3 repository
func (s3 *S3Backend) put(bytes []byte, repoGroup, repoName, filename string) error {
	file, err := s3.toFile(bytes)
	if err != nil {
		return err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat %s file: %s\n", filename, err)
	}
	_, err = s3.client.PutObject(context.Background(),
		s3.bucketName(repoGroup, repoName),
		filename,
		file,
		fileStat.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("cannot put object: %s\n", err)
	}
	return nil
}
