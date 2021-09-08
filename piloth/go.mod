module github.com/gatblau/onix/piloth

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/pilotctl => ../pilotctl
	github.com/gatblau/oxc => ../../oxc
)

require (
	github.com/AlecAivazis/survey/v2 v2.3.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-00010101000000-000000000000
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jackc/pgconn v1.8.1 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/rs/zerolog v1.24.0
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
