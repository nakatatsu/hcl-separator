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
{
  "start_line": 9,
  "end_line": 14,
  "content": "resource \"aws_internet_gateway\" \"main\" { # test comment2\n  vpc_id = aws_vpc.main.id # test comment3\n  tags = { # test comment4\n    Name = \"${var.environment}-${var.project}-igw\" # test comment5\n  } # test comment6\n} # test comment7"
}
```

### Build

```
go build -o tfpick
```
