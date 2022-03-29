module github.com/gatblau/onix/piloth

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/oxlib => ../oxlib
	github.com/gatblau/onix/pilotctl => ../pilotctl
)

require (
	github.com/ProtonMail/gopenpgp/v2 v2.2.4
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-00010101000000-000000000000
	github.com/pkg/profile v1.6.0
	github.com/reugn/go-quartz v0.3.6
	github.com/rs/zerolog v1.24.0
	github.com/schollz/peerdiscovery v1.6.9
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.3.2 // indirect
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
