module github.com/gatblau/onix/artisan/art

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../
	github.com/gatblau/onix/oxlib => ../../oxlib
)

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/minio/minio-go/v7 v7.0.21 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	gopkg.in/yaml.v2 v2.4.0
)
