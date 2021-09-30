module github.com/gatblau/onix/piloth

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/client => ../client
	github.com/gatblau/onix/pilotctl => ../pilotctl
)

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.3.0
	github.com/rs/zerolog v1.24.0
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
