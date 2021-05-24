module github.com/gatblau/onix/piloth

go 1.15

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/rem => ../rem
)

require (
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/gatblau/onix/pilot v0.0.0-20210519181728-1ecec3949f84
	github.com/gatblau/onix/rem v0.0.0-00010101000000-000000000000
	github.com/gatblau/oxc v0.0.0-20210523084722-f08170feef8e
	github.com/joho/godotenv v1.3.0
	github.com/rs/zerolog v1.22.0
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/tklauser/go-sysconf v0.3.6 // indirect
)
