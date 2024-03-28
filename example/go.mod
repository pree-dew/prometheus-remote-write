module sample-app

go 1.21.0

require github.com/prometheus/client_golang v1.19.0

replace github.com/pree-dew/prometheus-remote-write => ../

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/pree-dew/prometheus-remote-write v0.0.0-00010101000000-000000000000 // indirect
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/common v0.49.1-0.20240306132007-4199f18c3e92 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/prometheus/prometheus v0.51.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
