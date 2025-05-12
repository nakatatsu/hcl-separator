# Readme

This script extracts only the specified resource block from Terraform HCL file based on the resource name or line no.

## Usage (Ubuntu)

## command

```
tfpick <file> <resource_type.name>|<line no>
```

e.g.,

```
$ tfpick main.tf aws_alb.api
{
  "start_line": 26,
  "end_line": 40,
  "content": "resource \"aws_alb\" \"api\" {\n  name                       = \"${var.environment}-${var.project}-api\"\n  internal                   = false\n  load_balancer_type         = \"application\"\n  security_groups            = [var.api_alb_security_group_id]\n  subnets                    = var.public_subnet_ids\n  enable_deletion_protection = true\n  drop_invalid_header_fields = true\n\n  access_logs {\n    bucket  = aws_s3_bucket.logs.bucket\n    prefix  = \"${var.environment}-${var.project}-api-alb-logs\"\n    enabled = true\n  }\n}"
}
$ tfpick main.tf 32
{
  "start_line": 26,
  "end_line": 40,
  "content": "resource \"aws_alb\" \"api\" {\n  name                       = \"${var.environment}-${var.project}-api\"\n  internal                   = false\n  load_balancer_type         = \"application\"\n  security_groups            = [var.api_alb_security_group_id]\n  subnets                    = var.public_subnet_ids\n  enable_deletion_protection = true\n  drop_invalid_header_fields = true\n\n  access_logs {\n    bucket  = aws_s3_bucket.logs.bucket\n    prefix  = \"${var.environment}-${var.project}-api-alb-logs\"\n    enabled = true\n  }\n}"
}
$
```

### Build

```
go build -o tfpick

```
