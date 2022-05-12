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
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/rs/zerolog v1.24.0
	github.com/schollz/progressbar/v3 v3.8.6 // indirect
	github.com/shirou/gopsutil v3.21.8+incompatible
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/swaggo/swag v1.7.9
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/xuri/excelize/v2 v2.5.0
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
)
