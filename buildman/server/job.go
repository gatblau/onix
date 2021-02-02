package server

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os/exec"
	"strings"
	"time"
)

// check that base images have not changed
type CheckImageJob struct {
	cfg *policyConfig
	k8s *K8S
}

func NewCheckImageJob() (*CheckImageJob, error) {
	conf, err := NewPolicyConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot create job: %s", err)
	}
	k8s, err := NewK8S()
	if err != nil {
		return nil, fmt.Errorf("cannot create K8S client: %s", err)
	}
	return &CheckImageJob{
		cfg: conf,
		k8s: k8s,
	}, nil
}

func (c *CheckImageJob) Execute() {
	for _, policy := range c.cfg.Policies {
		if policy.PollBase {
			log.Printf("info: executing policy: %s\n", policy.Name)
			appImgBuildDate, appImgBaseBuildDate, baseImgBuildDate, baseImgBaseBuildDate, err := getImgProps(policy)
			if err != nil {
				log.Printf("error: cannot get image information for %s\n%s\nskipping policy\n", policy.Base, err)
				continue
			}
			// if the app image is not there
			// if the base image creation date label happened after the application image creation date label, or
			// if the base image creation date happened after the time the application image was created
			if appImgBuildDate == nil || baseImgBaseBuildDate.After(*appImgBaseBuildDate) || baseImgBuildDate.After(*appImgBuildDate) {
				log.Printf("info: base image change detected: %s\n", policy.Base)
				log.Printf("info: launching build\n")
				err = c.k8s.NewImagePipeline(policy.Name, policy.Namespace)
				if err != nil {
					log.Printf("error: cannot start image build: %s\n", err)
				} else {
					// if the start of the build was successful then
				}
			} else {
				log.Printf("info: base image unchanged, nothing to do: %s\n", policy.Name)
			}
		}
	}
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
		command = exec.Command("skopeo", "inspect", fmt.Sprintf("--creds=%s:%s", user, pwd), fmt.Sprintf("docker://%s", imageName))
	} else {
		command = exec.Command("skopeo", "inspect", fmt.Sprintf("docker://%s", imageName))
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

// returns all required image dates to work out if new build is required
func getImgProps(policy *policyConf) (appImgBuildDate, appImgBaseBuildDate, baseImgBuildDate, baseImgBaseBuildDate *time.Time, err error) {
	// first retrieves application image information
	appImgInfo, err := getImgInfo(policy.App, policy.AppUser, policy.AppPwd)
	if err != nil {
		// assume app image is not there - this could also happen if login credentials are wrong so it needs improved logic
		// do nothing, let appImgInfo be nil and deal with it down the line
		// return nil, nil, nil, nil, fmt.Errorf("cannot retrieve application image information: %s", err)
	}
	// second retrieves application base image information
	baseImgInfo, err := getImgInfo(policy.Base, policy.BaseUser, policy.BasePwd)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot retrieve application base image information: %s", err)
	}
	if appImgInfo != nil {
		appImgBuildDate = parseTime(appImgInfo.Created)
		if appImgBuildDate == nil {
			return nil, nil, nil, nil, fmt.Errorf("cannot parse created date on app image manifest: '%s'", appImgInfo.Created)
		}
		appImgBaseCreatedOn := appImgInfo.Labels[policy.BaseCreated]
		if len(appImgBaseCreatedOn) == 0 {
			return nil, nil, nil, nil, fmt.Errorf("cannot find base image build date based on label '%s' in image: %s", policy.BaseCreated, policy.App)
		}
		appImgBaseBuildDate, err = getLabelDate(policy.BaseCreated, appImgInfo)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("cannot parse app base image build date based on label '%s': %s", policy.BaseCreated, appImgBaseBuildDate)
		}
	}
	baseImgBuildDate = parseTime(baseImgInfo.Created)
	if baseImgBuildDate == nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot parse created date on base image manifest: '%s'", baseImgInfo.Created)
	}
	baseImgBaseBuildDate, err = getLabelDate(policy.BaseCreated, baseImgInfo)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("cannot parse base image build date based on label '%s': %s", policy.BaseCreated, baseImgBaseBuildDate)
	}
	return appImgBuildDate, appImgBaseBuildDate, baseImgBuildDate, baseImgBaseBuildDate, nil
}

func getLabelDate(label string, info *ImgInfo) (*time.Time, error) {
	value := info.Labels[label]
	if len(value) == 0 {
		return nil, fmt.Errorf("no label %s found in image %s", label, info.Name)
	}
	timeValue := parseTime(value)
	if timeValue == nil {
		return nil, fmt.Errorf("cannot parse time %s in label %s", value, label)
	}
	return timeValue, nil
}

// parses a time in string format trying different formatting
func parseTime(timeString string) *time.Time {
	var result time.Time
	result, err := time.Parse(time.ANSIC, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC822, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC822Z, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC850, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC1123, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC1123Z, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC3339, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(time.RFC3339Nano, timeString)
	if err == nil {
		return &result
	}
	result, err = time.Parse(Rfc3339Custom, timeString)
	if err == nil {
		return &result
	}
	return nil
}

func hashCode(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32())
}

// custom format in docker image
const Rfc3339Custom = "2006-01-02T15:04:05.999999"
