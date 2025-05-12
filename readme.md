# Readme

This script extracts only the specified resource block from Terraform HCL file based on the resource name.

## Usage (Ubuntu)

## command

```
tfpick <file> <resource_type.name>
```

e.g.,

```
$ ./tfpick example.tf aws_internet_gateway.main
resource "aws_internet_gateway" "main" { # test comment2
  vpc_id = aws_vpc.main.id # test comment3
  tags = { # test comment4
    Name = "${var.environment}-${var.project}-igw" # test comment5
  } # test comment6
} # test comment7
```

### Build

```
go build -o tfpick
```
