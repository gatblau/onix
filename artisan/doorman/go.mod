module github.com/gatblau/onix/artisan/doorman

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../
	github.com/gatblau/onix/oxlib => ../../oxlib
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/swaggo/swag v1.7.0
	go.mongodb.org/mongo-driver v1.8.3
)
