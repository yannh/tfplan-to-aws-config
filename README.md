# TFPlan-to-AWS-Config

*Currently in proof-of-concept-phase*

AWS Config supports [the recording of third-party resources](https://docs.aws.amazon.com/config/latest/APIReference/API_PutResourceConfig.html) through their API.
This is a proof-of-concept to automatically record the state of non-AWS resources maintained using
Terraform to AWS Config.

```
$ terraform plan -state tfplan
$ terraform show -json tfplan > tfplan.json
$ tfplan-to-aws-config tfplan.json
```