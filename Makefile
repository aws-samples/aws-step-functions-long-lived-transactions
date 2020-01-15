OUTPUT = main # Referenced as Handler in template.yaml
PACKAGED_TEMPLATE = packaged.yaml
S3_BUCKET := $(S3_BUCKET)
STACK_NAME := $(STACK_NAME)
TEMPLATE = template.yaml

.PHONY: install test clean build lambda teardown package

install:
		go get -u ./...

test: 
		go test ./...

clean:
		echo "Begin cleaning..."
		rm -f */**/$(OUTPUT) $(PACKAGED_TEMPLATE)
		$(MAKE) clean -C inventory/release
	 	$(MAKE) clean -C inventory/reserve
	 	$(MAKE) clean -C order/new
	 	$(MAKE) clean -C order/update
	 	$(MAKE) clean -C payment/pay
	 	$(MAKE) clean -C payment/refund
		echo "Done cleaning..."

lambda:
	echo "Begin release build..."
	$(MAKE) -C inventory/release
	$(MAKE) -C inventory/reserve
	$(MAKE) -C order/new
	$(MAKE) -C order/update
	$(MAKE) -C payment/pay
	$(MAKE) -C payment/refund
	echo "Done release build."

build:  clean lambda

deploy: build
	sam package \
    --template-file $(TEMPLATE) \
    --output-template-file $(PACKAGED_TEMPLATE) \
    --s3-bucket $(S3_BUCKET)
	
	sam deploy \
    --template-file $(PACKAGED_TEMPLATE) \
    --stack-name $(STACK_NAME) \
    --capabilities CAPABILITY_IAM

.PHONY: teardown
teardown:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)