# Readme

This script extracts resource blocks from HCL files for Terraform and saves them as separate text files.

## Usage

### Build

```
go build -o hcl-separator
```
#### install

```
mv ./hcl-separator ~/.local/bin/
chmod +x ~/.local/bin/hcl-separator
```

### Separate

```
hcl-separator /tmp/example.tf ./output
```

### Join separated files

```
bash ./scripts/merger.sh ./output ./main.tf
```

