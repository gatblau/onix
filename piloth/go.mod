module github.com/gatblau/onix/piloth

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/oxlib => ../oxlib
	github.com/gatblau/onix/pilotctl => ../pilotctl
)

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.3.0
	github.com/reugn/go-quartz v0.3.6 // indirect
	github.com/rs/zerolog v1.24.0
	github.com/schollz/peerdiscovery v1.6.9 // indirect
	golang.org/x/net v0.0.0-20211007125505-59d4e928ea9d // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
