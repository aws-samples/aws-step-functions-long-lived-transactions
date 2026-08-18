# Setup Information

You deploy this sample from source using the AWS SAM CLI.

## Requirements

* [AWS CLI](https://aws.amazon.com/cli/) configured with permissions to create IAM roles, Lambda functions, DynamoDB tables, SNS/SQS resources, and Step Functions state machines
* [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
* [Go](https://go.dev/doc/install) 1.26 or later
* Make (optional — convenience targets in the root `Makefile`)

> Docker is **not** required: the functions are built with SAM's native Go builder, which cross-compiles for the `provided.al2023` arm64 runtime on any host.

## Clone the repository

Go modules are used throughout, so you can clone the repository anywhere (no `GOPATH` setup needed):

```shell
git clone https://github.com/aws-samples/aws-step-functions-long-lived-transactions.git
cd aws-step-functions-long-lived-transactions
```

The repository contains one Go module per Lambda function plus a shared `models` module, tied together for local development with a `go.work` workspace file. To run all unit tests:

```shell
make test
```

## Deploy the sample application

To build and deploy the application for the first time, run:

```shell
sam build
sam deploy --guided

Configuring SAM deploy
======================

        Looking for config file [samconfig.toml] :  Found
        Reading default arguments  :  Success

        Setting default arguments for 'sam deploy'
        =========================================
        Stack Name [sfn-saga]: 
        AWS Region [ap-southeast-2]: 
        #Shows you resources changes to be deployed and require a 'Y' to initiate deploy
        Confirm changes before deploy [y/N]: N
        #SAM needs permission to be able to create roles to connect to the resources in your template
        Allow SAM CLI IAM role creation [Y/n]: Y
        Save arguments to configuration file [Y/n]: Y
        SAM configuration file [samconfig.toml]: 
        SAM configuration environment [default]: 
```

The first command builds all six functions. The second packages and deploys the application to AWS with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.

* **AWS Region**: The AWS region you want to deploy your app to.

* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.

* **Allow SAM CLI IAM role creation**: This template creates scoped-down AWS IAM roles for the Lambda functions and state machine. To deploy a stack that creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.

* **Save arguments to samconfig.toml**: If set to yes, your choices are saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

The following command describes the outputs defined within the CloudFormation stack:

```shell
aws cloudformation describe-stacks \
    --stack-name aws-sfn-saga --query 'Stacks[].Outputs'
```

## Completion

Once you have successfully deployed the application, go ahead and start testing the saga.

See the [Session Guide](guide.md) for more information.

## Clean up

To delete the sample application that you created, use the AWS CLI. Assuming you used your project name for the stack name, you can run the following:

```shell
aws cloudformation delete-stack --stack-name aws-sfn-saga
```
