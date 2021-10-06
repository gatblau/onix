module github.com/gatblau/onix/pilotctl/receivers/mongo

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../../../artisan
	github.com/gatblau/onix/oxlib => ../../../oxlib // needed as it is a pilotctl dependency
	github.com/gatblau/onix/pilotctl => ../../../pilotctl
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/pilotctl v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/swaggo/swag v1.7.1
	go.mongodb.org/mongo-driver v1.7.2
)
