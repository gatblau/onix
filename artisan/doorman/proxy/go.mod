module github.com/gatblau/onix/artisan/doorman/proxy

go 1.16

replace github.com/gatblau/onix/oxlib => ../../../oxlib

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gatblau/onix/oxlib v0.0.0-20220218080420-10c5cf8ab357
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/swaggo/swag v1.7.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
)
