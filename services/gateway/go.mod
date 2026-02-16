module github.com/Molov30/go-microservices/services/gateway

go 1.25.6

replace github.com/Molov30/go-microservices/generated => ../../generated

require (
	github.com/Molov30/go-microservices/generated v0.0.0
	github.com/caarlos0/env/v10 v10.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.18.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang/protobuf v1.5.4
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.34.0
	google.golang.org/grpc v1.79.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260203192932-546029d2fa20 // indirect
)
