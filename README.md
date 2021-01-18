# HCL to Markdown

**Warning:** We struggle to support `type = string` without quotes.

To write good documentation for your terraform module, quite often we need to print all our variables with all descriptions.

We offer you to avoid long bash scripts and update your README.md as quick as it was released.

## How to use

Download the latest binary for your OS from Releases Page and run it over required file

```shell
$ ./hcltomd --file ./testdata/inputs.tf
|      NAME      |  TYPE  |        DEFAULT        |                      DESCRIPTION                       |
|----------------|--------|-----------------------|--------------------------------------------------------|
| aws_region     | string | us-east-1             | The AWS region in which all resources will be created. |
| vpc_id         | string | 1                     | The id of the specific VPC to retrieve.                |
| instance_count | number | 2                     | The count of desired instances of EC2.                 |
| zones          | list   | [us-east-1 us-east-2] | The selected zones                                     |
```

That's it, you can copy it into your Markdown documentation
