package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <input_file> <output_dir>", os.Args[0])
	}

	inputFile := os.Args[1]
	outputDir := os.Args[2]

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Read the HCL file
	src, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", inputFile, err)
	}

	// Parse the HCL file
	file, diags := hclwrite.ParseConfig(src, inputFile, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		log.Fatalf("Failed to parse file %s: %v", inputFile, diags.Error())
	}

	// Iterate over all blocks
	for _, block := range file.Body().Blocks() {
		if block.Type() == "resource" {
			labels := block.Labels()
			if len(labels) < 2 {
				log.Println("Skipping malformed resource block")
				continue
			}

			resourceType := labels[0]
			resourceName := labels[1]
			filename := filepath.Join(outputDir, fmt.Sprintf("%s_%s.hcl", resourceType, resourceName))

			// Create a new HCL file for this resource block
			newFile := hclwrite.NewEmptyFile()
			newBody := newFile.Body()
			newBody.AppendBlock(block)

			// Write to the new file
			if err := os.WriteFile(filename, newFile.Bytes(), 0644); err != nil {
				log.Fatalf("Failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Extracted: %s\n", filename)
		}
	}
}
