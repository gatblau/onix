module github.com/gatblau/onix/dbman

go 1.13

//replace github.com/gatblau/onix/dbman => ./

require (
	github.com/gatblau/oxc v0.0.0-20200518102735-38237e6a1005
	github.com/ghodss/yaml v1.0.0
	github.com/google/wire v0.4.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/terraform v0.12.26
	github.com/jackc/pgconn v1.5.0
	github.com/jackc/pgx/v4 v4.6.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/prometheus/client_golang v1.0.0
	github.com/rs/zerolog v1.18.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	google.golang.org/api v0.13.0
	gopkg.in/yaml.v2 v2.3.0
)
