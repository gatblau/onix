module github.com/gatblau/onix/artisan/doorman

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../
	github.com/gatblau/onix/oxlib => ../../oxlib
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/eclipse/paho.mqtt.golang v1.3.5 // indirect
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/rs/zerolog v1.24.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/swaggo/swag v1.7.0
	go.mongodb.org/mongo-driver v1.8.3
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sys v0.0.0-20210902050250-f475640dd07b // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
)
