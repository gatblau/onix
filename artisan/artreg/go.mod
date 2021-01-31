module github.com/gatblau/onix/artisan/artreg

go 1.15

replace github.com/gatblau/onix/artisan => ../

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/prometheus/client_golang v1.8.0
	github.com/swaggo/http-swagger v1.0.0
	gopkg.in/yaml.v2 v2.4.0
)
