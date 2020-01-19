# Setup Information

## Using AWS CloudFormation (preferred)

The easiest way to get started is to launch an AWS CloudFormation template that will deploy the resources for this workshop.

Region| Launch
------|-----
US East (N. Virginia) | [![Launch in us-east-1](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/images/cloudformation-launch-stack-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/create/review?stackName=reinvent-builder-session-401&templateURL=https://s3.amazonaws.com/aws-step-functions-long-lived-transactions-us-east-1/template.yaml)
US West (Oregon) | [![Launch in us-west-2](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/images/cloudformation-launch-stack-button.png)](https://console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/create/review?stackName=reinvent-builder-session-401&templateURL=https://s3-us-west-2.amazonaws.com/aws-step-functions-long-lived-transactions-us-west-2/template.yaml)
Asia Pacific (Sydney) | [![Launch in ap-southeast-2](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/images/cloudformation-launch-stack-button.png)](https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-2#/stacks/create/review?stackName=reinvent-builder-session-401&templateURL=https://s3-ap-southeast-2.amazonaws.com/aws-step-functions-long-lived-transactions-ap-southeast-2/template.yaml)

### CloudFormation Setup Instructions

1. Click the **Launch Stack** link above for the region of your choice.

1. Name the stack (or leave the default `reinvent-builder-session-401`)

1. In the Capabilities section acknowledge that CloudFormation will create IAM resources and click **Create**.
    ![Acknowledge IAM Screenshot](images/capabilities.png)

1. Select `Create Change Set`

1. Once the Change Set has been successfully created, select `Execute` to create the stack.

</p></details>

## Deploy from Source
Alternatively, you can deploy the source from you local development environment. Please note, this requires the following environment setup.

<details>
<summary><strong>Deploy from Source(expand for details)</strong></summary><p>

### Requirements

* [aws-cli](https://aws.amazon.com/cli/) already configured with Administrator permissions.
* [sam-cli](https://github.com/awslabs/aws-sam-cli) AWS SAM CLI tool for local development and testing of Serverless applications
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* Make (see instructions below)

<br/>
<details>
<summary><strong>Installing SAM CLI</strong></summary><p>

**Brew for Mac and Linux**

You can install SAM CLI using brew, a popular package manager for installing the packages you need. Installation is as simple as:

```shell
brew tap aws/tap
brew install aws-sam-cli
```

> **NOTE:** On a Mac you use [Homebrew](https://brew.sh/), and on Linux you use [Linuxbrew](http://linuxbrew.sh/) (a fork of the Homebrew package manager).

**MSI for Windows**

You can now download an MSI to install SAM CLI on Windows. Get the MSI you need here:

* [64-bit](https://github.com/awslabs/aws-sam-cli/releases/download/v0.6.2/AWS_SAM_CLI_64_PY3.msi)
* [32-bit](https://github.com/awslabs/aws-sam-cli/releases/download/v0.6.2/AWS_SAM_CLI_32_PY3.msi)

</p></details>

<br/>
<details>
<summary><strong>Installing Golang</strong></summary><p>

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quick-start way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

### Configuring GoPATH

In order to build the source locally you need to set up our Go development environment.

Follow the instructions as outlined here https://github.com/golang/go/wiki/SettingGOPATH

</p></details>

<br/>
<details>
<summary><strong>Installing Make</strong></summary><p>
**Make for Mac**

`xcode-select --install`

**Make for Windows**

Installation files can be downloaded from http://gnuwin32.sourceforge.net/packages/make.htm
</p></details>

### Clone the repository

Once you have you GOPATH configured clone the builder session repository into the following directory under your GOPATH

```shell
mkdir $GOPATH/src/github.com/aws-samples

git clone https://github.com/aws-samples/aws-step-functions-long-lived-transactions.git
```

### Installing Dependencies

In this project, we use the Makefile to execute all the go commands. The first thing we need to do download all the solution dependencies:

```shell
make install
```

> This will execute the built-in `go get` and download all packages to your specified GOPATH.

### Building

Golang is a statically compiled language, meaning in order to run it you have to build the executable target. As there are a number of functions to create we will use the Makefile to build all projects. You can issue the following command in a shell to build it:

```shell
make build
```

> The Makefile executes the following `go build` command for each function to make sure they are compatible with the AWS Lambda system architecture.
>
> **Example:** `GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world`
>
> **NOTE**: If you're not building the function on a Linux machine, you will need to specify the `GOOS` and `GOARCH` environment variables, this allows Golang to build your function for another system architecture and ensure compatibility.

### Deploy

We will use the AWS SAM CLI to install the function.

> **See [Serverless Application Model (SAM) HOWTO Guide](https://github.com/awslabs/serverless-application-model/blob/master/HOWTO.md) for more details in how to get started.**

First and foremost, we need a `S3 bucket` where we can upload our Lambda functions packaged as ZIP before we deploy anything - If you don't have a S3 bucket to store code artefacts then this is a good time to create one:

```shell
aws s3 mb s3://BUCKET_NAME --region YOUR_AWS_REGION
```

In your terminal, execute the following commands to package, deploy the serverless template:

```shell
sam package \
    --template-file template.yaml \
    --output-template-file packaged.yaml \
    --s3-bucket REPLACE_THIS_WITH_YOUR_S3_BUCKET_NAME \
    --region YOUR_AWS_REGION

sam deploy \
    --template-file packaged.yaml \
    --stack-name reinvent-builder-session-401 \
    --capabilities CAPABILITY_IAM
```

The following command describes the outputs defined within the cloudformation stack:

```shell
aws cloudformation describe-stacks \
    --stack-name reinvent-builder-session-401 --query 'Stacks[].Outputs'
```
</p></details>

## Completion

Once you have successfully deployed the functions, go ahead and start building your saga.

See the [Session Guide](guide.md) for more information.