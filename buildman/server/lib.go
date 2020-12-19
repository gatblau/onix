package server

import (
	"bytes"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
)

// gets an instance of the k8s client
func getKubeClient() (kubernetes.Interface, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v.", err)
		return nil, err
	}
	return client, nil
}

// gets the K8S client configuration either inside or outside of the cluster depending on
// whether the kube config file could be found
func getKubeConfig() (*rest.Config, error) {
	// k8s client configuration
	var config *rest.Config

	// gets the path to the kube config file
	kubeConfigFile := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))

	if _, err := os.Stat(kubeConfigFile); err == nil {
		log.Print("Kube config file found: attempting out of cluster configuration.")
		// if the kube config file exists then do an outside of cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			log.Printf("Could not create out of cluster configuration: %v.", err)
			return nil, err
		}
	} else if os.IsNotExist(err) {
		log.Print("Kube config file not found: attempting in cluster configuration.")
		// the kube config file was not found then do in cluster configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("could not find the K8S client configuration. "+
				"are you running Sentinel in a container that has not been deployed in Kubernetes?.\n "+
				"the error message was: %v.", err)
			return nil, err
		}
	} else {
		// kube config might be there or not but it failed anyway :(
		if err != nil {
			log.Printf("could not figure out the Kube client configuration: %v", err)
			return nil, err
		}
	}
	return config, nil
}

// returns a merged template with the passed-in data
func merge(templateText string, data interface{}) (string, error) {
	w := new(bytes.Buffer)
	t, err := template.New("template").Parse(templateText)
	if err != nil {
		return "", err
	}
	err = t.Execute(w, data)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

type pipelineRunData struct {
	AppName   string
	TimeStamp string
	Namespace string
}

func newPipelineRunData(appName, namespace string) *pipelineRunData {
	t := time.Now()
	return &pipelineRunData{
		AppName:   appName,
		Namespace: namespace,
		TimeStamp: fmt.Sprintf("%02d%02d%02s%02d%02d%02d%s", t.Day(), t.Month(), strconv.Itoa(t.Year())[:2], t.Hour(), t.Minute(), t.Second(), strconv.Itoa(t.Nanosecond())[:3]),
	}
}

const imageBuilderRun = `
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: {{.AppName}}-image-pr-{{.TimeStamp}}
  namespace: {{.Namespace}}
spec:
  serviceAccountName: pipeline
  pipelineRef:
    name: {{.AppName}}-image-builder
  params:
    - name: deployment-name
      value: {{.AppName}}
`
