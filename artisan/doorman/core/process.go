/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"bytes"
	"fmt"
	"github.com/gatblau/onix/artisan/build"
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/doorman/types"
	"github.com/gatblau/onix/artisan/export"
	"github.com/gatblau/onix/artisan/merge"
	"github.com/gatblau/onix/artisan/registry"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Processor struct {
	id   string
	path string
	tmp  string
	out  *bytes.Buffer
	db   Db
	reg  *registry.LocalRegistry
}

func NewProcessor(deploymentId, bucketPath string) Processor {
	p := Processor{}
	// p.bucketURI, p.releaseURI = getURIs(id)
	p.id = deploymentId
	p.path = bucketPath
	p.out = new(bytes.Buffer)
	p.db = *NewDb()
	p.reg = registry.NewLocalRegistry()
	return p
}

func (p Processor) Info(format string, a ...interface{}) {
	format = fmt.Sprintf("%s INFO %s\n", time.Now().Format("02-01-06 15:04:05"), format)
	_, err := p.out.WriteString(fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Printf("cannot log INFO: %s\n", err)
	}
}

func (p Processor) Error(format string, a ...interface{}) error {
	format = fmt.Sprintf("%s ERROR %s", time.Now().Format("02-01-06 15:04:05"), format)
	msg := fmt.Sprintf(format, a...)
	p.out.WriteString(fmt.Sprintf("%s\n", msg))
	return fmt.Errorf(msg)
}

func (p Processor) Warn(format string, a ...interface{}) {
	format = fmt.Sprintf("%s WARN %s\n", time.Now().Format("02-01-06 15:04:05"), format)
	p.out.WriteString(fmt.Sprintf(format, a...))
}

// Start starts processing a pipeline asynchronously
func (p Processor) Start() error {
	return p.process()
}

