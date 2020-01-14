module gatblau.org/onix/terra

go 1.13

replace gatblau.org/onix/wapic => ../wapic

require (
	gatblau.org/onix/wapic v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
	github.com/prometheus/client_golang v0.9.3
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
)
