# HCL to Markdown

[![Build Status](https://cloud.drone.io/api/badges/idestis/hcltomd/status.svg)](https://cloud.drone.io/idestis/hcltomd)
## About

To write a good documentation for terraform module, quite often we just need to print all our input variables as a fancy table.

This tool offer you to avoid manual work and print variables as ready-to-go Markdown table. :rocket:

Hope this will save your time. :blush:

## How to Use

Download the latest binary for your OS from [Releases Page](https://github.com/idestis/hcltomd/releases) and run it over file.

Example variable file you can get from the current repository [inputs.tf](./example/inputs.tf)

```tf
variable "aws_region" {
  description = "The AWS region in which all resources will be created."
  type        = "string"
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "The id of the specific VPC to retrieve."
  type        = "string"
  default     = "1"
}

variable "instance_count" {
  description = "The count of desired instances of EC2."
  type        = "number"
  default     = 2
}

variable "zones" {
  description = "The selected zones for deployment."
  type        = "list(string)"
  default = [
    "us-east-1", "us-east-2"
  ]
}
```

Execute with flag `--file`

```bash
$ ./hcltomd --file ./example/inputs.tf
|      NAME      |  TYPE          |        DEFAULT        |                      DESCRIPTION                       |
|----------------|----------------|-----------------------|--------------------------------------------------------|
| aws_region     | string         | us-east-1             | The AWS region in which all resources will be created. |
| vpc_id         | string         | 1                     | The id of the specific VPC to retrieve.                |
| instance_count | number         | 2                     | The count of desired instances of EC2.                 |
| zones          | list(string)   | [us-east-1 us-east-2] | The selected zones                                     |
```

That's it, now you can copy output into your documentation.

## Contribute

I will be happy to welcome a new contributors for this repo, feel free to read CONTRIBUTE.md

## Special Thanks

@olekukonko for [tablewriter](https://github.com/olekukonko/tablewriter) package.

## TBD

- Abillity to select the table format (e.g. Markdown, Confluence Markdown, etc.)
- Near to 100% test coverage.
