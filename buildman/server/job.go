package server

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os/exec"
	"strings"
)

// check that base images have not changed
type CheckImageJob struct {
	cfg    *policyConfig
	digest *digestCache
	k8s    *K8S
}

func NewCheckImageJob() (*CheckImageJob, error) {
	conf, err := NewPolicyConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot create job: %s", err)
	}
	dig, err := NewDigests()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve cached digests: %s", err)
	}
	err = dig.load()
	if err != nil {
		return nil, fmt.Errorf("cannot load cached digests: %s", err)
	}
	k8s, err := NewK8S()
	if err != nil {
		return nil, fmt.Errorf("cannot create K8S client: %s", err)
	}
	return &CheckImageJob{
		cfg:    conf,
		digest: dig,
		k8s:    k8s,
	}, nil
}

func (c *CheckImageJob) Execute() {
	for _, policy := range c.cfg.Policies {
		if policy.PollBase {
			log.Printf("info: executing policy: %s\n", policy.Name)
			info, err := getImgInfo(policy.Base, policy.User, policy.Pwd)
			if err != nil {
				log.Printf("error: cannot get image infromation for %s\n%sskipping policy\n", policy.Base, err)
				continue
			}
			// get the base image digest
			digest := info.Config.Digest
			// compare the digest with the last recorded one
			if c.digest.changed(policy.Base, digest) {
				log.Printf("info: base image change detected: %s\n", policy.Base)
				log.Printf("info: launching build\n")
				err = c.k8s.NewImagePipeline(policy.Name, policy.Namespace)
				if err != nil {
					log.Printf("error: cannot start image build: %s\n", err)
				}
			} else {
				log.Printf("info: base image unchanged, nothing to do: %s\n", policy.Name)
			}
		}
	}
	c.digest.save()
}

func (c *CheckImageJob) Description() string {
	return "check for changes in container images and triggers image builds"
}

func (c *CheckImageJob) Key() int {
	return hashCode(c.Description())
}

// gets the remote image information
func getImgInfo(imageName, user, pwd string) (*ImgInfo, error) {
	var command *exec.Cmd
	if len(user) > 0 && len(pwd) > 0 {
		command = exec.Command("skopeo", "inspect", fmt.Sprintf("--creds=%s:%s", user, pwd), "--raw", fmt.Sprintf("docker://%s", imageName))
	} else {
		command = exec.Command("skopeo", "inspect", "--raw", fmt.Sprintf("docker://%s", imageName))
	}
	result, err := command.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf(string(e.Stderr))
		}
		return nil, err
	}
	txtResult := strings.TrimRight(string(result), "\n")
	info := new(ImgInfo)
	err = json.Unmarshal([]byte(txtResult), info)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal image manifest: %s", err)
	}
	return info, nil
}

func hashCode(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32())
}
