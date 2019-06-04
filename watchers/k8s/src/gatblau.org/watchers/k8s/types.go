package main

type Config struct {
	Handler EventHandler
	Observe ObservedResources
}

/*

 */
type ObservedResources struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"pod"`
	Job                   bool `json:"job"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"namespace"`
	Secret                bool `json:"secret"`
	ConfigMap             bool `json:"configmap"`
	Ingress               bool `json:"ingress"`
}

/*

 */
type EventHandler interface {
	OnCreate(e Event, o interface{})
	OnDelete(e Event, o interface{})
	OnUpdate(e Event, o interface{})
}

type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Event represent an event got from k8s api server
// Events from different endpoints need to be casted to KubewatchEvent
// before being able to be handled by Handler
type KubeEvent struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Status    string
	Name      string
}