func (p Processor) process() error {
	p.Info("start processing release Id=%s - Path=%s", p.id, p.path)
	pipes, err := FindPipelinesByBucketId(p.id)
	if err != nil {
		return p.Error("cannot retrieve pipelines for bucket Id='%s': %s\n", p.id, err)
	}
	if len(pipes) == 0 {
		return p.Error("no pipeline configuration found for release Id=%s\n", p.id)
	}
	for _, pipe := range pipes {
		err = p.processPipeline(pipe)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Processor) processPipeline(pipe types.Pipeline) error {
	_ = registry.NewLocalRegistry()
	// create a new temp folder for processing
	tmp, err := core.NewTempDir()
	if err != nil {
		return p.Error("cannot create temporary folder: %s\n", err)
	}
	p.tmp = tmp
	// find inbound route
	for _, inRoute := range pipe.InboundRoutes {
		// find the inbound route matching the bucket Id
		if inRoute.BucketId == p.id {
			// process the inbound route and run any specified commands
			err = p.processInboundRoute(pipe, inRoute)
			if err != nil {
				return err
			}
			// process the outbound route(s)
			for _, outRoute := range pipe.OutboundRoutes {
				err = p.processOutboundRoute(inRoute, outRoute)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func (p Processor) processInboundRoute(pipe types.Pipeline, route types.InRoute) error {
	// download spec
	p.Info("downloading specification: started")
	err := export.DownloadSpec(fmt.Sprintf("%s/%s", route.BucketURI, p.path), fmt.Sprintf("%s:%s", route.User, route.Pwd), p.tmp)
	if err != nil {
		return p.Error("cannot download specification: %s\n", err)
	}
	p.Info("downloading specification: complete")
	// execute commands
	p.Info("verifying downloaded files: started")
	for _, command := range pipe.Commands {
		c := strings.ReplaceAll(command.Value, "${path}", p.tmp)
		p.Info("executing verification task: %s", c)
		out, exeErr := build.ExeAsync(c, ".", merge.NewEnVarFromSlice([]string{}), false)
		if exeErr != nil {
			return p.Error("execution failed: %s", err)
		}
		// use the regex in the command definition to decide if the command execution failed based on the content of the output
		matched, regexErr := regexp.MatchString(command.ErrorRegex, out)
		if regexErr != nil {
			return p.Error("invalid regex %s: %s", command.ErrorRegex, regexErr)
		}
		// if the regex matched return error and content of command output
		if matched {
			return p.Error("command %s failed:\n%s", out)
		}
	}
	p.Info("verifying downloaded files: complete")
	// load public key
	p.Info("loading verification key %s", route.PublicKey)
	pubKey, keyErr := p.db.FindKeyByName(route.PublicKey)
	if keyErr != nil {
		return p.Error("cannot load verification key %s", route.PublicKey)
	}
	// stores the public key in tmp folder
	err = os.WriteFile(p.verifyKeyFile(), []byte(pubKey.Value), 0660)
	if err != nil {
		return p.Error("cannot persist verification key %s to working folder %s", route.PublicKey, p.tmp)
	}
	return nil
}

func (p Processor) processOutboundRoute(inRoute types.InRoute, outRoute types.OutRoute) error {
	var userPwd string
	p.Info("processing outbound route %s: started", outRoute.Name)
	if outRoute.S3Store != nil {
		// import spec
		if err := p.importSpec(); err != nil {
			return err
		}
		// if S3 requires re-signing
		if outRoute.S3Store.Sign {
			// prepare the private key
			p.Info("loading S3 store signing key %s", outRoute.S3Store.PrivateKey)
			privKey, keyErr := p.db.FindKeyByName(outRoute.S3Store.PrivateKey)
			if keyErr != nil {
				return p.Error("cannot load signing key %s", outRoute.S3Store.PrivateKey)
			}
			// stores the public key in tmp folder
			err := os.WriteFile(p.signKeyS3File(), []byte(privKey.Value), 0660)
			if err != nil {
				return p.Error("cannot persist signing key %s to working folder %s", outRoute.PackageRegistry.PrivateKey, p.tmp)
			}
			// resign packages
			p.Info("re-signing packages with key %s: started", outRoute.PackageRegistry.PrivateKey)
			r := registry.NewLocalRegistry()
			err = r.Sign("", p.signKeyS3File(), "")
			if err != nil {
				return p.Error("cannot re-sign spec artefacts with key %s: %s", outRoute.PackageRegistry.PrivateKey, err)
			}
			p.Info("re-signing packages with key %s: completed", outRoute.PackageRegistry.PrivateKey)
		}
		// export packages
		p.Info("exporting re-signed packages: started")
		spec, err := export.NewSpec(p.tmp, "")
		if err != nil {
			return p.Error("cannot load spec.yaml from working folder: %s", err)
		}
		targetURI := fmt.Sprintf("%s/%s", outRoute.S3Store.BucketURI, p.path)
		err = export.ExportSpec(*spec, targetURI, "", outRoute.S3Store.Creds(), "")
		if err != nil {
			return p.Error("cannot export spec to %s: %s", targetURI, err)
		}
		p.Info("exporting re-signed packages: completed")
		userPwd = fmt.Sprintf("%s:%s", outRoute.S3Store.User, outRoute.S3Store.Pwd)
		p.Info("uploading to S3 store %s: started", outRoute.S3Store.BucketURI)
		err = export.UploadSpec(outRoute.S3Store.BucketURI, userPwd, "")
		if err != nil {
			return p.Error("cannot upload spec tarball files to S3 store %s: %s", outRoute.S3Store.BucketURI, err)
		}
		p.Info("uploading to S3 store %s: completed", outRoute.S3Store.BucketURI)
	}
	if outRoute.PackageRegistry != nil {
		if err := p.importSpec(); err != nil {
			return err
		}
		// if resigning of packages is required
		if outRoute.PackageRegistry.Sign {
			// prepare the private key
			p.Info("loading package registry signing key %s", outRoute.PackageRegistry.PrivateKey)
			privKey, keyErr := p.db.FindKeyByName(outRoute.PackageRegistry.PrivateKey)
			if keyErr != nil {
				return p.Error("cannot load signing key %s", outRoute.PackageRegistry.PrivateKey)
			}
			// stores the public key in tmp folder
			err := os.WriteFile(p.signKeyArtFile(), []byte(privKey.Value), 0660)
			if err != nil {
				return p.Error("cannot persist signing key %s to working folder %s", outRoute.PackageRegistry.PrivateKey, p.tmp)
			}
			// resign packages
			p.Info("re-signing packages with key %s: started", outRoute.PackageRegistry.PrivateKey)
			r := registry.NewLocalRegistry()
			err = r.Sign("", p.signKeyArtFile(), "")
			if err != nil {
				return p.Error("cannot re-sign spec artefacts with key %s: %s", outRoute.PackageRegistry.PrivateKey, err)
			}
			p.Info("re-signing packages with key %s: completed", outRoute.PackageRegistry.PrivateKey)
		}
		// tagging artefacts & pushing
		p.Info("tagging and pushing artefacts to package registry: started")
		userPwd = fmt.Sprintf("%s:%s", outRoute.PackageRegistry.User, outRoute.PackageRegistry.Pwd)
		err := export.PushSpec(p.tmp, outRoute.PackageRegistry.Domain, outRoute.PackageRegistry.Group, userPwd, "", false, true, false)
		if err != nil {
			return p.Error("cannot push spec artefacts to package registry: %s", err)
		}
		p.Info("tagging and pushing artefacts to package registry: completed")
	}
	if outRoute.ImageRegistry != nil {
		// tagging images & pushing
		p.Info("tagging and pushing images to docker registry: started")
		userPwd = fmt.Sprintf("%s:%s", outRoute.ImageRegistry.User, outRoute.ImageRegistry.Pwd)
		err := export.PushSpec(p.tmp, outRoute.PackageRegistry.Domain, outRoute.PackageRegistry.Group, userPwd, "", false, true, true)
		if err != nil {
			return p.Error("cannot push spec artefacts to image registry: %s", err)
		}
		p.Info("tagging and pushing artefacts to image registry: completed")

	}
	p.Info("processing outbound route %s: completed", outRoute.Name)
	return nil
}

func (p Processor) verifyKeyFile() string {
	return filepath.Join(p.tmp, "verify_key.pgp")
}

func (p Processor) signKeyS3File() string {
	return filepath.Join(p.tmp, "sign_key_s3.pgp")
}

func (p Processor) signKeyArtFile() string {
	return filepath.Join(p.tmp, "sign_key_art.pgp")
}

func (p Processor) importSpec() error {
	// import artefacts
	p.Info("importing specification files: started")
	// prune local registry
	err := p.reg.Prune()
	if err != nil {
		return p.Error("cannot prune local registry: %s", err)
	}
	// import spec in tmp folder
	err = export.ImportSpec(p.tmp, "", "", p.verifyKeyFile(), false)
	if err != nil {
		return p.Error("cannot import spec: %s", err)
	}
	p.Info("importing specification files: complete")
	return nil
}

// func getURIs(uri string) (bucketURI, releaseURI string) {
//     parts := strings.Split(uri, "/")
//     for i := 0; i < len(parts)-2; i++ {
//         if i == 0 {
//             bucketURI = parts[i]
//         } else {
//             bucketURI = fmt.Sprintf("%s/%s", bucketURI, parts[i])
//         }
//     }
//     return bucketURI, fmt.Sprintf("%s/%s", bucketURI, parts[len(parts)-2])
// }
