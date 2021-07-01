module aws-step-functions-long-lived-transactions/payment/pay

go 1.16

require (
	aws-step-functions-long-lived-transactions/models v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.26
	github.com/aws/aws-xray-sdk-go v1.0.0-rc.14
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
)

replace aws-step-functions-long-lived-transactions/models => ../../models
