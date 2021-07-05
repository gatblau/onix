module github.com/gatblau/onix/piloth

go 1.15

replace github.com/gatblau/onix/artisan => ../artisan

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-20210628201641-ee0039d656de
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-shellwords v1.0.11
	github.com/rs/zerolog v1.23.0
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/tklauser/go-sysconf v0.3.6 // indirect
)
