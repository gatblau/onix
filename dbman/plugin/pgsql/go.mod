module github.com/gatblau/onix/dbman/plugins/pgsql

go 1.13

replace github.com/gatblau/onix/dbman => ../../

require (
	github.com/gatblau/onix/dbman v0.0.0-20200623160749-05451f11f8c1
	github.com/hashicorp/go-plugin v1.3.0
	github.com/jackc/pgconn v1.6.0
	github.com/jackc/pgx/v4 v4.6.0
)
