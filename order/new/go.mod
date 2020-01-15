module aws-step-functions-long-lived-transactions/order/new

go 1.13

require (
	aws-step-functions-long-lived-transactions/models v0.0.0-00010101000000-000000000000
	github.com/aws-samples/aws-step-functions-long-lived-transactions v0.0.0-20191001081655-c5ca1f79a412 // indirect
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.26
	github.com/aws/aws-xray-sdk-go v0.9.4
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
)

replace aws-step-functions-long-lived-transactions/models => ../../models
