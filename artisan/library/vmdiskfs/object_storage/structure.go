package object_storage

// DiskConfiguration - image configuration structure
type DiskConfiguration struct {
	Cpu    int    `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
	Os     string `yaml:"os" json:"os"`
	Arch   string `yaml:"arch" json:"arch"`
	Source string `yaml:"source" json:"source"`
	Disk   []Disk `yaml:"disk" json:"disk"`
}
type Disk struct {
	Name     string `yaml:"name" json:"name"`
	DiskID   int    `yaml:"disk_id" json:"disk_id"`
	BootDisk bool   `yaml:"boot_disk" json:"boot_disk"`
	Size     int    `yaml:"size" json:"size"`
}

// WebHookMessage -  structure for deploy webhook
type WebHookMessage struct {
	StorageUrl string        `yaml:"url" json:"url"`
	Cpu        int           `yaml:"cpu" json:"cpu"`
	Memory     string        `yaml:"memory" json:"memory"`
	Os         string        `yaml:"os" json:"os"`
	Arch       string        `yaml:"arch" json:"arch"`
	Source     string        `yaml:"source" json:"source"`
	Disk       []WebHookDisk `yaml:"disk" json:"disk"`
}
type WebHookDisk struct {
	Name     string `yaml:"name" json:"name"`
	DiskID   int    `yaml:"disk_id" json:"disk_id"`
	BootDisk bool   `yaml:"boot_disk" json:"boot_disk"`
	Size     int    `yaml:"size" json:"size"`
	Url      string `yaml:"url" json:"url"`
}
