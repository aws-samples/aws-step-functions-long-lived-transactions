# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

# Set STACK_NAME to override the default stack name for deploy/teardown.
STACK_NAME ?= aws-sfn-saga

MODULES = models order-new order-update payment-debit payment-credit inventory-reserve inventory-release

.PHONY: test
test:
	@for m in $(MODULES); do \
		echo "=== $$m ==="; \
		(cd $$m && go test ./...) || exit 1; \
	done

.PHONY: vet
vet:
	@for m in $(MODULES); do \
		echo "=== $$m ==="; \
		(cd $$m && go vet ./...) || exit 1; \
	done

.PHONY: tidy
tidy:
	@for m in $(MODULES); do \
		echo "=== $$m ==="; \
		(cd $$m && go mod tidy) || exit 1; \
	done

.PHONY: build
build:
	sam build

.PHONY: validate
validate:
	sam validate --lint

.PHONY: deploy
deploy: build
	sam deploy --stack-name $(STACK_NAME) --resolve-s3 --capabilities CAPABILITY_IAM --no-confirm-changeset --no-fail-on-empty-changeset

.PHONY: clean
clean:
	rm -rf .aws-sam

.PHONY: teardown
teardown:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
