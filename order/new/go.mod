module aws-step-functions-long-lived-transactions/order/new

go 1.16

require (
	aws-step-functions-long-lived-transactions/models v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.69
	github.com/aws/aws-xray-sdk-go v1.5.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
)

replace aws-step-functions-long-lived-transactions/models => ../../models
