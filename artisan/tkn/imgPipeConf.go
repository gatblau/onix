package tkn

import (
	"bytes"
	"fmt"
	"text/template"
)

type imgPipeConf struct {
	AppName               string
	BuilderImage          string
	PushImageRegistry     string
	PushImageRepository   string
	PushImageName         string
	PushImageTag          string
	PullImageRegistry     string
	PullImageRegistryUser string
	PullImageRegistryPwd  string
	PushImageRegistryUser string
	PushImageRegistryPwd  string
	ArtefactRegistryUser  string
	ArtefactRegistryPwd   string
	GitToken              string
	DockerfileURL         string
	ArtefactName          string
	BuildFileURL          string
	ArtefactKeyName       string
}

func (c *imgPipeConf) buildTaskName() string {
	return fmt.Sprintf("%s-image-build-task", c.AppName)
}

func (c *imgPipeConf) getArgs() (string, error) {
	tmpl, err := template.New("image_build_args").Parse(imgBuildArgs)
	if err != nil {
		return "", fmt.Errorf("cannot create template to merge arguments: %s", err)
	}
	var buf = new(bytes.Buffer)
	err = tmpl.Execute(buf, c)
	if err != nil {
		return "", fmt.Errorf("cannot execute template to merge arguments: %s", err)
	}
	return buf.String(), nil
}

const imgBuildArgs = `
	- |-
	  # required vairables for buidconfig
	  export APPLICATION_NAME={{.AppName}}
	  export BUILDAH_IMAGE={{.BuilderImage}}
	  export PUSH_IMAGE_REGISTRY={{.PushImageRegistry}}
	  export PUSH_IMAGE_REPO={{.PushImageRepository}}
	  export PUSH_IMAGE_NAME={{.PushImageName}}
	  export PUSH_IMAGE_VERSION={{.PushImageTag}}
	  export PULL_IMAGE_REGISTRY={{.PullImageRegistry}}
	  export PULL_IMAGE_REGISTRY_UNAME={{.PullImageRegistryUser}}
	  export PULL_IMAGE_REGISTRY_PWD={{.PullImageRegistryPwd}}
	  export PUSH_IMAGE_REGISTRY_UNAME={{.PushImageRegistryUser}}
	  export PUSH_IMAGE_REGISTRY_PWD={{.PushImageRegistryPwd}}
	  export ARTEFACT_REG_USER={{.ArtefactRegistryUser}}
	  export ARTEFACT_REG_PWD={{.ArtefactRegistryPwd}}
	  export GIT_TOKEN={{.GitToken}}
	  export DOCKER_FILE_URL={{.DockerfileURL}}
	  export ARTEFACT_NAME={{.ArtefactName}}
	  export BUILD_FILE_URL={{.BuildFileURL}}
	  export ARTEFACT_KEY_NAME={{.ArtefactKeyName}}

	  art merge /tmp/buildconfig.yaml.tem
	  if [ "$?" -ne 0 ]; then
		echo "failed to merge /tmp/buildconfig.yaml.tem"
		exit 1
	  fi

	  oc create -f /tmp/buildconfig.yaml

	  oc start-build sap-equip-jvm-custom-build -F | tee buildconfig.log

	  oc delete -f /tmp/buildconfig.yaml

	  error_count=$(grep -c error buildconfig.log)
	  if [ $error_count -gt 0 ]; then
		echo "Failed to build an image"
		exit 1
	  fi
`
