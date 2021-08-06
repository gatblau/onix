package backend

/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
import (
	"github.com/gatblau/onix/artisan/data"
	"github.com/gatblau/onix/artisan/registry"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"os"
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
	panic("implement me")
}

// GetAllRepositoryInfo get information for all repositories in the remote repository
func (s3 *S3Backend) GetAllRepositoryInfo(user, pwd string) ([]*registry.Repository, error) {
	panic("implement me")
}

// GetRepositoryInfo get information for a specific repository in the remote repository
func (s3 *S3Backend) GetRepositoryInfo(group, name, user, pwd string) (*registry.Repository, error) {
	panic("implement me")
}

// GetPackageInfo get package information
func (s3 *S3Backend) GetPackageInfo(group, name, id, user, pwd string) (*registry.Package, error) {
	panic("implement me")
}

// UpdatePackageInfo update package information
func (s3 *S3Backend) UpdatePackageInfo(group, name string, packageInfo *registry.Package, user string, pwd string) error {
	panic("implement me")
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
