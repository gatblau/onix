module github.com/gatblau/onix/piloth

go 1.15

replace (
	github.com/gatblau/onix/artisan => ../artisan
	//github.com/gatblau/onix/pilotctl => ../pilotctl
)

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/gatblau/onix/rem v0.0.0-20210605115117-6b8709b4a45d
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/rs/zerolog v1.22.0
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/tklauser/go-sysconf v0.3.6 // indirect
)
