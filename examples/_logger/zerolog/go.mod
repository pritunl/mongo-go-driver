module go.mongodb.go/mongo-driver/examples/logger/zerolog

go 1.20

replace github.com/pritunl/mongo-go-driver => ../../../

require (
	github.com/go-logr/zerologr v1.2.2
	github.com/rs/zerolog v1.28.0
	// Note that the Go driver version is replaced with the local Go driver code
	// by the replace directive above.
	github.com/pritunl/mongo-go-driver v1.11.7
)

require (
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
