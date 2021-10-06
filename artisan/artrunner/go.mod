module github.com/gatblau/onix/artisan/artrunner

go 1.15

replace (
	github.com/gatblau/onix/artisan => ../
	github.com/gatblau/onix/oxlib => ../../oxlib
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gatblau/oxc v0.0.0-20210810120109-3c7f200d87d2
	github.com/gorilla/mux v1.8.0
	github.com/swaggo/swag v1.7.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
)
