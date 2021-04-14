module github.com/gatblau/onix/rem

go 1.15

replace github.com/gatblau/onix/artisan => ../artisan

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/swaggo/swag v1.7.0
)
