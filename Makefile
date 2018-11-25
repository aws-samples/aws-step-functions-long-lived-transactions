OUTPUT = main # Referenced as Handler in template.yaml
PACKAGED_TEMPLATE = packaged.yaml
S3_BUCKET := $(S3_BUCKET)
STACK_NAME := $(STACK_NAME)
TEMPLATE = template.yaml

.PHONY: install test clean build package

install:
		go get -u ./...

test: 
		go test ./...

clean:
		echo "Begin cleaning..."
		rm -f */**/$(OUTPUT) $(PACKAGED_TEMPLATE)
		echo "Done cleaning..."

build:  clean
		echo "Begin release build..."
		GOOS=linux GOARCH=amd64 go build -o ./order/new/main ./order/new
		GOOS=linux GOARCH=amd64 go build -o ./order/update/main ./order/update
		GOOS=linux GOARCH=amd64 go build -o ./payment/pay/main ./payment/pay
		GOOS=linux GOARCH=amd64 go build -o ./payment/refund/main ./payment/refund
		GOOS=linux GOARCH=amd64 go build -o ./inventory/reserve/main ./inventory/reserve
		GOOS=linux GOARCH=amd64 go build -o ./inventory/release/main ./inventory/release
		echo "Done release build."

debug:
		echo "Begin debug build..."
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./order/new/main ./order/new
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./order/update/main ./order/update
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./payment/pay/main ./payment/pay
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./payment/refund/main ./payment/refund
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./inventory/reserve/main ./inventory/reserve
		GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o ./inventory/release/main ./inventory/release
		echo "Done debug build."

deploy: build
	sam package \
    --template-file $(TEMPLATE) \
    --output-template-file $(PACKAGED_TEMPLATE) \
    --s3-bucket $(S3_BUCKET)
	
	sam deploy \
    --template-file $(PACKAGED_TEMPLATE) \
    --stack-name $(STACK_NAME) \
    --capabilities CAPABILITY_IAM