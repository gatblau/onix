module github.com/gatblau/onix/artisan/runner/host

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../../../artisan
	github.com/gatblau/onix/oxlib => ../../../oxlib
)

require (
	github.com/gatblau/onix/artisan v0.0.0-20220216112625-36146b593961
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/swaggo/swag v1.8.0
)
