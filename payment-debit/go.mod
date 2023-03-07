module aws-step-functions-long-lived-transactions/paymentment-debit

go 1.16

require (
	aws-step-functions-long-lived-transactions/models v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.70
	github.com/aws/aws-xray-sdk-go v1.5.0
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.34.0 // indirect
	google.golang.org/genproto v0.0.0-20210630183607-d20f26d13c79 // indirect
	google.golang.org/grpc v1.39.0 // indirect
)

replace aws-step-functions-long-lived-transactions/models => ../models
