module github.com/gatblau/onix/dbman

go 1.16

replace github.com/gatblau/onix/oxlib => ../oxlib

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/gatblau/oxc v0.0.0-20210810120109-3c7f200d87d2
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/go-plugin v1.4.3
	github.com/jackc/pgconn v1.10.1
	github.com/jackc/pgtype v1.9.1
	github.com/jackc/pgx/v4 v4.14.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.19.0 // indirect
	github.com/spf13/afero v1.3.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/swaggo/http-swagger v1.1.2
	github.com/swaggo/swag v1.7.6
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
