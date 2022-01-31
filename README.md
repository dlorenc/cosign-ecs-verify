# cosign-ecs-verify

In this demo, we build an analog of a [Kubernetes admission controller] for
Amazon [Elastic Container Service (ECS)][ECS] that checks all images to be run
for a valid [cosign] signature with a given key in AWS [KMS].

**NOTE:** This is demonstration code and as such shouldn't be used in
production. In the event of misconfiguration or a bug, it can prevent all ECS
containers from running.

[Kubernetes admission controller]: https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/
[ECS]: https://aws.amazon.com/ecs/
[cosign]: https://github.com/sigstore/cosign
[KMS]: https://aws.amazon.com/kms/

## How it works

![](aws-ecs-cosign-verify.png)

1. Start an ECS task in the cluster
2. The task definition has the container image stored in ECR
3. EventBridge sends a notification to Lambda
4. Cluster and Task definition is sent to function 
5. KMS key that has signed an image 
6. Lambda function evaluates if container image is signed w/ KMS
7. If not signed with specified key it does two things
   1. Stop task definition
   2. SNS notification email to alert that the service/task has been stopped

## Requirements and preliminaries.

For this demo, you will need the following tools installed:

- `make` (e.g., [GNU make])
- [Terraform]: for local testing
- [AWS CLI] and [AWS SAM CLI]: for deploying
- [`cosign`]: for generating keys
- [`docker`]: if you need to make images

[AWS CLI]: https://aws.amazon.com/cli/
[AWS SAM CLI]: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
[GNU make]: https://www.gnu.org/software/make/
[Terraform]: https://www.terraform.io/downloads
[`cosign`]: https://github.com/sigstore/cosign
[`docker`]: https://docs.docker.com/get-docker/

You should [configure the AWS CLI] for your project and account.

[configure the AWS CLI]: https://docs.aws.amazon.com/cli/latest/reference/configure/

We need a key against which to verify image signatures. If you have an existing
keypair for cosign in AWS KMS, set it:

## Deploy

To deploy, run:

```shell
make sam_deploy
```

This uses a SAM template (`template.yml`) to create:

- The Serverless function (source in cosign-ecs-function)
  - Triggered on cloud watch event: ecs task/container state change
  - for each container in the event
    - get the key corresponding to the region
    - verify the container image
- An Amazon [SNS] topic: if the function stops an unsigned container image, it
  will send a message to this topic.
  - You can [configure email notifications][sns-email] for this topic to be
    alerted whenever an unverified image is stopped.
  - Messages sent to the topic are [encrypted using a key in KMS][sns-kms].

[SNS]: https://aws.amazon.com/sns/
[sns-email]: https://docs.aws.amazon.com/sns/latest/dg/sns-email-notifications.html
[sns-kms]: https://aws.amazon.com/blogs/compute/encrypting-messages-published-to-amazon-sns-with-aws-kms/
    

## Test it

### Deploy a cluster and run tasks

The `terraform` subdirectory contains a Terraform template for an ECS cluster
and task definitions for running our signed/unsigned tasks. First, initialize
it (this will download required providers):

``` shell
make tf_init
```

Then, deploy the template:

``` shell
make tf_apply  # run `make tf_plan` to see the plan first
```

We can then run our tasks:

``` shell
make run_unsigned_task
make run_signed_task
```

*Note:* this will run on the tasks on a subnet of the [default VPC].

[default VPC]: https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html

Check:

``` shell
make task_status
```

You should see the unsigned task in the `STOPPED` tasks and the signed task in the `RUNNING` tasks.


### Cleanup

``` shell
make stop_tasks
make tf_destroy
make sam_delete
```


## Local Dev

- Go 1.17

``` shell
make sam_local 
make sam_local_debug
```
