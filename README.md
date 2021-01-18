# HCL to Markdown
[![Build Status](https://cloud.drone.io/api/badges/idestis/hcltomd/status.svg)](https://cloud.drone.io/idestis/hcltomd)

**Note:** In progress to add an abillity to parse unquoted type.

## About

To write good documentation for your terraform module, quite often we need to print all our variables with all descriptions.

We offer you to avoid long bash scripts and update your README.md as quick as it was released.

## How to Use

Download the latest binary for your OS from [Releases Page](https://github.com/idestis/hcltomd/releases) and run it over required file

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

## Contribute

I will be happy to welcome new contributers for this repo, feel free to read CONTRIBUTE.md
