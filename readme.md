# Readme

This script extracts resource blocks from HCL files for Terraform and saves them as separate text files.

## Usage

### Build

```
go build -o hcl-separator
```

### Install (for Ubuntu)

```
mv ./hcl-separator ~/.local/bin/
```

### Extract resource blocks

```
hcl-separator /tmp/example.tf ./output
```

### Merge extracted files back into a single Terraform file

```
bash ./scripts/merger.sh ./output ./main.tf
```
