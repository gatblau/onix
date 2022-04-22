module github.com/gatblau/onix/dbman

go 1.16

replace github.com/gatblau/onix/oxlib => ../oxlib

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/go-plugin v1.4.3
	github.com/jackc/pgconn v1.10.1
	github.com/jackc/pgtype v1.9.1
	github.com/jackc/pgx/v4 v4.14.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.19.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.10.1
	github.com/swaggo/http-swagger v1.2.6
	github.com/swaggo/swag v1.7.9
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
