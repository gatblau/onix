module github.com/gatblau/onix/artisan/art

go 1.15

replace (
	github.com/gatblau/onix/artisan => ../
	github.com/gatblau/oxc => ../../../oxc
)

require (
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	gopkg.in/yaml.v2 v2.4.0
)
