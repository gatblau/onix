module github.com/gatblau/onix/artisan

go 1.18

replace github.com/gatblau/onix/oxlib => ../oxlib

require (
	github.com/AlecAivazis/survey/v2 v2.3.1
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/compose-spec/compose-go v1.0.8
	github.com/eclipse/paho.mqtt.golang v1.3.5 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gatblau/oxc v0.0.0-20210810120109-3c7f200d87d2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/kr/pty v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-shellwords v1.0.12
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/ohler55/ojg v1.12.5
	github.com/pelletier/go-toml v1.9.4
	github.com/rs/zerolog v1.24.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e
	gopkg.in/yaml.v2 v2.4.0
)
