module aws-step-functions-long-lived-transactions/paymentment-debit

go 1.16

require (
	aws-step-functions-long-lived-transactions/models v0.0.0-00010101000000-000000000000
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.70
	github.com/aws/aws-xray-sdk-go v1.5.0
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.28.0 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20210630183607-d20f26d13c79 // indirect
	google.golang.org/grpc v1.39.0 // indirect
)

replace aws-step-functions-long-lived-transactions/models => ../models
