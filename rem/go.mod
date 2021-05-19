module github.com/gatblau/onix/rem

go 1.15

replace github.com/gatblau/onix/artisan => ../artisan

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/artisan v0.0.0-00010101000000-000000000000
	github.com/gatblau/oxc v0.0.0-20201021102807-48791fe073d6
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/reugn/go-quartz v0.3.4
	github.com/swaggo/swag v1.7.0
)