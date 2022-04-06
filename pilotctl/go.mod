module github.com/gatblau/onix/pilotctl

go 1.16

replace (
	github.com/gatblau/onix/artisan => ../artisan
	github.com/gatblau/onix/oxlib => ../oxlib
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/onix/oxlib v0.0.0-00010101000000-000000000000
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/rs/zerolog v1.24.0
	github.com/shirou/gopsutil v3.21.8+incompatible
	github.com/swaggo/swag v1.7.4
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/tools v0.1.7 // indirect
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
