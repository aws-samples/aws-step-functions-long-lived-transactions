# Session Guide

## What's in this repository?

This is a sample template for Managing Long Lived Transactions with AWS Step Functions. Below is a brief explanation of what we have created for you:

``` bash
.
├── Makefile              <-- Make to automate build
├── docs                  <-- Workshop guide and setup instructions
│   ├── guide.md
│   └── setup.md
├── inventory
│   ├── release           <-- Lambda function code represents compensating transaction to release inventory
│   └── reserve           <-- Lambda function code represents task to reserve order items from the inventory
├── models                <-- Models package that defines the types used by the various functions and state data
│   ├── inventory.go
│   ├── order.go
│   └── payment.go
├── order                 <-- Lambda function code represents task to create a new order and set status to "new order"
│   ├── new
│   └── update
├── payment
│   ├── pay               <-- Lambda function code represents task to process financial transaction for the order
│   └── refund            <-- Lambda function code represents the compensating transaction to refund customer order
├── state-machine.json    <-- Sample saga implementation with Step Functions
└── template.yaml         <-- SAM template for defining and deploying serverless application resources
                              [USE THIS AS A GUIDE IF YOU GET STUCK]

```

## Amazon States Language

A full description of the how to describe your state machine can be found on the Amazon States Language specification (see the resources section at the bottom of this page).

Please review the "Templates" section in the [AWS Console](https://console.aws.amazon.com/states/home) for examples of how you can implement various states.

### Useful snippets

#### Task state

The Task State (identified by "Type":"Task") causes the interpreter to execute the work identified by the state's “Resource” field.

```json
"TaskName": {
  "Comment": "Add comment...",
  "Type": "Task",
  "Resource": "AWS_LAMBDA_ARN",
  "Next": "",
  "TimeoutSeconds": 300,
  "HeartbeatSeconds": 60
}
```

#### Catch
Any state can encounter runtime errors. Errors can arise because of state machine definition issues (e.g. the “ResultPath” problem discussed immediately above), task failures (e.g. an exception thrown by a Lambda function) or because of transient issues, such as network partition events.

```json
"Catch": [
  {
    "ErrorEquals": ["CustomError"],
    "Next": ""
  },
  {
    "ErrorEquals": ["States.ALL"],
    "Next": ""
  }
]
```

#### Retry
Task States and Parallel States MAY have a field named “Retry”, whose value MUST be an array of objects, called Retriers.

When a	state reports an error, the interpreter scans through the Retriers and, when the Error Name appears in the value of of a Retrier’s “ErrorEquals” field, implements the retry policy described in that Retrier.

```json
"Retry": [
  {
    "ErrorEquals": ["CustomError"],
    "IntervalSeconds": 1,
    "MaxAttempts": 2,
    "BackoffRate": 2.0
  },
  {
    "ErrorEquals": ["States.ALL"],
    "IntervalSeconds": 30,
    "MaxAttempts": 2,
    "BackoffRate": 2.0
  }
]
```

## Invoking your Step Function via CLI

The AWS CLI command will trigger a execution of your state machine. Make sure you substitute the ARN for the state machine in your account. You can find the ARN in the AWS CloudFormation Output section or in the AWS Step Functions console.

![CloudFormation Output](images/cfn-output.png)

> `--region` must match the region you have deployed the application stack into. This is optional if you're using your default region.

``` bash
aws stepfunctions start-execution \
    --state-machine-arn "arn:aws:states:[REGION]:[ACCOUNT NUMBER]:stateMachine:[STATEMACHINE-NAME]" \
    --input "{\"order_id\": \"40063fe3-56d9-4c51-b91f-71929834ce03\", \"order_date\": \"2018-10-19T10:50:16+08:00\", \"customer_id\": \"8d04ea6f-c6b2-4422-8550-839a16f01feb\", \"items\": [{ \"item_id\": \"567\", \"qty\": 1.0, \"description\": \"Cart item 1\", \"unit_price\": 199.99    }]}" \
    --region [AWS_REGION]
```

## Exceptions

The following is a list of all the exceptions thrown by the application and can be used in your state machine.

* **ErrProcessOrder** represents a process order error
* **ErrUpdateOrderStatus** represents a process order error
* **ErrProcessPayment** represents a process payment error
* **ErrProcessRefund** represents a process payment refund error
* **ErrReserveInventory** represents a inventory update error
* **ErrReleaseInventory** represents a inventory update reversal error

## Testing Scenarios

The AWS Step Functions implementation has been configured for you to be easily test the various scenarios of the saga implementation. Modifying your `order_id` with a specified prefix will trigger an error in the each Task.

OrderID Prefix | Will error with | Example | Expected execution
------------ | ------------- | --- | ---
1 | ErrProcessOrder | 1ae4501d-ed92-4b27-bf0e-fd978ed45127 | ![1](images/paths-breakdown-1.png) 
11 | ErrUpdateOrderStatus | 11328abd-368d-43fd-bd4f-db15b5b63951 | ![11](images/paths-breakdown-11.png)
2 | ErrProcessPayment | 20b0b599-441b-45c3-910e-ad63fe992c43 | ![2](images/paths-breakdown-2.png)
22 | ErrProcessRefund | 222f741b-0292-4f93-a2f7-503f92486955 | ![22](images/paths-breakdown-22.png)
3 | ErrReserveInventory | 3a7dc768-6f32-495d-a140-3d330c246f50 | ![3](images/paths-breakdown-3.png)
33 | ErrReleaseInventory | 33a49007-a815-4079-9b9b-e30ae7eca11f | ![3](images/paths-breakdown-33.png)
4 | No error | Order IDs beginning with 4 - 9 will pass successfully. For example:<br>47063fe3-56d9-4c51-b91f-71929834ce03<br>875e6c0a-9cd1-448d-94be-1f110fc3e5b3<br>9c57e9ae-966b-49d3-8b9f-5ba04cfe2533| ![4-9](images/paths-breakdown-7.png)


## Additional Resources

### Step Functions

* [AWS Step Functions](https://aws.amazon.com/step-functions/)
* [AWS Step Functions Developer Guide](https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html)
* [AWS Step Function Tutorials](https://docs.aws.amazon.com/step-functions/latest/dg/tutorials.html)
* [statelint](https://github.com/awslabs/statelint)
* [Amazon States Language](https://states-language.net/spec.html)

## How else can you implement this solution?

Is there any other way you can think of how to break this problem down? What other features of Step Functions could be employed to implement a saga pattern?