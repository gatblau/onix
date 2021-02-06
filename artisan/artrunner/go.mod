module github.com/gatblau/onix/artisan/artrunner

go 1.15

replace github.com/gatblau/onix/artisan => ../

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/swaggo/http-swagger v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
)
